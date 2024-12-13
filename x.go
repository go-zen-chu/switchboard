package switchboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

type XClient interface {
	Post(ctx context.Context, content string) (*XPost, error)
}

type XPost struct {
	ID string
}

type xclient struct {
	gotwiCli *gotwi.Client
}

type BearerAuthorizer struct {
	Token string
}

func (a BearerAuthorizer) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func NewXClient(ctx context.Context, oauthToken, oauthTokenSecret, apiKey, apiKeySecret string) (XClient, error) {
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
		return nil, fmt.Errorf("managetweet create tweet: %w", err)
	}
	p := &XPost{
		ID: *res.Data.ID,
	}
	return p, nil
}
