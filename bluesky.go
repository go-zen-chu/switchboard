//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
)

const BlueskyHost = "https://bsky.social"

type BlueskyClient interface {
	GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int64) ([]BlueskyPost, error)
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

func (bc *blueskyClient) GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int64) ([]BlueskyPost, error) {
	profile, err := bsky.ActorGetProfile(ctx, bc.cli, bc.cli.Auth.Did)
	if err != nil {
		return nil, fmt.Errorf("Bluesky get profile: %w", err)
	}
	feeds, err := bsky.FeedGetAuthorFeed(ctx, bc.cli, profile.Did, "", "", false, numPosts)
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
		uriParts := strings.Split(string(f.Post.Uri), "/")
		url := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", uriParts[2], uriParts[4])
		posts = append(posts, BlueskyPost{
			Cid:       f.Post.Cid,
			Content:   fp.Text,
			CreatedAt: ca,
			URL:       url,
			Reply:     rep,
		})
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
	return posts, nil
}
