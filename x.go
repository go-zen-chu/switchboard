//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"golang.org/x/text/unicode/norm"
)

const (
	XMaxTweetLength      = 280
	XShortenedLinkLength = 23
)

type XPost struct {
	ID string
}

type ErrXDuplicatePost struct {
	GoTwiError *gotwi.GotwiError
}

func (e *ErrXDuplicatePost) Error() string {
	if e.GoTwiError == nil {
		return "unexpected error: gotwi error is nil"
	}
	return fmt.Sprintf("duplicate post exists in X (status code %d, title %s, detail %s)", e.GoTwiError.StatusCode, e.GoTwiError.Title, e.GoTwiError.Detail)
}

type XClient interface {
	Post(ctx context.Context, content string) (*XPost, error)
}

type xclient struct {
	gotwiCli *gotwi.Client
}

func NewXClient(ctx context.Context, oauthToken, oauthTokenSecret, apiKey, apiKeySecret string) (XClient, error) {
	for k, v := range map[string]string{
		"oauthToken":       oauthToken,
		"oauthTokenSecret": oauthTokenSecret,
		"apiKey":           apiKey,
		"apiKeySecret":     apiKeySecret,
	} {
		if v == "" {
			return nil, fmt.Errorf("%s is empty", k)
		}
	}
	gotwiCli, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           oauthToken,
		OAuthTokenSecret:     oauthTokenSecret,
		APIKey:               apiKey,
		APIKeySecret:         apiKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("init gotwi client: %w", err)
	}

	return &xclient{
		gotwiCli: gotwiCli,
	}, nil
}

func (c *xclient) Post(ctx context.Context, content string) (*XPost, error) {
	// TODO: content length must be < 280 letters
	ci := &types.CreateInput{
		Text: gotwi.String(content),
	}
	// TODO: support reply
	res, err := managetweet.Create(ctx, c.gotwiCli, ci)
	if err != nil {
		ge := err.(*gotwi.GotwiError)
		if !ge.OnAPI {
			return nil, fmt.Errorf("managetweet create tweet: %w", err)
		}
		slog.Warn("create tweet",
			"error title", ge.Title,
			"error detail", ge.Detail,
			"error type", ge.Type,
			"error status", ge.Status,
			"error status code", ge.StatusCode,
		)
		if strings.Contains(ge.Detail, "not allowed to create a Tweet with duplicate content") {
			return nil, &ErrXDuplicatePost{GoTwiError: ge}
		}
		return nil, fmt.Errorf("managetweet create tweet: %w", err)
	}
	p := &XPost{
		ID: *res.Data.ID,
	}
	return p, nil
}

var (
	urlRegex   = regexp.MustCompile(`https?://\S+`)
	emojiRegex = regexp.MustCompile(`[\p{So}\p{Sk}][\p{Mn}\p{Me}\x{FE0F}\x{20E3}]?[\x{200D}\p{Zs}]?[\p{So}\p{Sk}]?`)
	cjkRegex   = regexp.MustCompile(`\p{Han}|\p{Hiragana}|\p{Katakana}|\p{Hangul}`)
)

// countTweetCharacters counts the number of characters in a tweet.
// The actual counting algorithm is described here: https://developer.x.com/en/docs/counting-characters
// This function does not aim to follow the actual counting algorithm above.
func CountTweetCharacters(content string) int {
	normText := norm.NFC.String(content)
	// X count any URL to 23 characters
	urlReplacedText := urlRegex.ReplaceAllString(normText, strings.Repeat("x", XShortenedLinkLength))
	countX := 0
	for _, r := range urlReplacedText {
		switch {
		case emojiRegex.MatchString(string(r)):
			countX += 2
		case cjkRegex.MatchString(string(r)):
			countX += 2
		default:
			countX++
		}
	}
	return countX
}

func TruncateTweet(content string, suffixLength int) string {
	ellipsis := "..."
	normText := norm.NFC.String(content)
	countX := 0
	countStr := 0
	for _, r := range normText {
		switch {
		case emojiRegex.MatchString(string(r)):
			countX += 2
		case cjkRegex.MatchString(string(r)):
			countX += 2
		default:
			countX++
		}
		// when countX surpassed the limit, truncate normText with one character before
		if countX >= XMaxTweetLength-suffixLength-len(ellipsis) {
			fmt.Printf("len(content): %d, len(normText): %d, countX: %d, countStr: %d, string: %s\n", len(content), len(normText), countX, countStr, normText[:countStr])
			return normText[:countStr] + ellipsis
		}
		countStr += len(string(r))
	}
	return content
}
