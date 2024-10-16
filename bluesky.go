package main

import (
	"context"
	"fmt"
	"sort"
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
	Cid       string    `json:"cid"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Uri       string    `json:"uri"`
}

func NewBlueskyClient(ctx context.Context, identifier, password string) (BlueskyClient, error) {
	bcli := &xrpc.Client{
		Host: BlueskyHost,
	}
	if identifier == "" {
		return nil, fmt.Errorf("identifier is empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password is empty")
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
		posts = append(posts, BlueskyPost{
			Cid:       f.Post.Cid,
			Content:   fp.Text,
			CreatedAt: ca,
			Uri:       f.Post.Uri,
		})
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
	return posts, nil
}
