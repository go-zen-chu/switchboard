//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
)

const BlueskyHost = "https://bsky.social"

type BlueskyClient interface {
	GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int) ([]BlueskyPost, error)
}

type blueskyClient struct {
	cli *xrpc.Client
}

type BlueskyPost struct {
	Cid       string        `json:"cid"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	URL       string        `json:"url"`
	Reply     *BlueskyReply `json:"reply,omitempty"`
}

type BlueskyReply struct {
	RootCid   string `json:"root_cid"`
	ParentCid string `json:"parent_cid"`
}

func NewBlueskyClient(ctx context.Context, identifier, password string) (BlueskyClient, error) {
	for k, v := range map[string]string{
		"identifier": identifier,
		"password":   password,
	} {
		if v == "" {
			return nil, fmt.Errorf("%s is empty", k)
		}
	}
	bcli := &xrpc.Client{
		Host: BlueskyHost,
	}
	sout, err := atproto.ServerCreateSession(ctx, bcli, &atproto.ServerCreateSession_Input{
		Identifier: identifier,
		Password:   password,
	})
	if err != nil {
		return nil, fmt.Errorf("Bluesky create session: %w", err)
	}
	bcli.Auth = &xrpc.AuthInfo{
		AccessJwt:  sout.AccessJwt,
		RefreshJwt: sout.RefreshJwt,
		Handle:     sout.Handle,
		Did:        sout.Did,
	}
	return &blueskyClient{
		cli: bcli,
	}, nil
}

func (bc *blueskyClient) GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int) ([]BlueskyPost, error) {
	profile, err := bsky.ActorGetProfile(ctx, bc.cli, bc.cli.Auth.Did)
	if err != nil {
		return nil, fmt.Errorf("Bluesky get profile: %w", err)
	}
	feeds, err := bsky.FeedGetAuthorFeed(ctx, bc.cli, profile.Did, "", "", false, int64(numPosts))
	if err != nil {
		return nil, fmt.Errorf("Bluesky get author feed: %w", err)
	}
	posts := make([]BlueskyPost, 0, len(feeds.Feed))
	for _, f := range feeds.Feed {
		fp := f.Post.Record.Val.(*bsky.FeedPost)
		if fp == nil {
			return nil, fmt.Errorf("cast bluesky feed post results nil: %v", f.Post.Record)
		}
		ca, err := time.Parse("2006-01-02T15:04:05.999Z", fp.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("parse bluesky post CreatedAt(%s): %w", fp.CreatedAt, err)
		}
		var rep *BlueskyReply
		if fp.Reply != nil {
			rep = &BlueskyReply{
				RootCid:   fp.Reply.Root.Cid,
				ParentCid: fp.Reply.Parent.Cid,
			}
		}
		// urlReplacedContent := replaceAbbreviatedURLToOriginal(fp)

		posts = append(posts, BlueskyPost{
			Cid:       f.Post.Cid,
			Content:   replaceAbbreviatedURLToOriginal(fp),
			CreatedAt: ca,
			URL:       buildPostURL(f),
			Reply:     rep,
		})
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
	return posts, nil
}

func buildPostURL(feed *bsky.FeedDefs_FeedViewPost) string {
	if feed.Post == nil {
		return ""
	}
	uriParts := strings.Split(string(feed.Post.Uri), "/")
	if len(uriParts) < 5 {
		return ""
	}
	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", uriParts[2], uriParts[4])
}

func replaceAbbreviatedURLToOriginal(feedPost *bsky.FeedPost) string {
	if feedPost.Facets == nil {
		return feedPost.Text
	}
	resultText := feedPost.Text
	tokenMap := make(map[string]string)
	tokenId := 0
	for _, facets := range feedPost.Facets {
		if facets.Features == nil {
			continue
		}
		for _, feature := range facets.Features {
			if feature.RichtextFacet_Link == nil {
				continue
			}
			originalURLStr := feature.RichtextFacet_Link.Uri
			originalURL, err := url.Parse(originalURLStr)
			if err != nil {
				slog.Warn("parse original url failed", "url", originalURLStr, "replacedText", resultText, "error", err)
				continue
			}
			host := originalURL.Hostname()
			hostRegexp := strings.ReplaceAll(host, ".", `\.`)
			abbrevURLRegexp := hostRegexp + `.*?\.\.\.`
			re, err := regexp.Compile(abbrevURLRegexp)
			if err != nil {
				slog.Warn("compile regexp failed", "regexp", abbrevURLRegexp, "replacedText", resultText, "error", err)
				continue
			}
			// NOTES: temporary replace abbreviated URL with token otherwise if we have same abbreviated URL in the post, it will be misreplaced
			// e.g. `github.com/... github.com/...` -> `https://https://github.com/test github.com/...`
			isFirstMatch := true
			tokenString := fmt.Sprintf("<<swbtoken%d>>", tokenId)
			resultText = re.ReplaceAllStringFunc(resultText, func(match string) string {
				if isFirstMatch {
					isFirstMatch = false
					return tokenString
				}
				return match
			})
			tokenMap[tokenString] = originalURLStr
			tokenId++
		}
	}
	for token, url := range tokenMap {
		resultText = strings.ReplaceAll(resultText, token, url)
	}
	return resultText
}
