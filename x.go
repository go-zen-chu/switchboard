package switchboard

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
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
