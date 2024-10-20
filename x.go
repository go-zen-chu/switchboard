package switchboard

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/timeline"
	"github.com/michimani/gotwi/tweet/timeline/types"
)

type XClient interface {
	GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int64) ([]XPost, error)
}

type XPost struct {
	ID      string
	Content string
}

type xclient struct {
	xID string
	cli *gotwi.Client
}

func NewXClient(ctx context.Context, xID, oauthToken, oauthTokenSecret, apiKey, apiKeySecret string) (XClient, error) {
	xcli, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           oauthToken,
		OAuthTokenSecret:     oauthTokenSecret,
		APIKey:               apiKey,
		APIKeySecret:         apiKeySecret,
	})
	if err != nil {
		return nil, fmt.Errorf("creating X client: %w", err)
	}
	return &xclient{
		xID: xID,
		cli: xcli,
	}, nil
}

func (c *xclient) GetMyLatestPostsCreatedAsc(ctx context.Context, numPosts int64) ([]XPost, error) {
	tweets, err := timeline.ListTweets(ctx, c.cli, &types.ListTweetsInput{
		ID:         c.xID,
		MaxResults: types.ListMaxResults(numPosts),
	})
	if err != nil {
		return nil, fmt.Errorf("getting latest posts from X: %w", err)
	}
	posts := make([]XPost, 0, len(tweets.Data))
	for _, tweet := range tweets.Data {
		posts = append(posts, XPost{
			ID:      *tweet.ID,
			Content: *tweet.Text,
		})
	}
	return posts, nil
}
