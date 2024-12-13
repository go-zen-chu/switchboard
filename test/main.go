package main

import (
	"context"
	"fmt"
	"os"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

func main() {
	cli, err := gotwi.NewClient(&gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           os.Getenv("X_ACCESS_TOKEN"),
		OAuthTokenSecret:     os.Getenv("X_ACCESS_SECRET"),
		APIKey:               os.Getenv("X_API_KEY"),
		APIKeySecret:         os.Getenv("X_API_SECRET"),
		Debug:                true,
	})
	if err != nil {
		panic(err)
	}
	p := &types.CreateInput{
		Text: gotwi.String("Hello from bluesky"),
	}
	res, err := managetweet.Create(context.Background(), cli, p)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Data.ID)
}
