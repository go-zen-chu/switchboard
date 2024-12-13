package switchboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	v2 "github.com/g8rswimmer/go-twitter/v2"
	"github.com/michimani/gotwi"
)

type XClient interface {
	Post(ctx context.Context, content string) (*XPost, error)
}

type XPost struct {
	ID string
}

type xclient struct {
	xID       string
	oauth1cli *gotwi.Client
	gotwiCli  *twitter.Client
	twCli     *v2.Client
}

type BearerAuthorizer struct {
	Token string
}

func (a BearerAuthorizer) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func NewXClient(ctx context.Context, xID, oauthToken, oauthTokenSecret, apiKey, apiKeySecret, bearerToken string) (XClient, error) {
	oauth1cli, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           oauthToken,
		OAuthTokenSecret:     oauthTokenSecret,
		APIKey:               apiKey,
		APIKeySecret:         apiKeySecret,
		Debug:                true,
	})
	if err != nil {
		return nil, fmt.Errorf("oauth1 X client: %w", err)
	}

	cnf := oauth1.NewConfig(apiKey, apiKeySecret)
	token := oauth1.NewToken(oauthToken, oauthTokenSecret)
	httpClient := cnf.Client(oauth1.NoContext, token)
	gotwiCli := twitter.NewClient(httpClient)

	twCli := &v2.Client{
		Authorizer: BearerAuthorizer{
			Token: bearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	return &xclient{
		xID:       xID,
		oauth1cli: oauth1cli,
		gotwiCli:  gotwiCli,
		twCli:     twCli,
	}, nil
}

// func (c *xclient) GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int64) ([]XPost, error) {
// 	tweets, err := timeline.ListTweets(ctx, c.cli, &types.ListTweetsInput{
// 		ID:         c.xID,
// 		MaxResults: types.ListMaxResults(numPosts),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("getting latest posts from X: %w", err)
// 	}
// 	posts := make([]XPost, 0, len(tweets.Data))
// 	for _, tweet := range tweets.Data {
// 		posts = append(posts, XPost{
// 			ID:      *tweet.ID,
// 			Content: *tweet.Text,
// 		})
// 	}

// 	return posts, nil
// }

func (c *xclient) Post(ctx context.Context, content string) (*XPost, error) {
	// // TODO: content length must be < 280 letters
	// ci := &types.CreateInput{
	// 	Text: gotwi.String(content),
	// }
	// // TODO: support reply
	// res, err := managetweet.Create(ctx, c.oauth1cli, ci)
	// if err != nil {
	// 	return nil, fmt.Errorf("post tweet: %w", err)
	// }
	// p := &XPost{
	// 	ID: *res.Data.ID,
	// }
	res, _, err := c.gotwiCli.Statuses.Update(content, nil)
	if err != nil {
		return nil, fmt.Errorf("post tweet: %w", err)
	}
	p := &XPost{
		ID: string(res.ID),
	}
	return p, nil
}
