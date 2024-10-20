/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-zen-chu/switchboard"
	"github.com/spf13/cobra"
)

const (
	// numSyncPosts is the number of posts to sync from Bluesky to X. Be careful not to exceed the rate limit.
	// https://developer.x.com/en/docs/x-api/lists/list-tweets/introduction
	numSyncPosts = 10
)

func NewBluesky2XCmd(ctx context.Context, bcli switchboard.BlueskyClient, xcli switchboard.XClient) *cobra.Command {
	// bluesky2xCmd represents the bluesky2x command
	var bluesky2xCmd = &cobra.Command{
		Use:   "bluesky2x",
		Short: "Send bluesky post to x",
		Long:  `Send bluesky post to x`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get posts from bluesky
			bposts, err := bcli.GetMyLatestPostsCreatedAsc(ctx, numSyncPosts)
			if err != nil {
				return fmt.Errorf("getting latest posts from Bluesky: %w\n", err)
			}
			// get posts from x
			xposts, err := xcli.GetMyLatestPostsCreatedAsc(ctx, numSyncPosts)
			if err != nil {
				return fmt.Errorf("getting latest posts from X: %w\n", err)
			}
			slog.Info("Got posts", "bluesky", bposts, "x", xposts)

			// post to x that are not in json

			// // save sent posts to json and git push
			// for _, bpost := range bposts {
			// 	posted[bpost.Cid] = bpost
			// }
			// jsonBytes, err := json.MarshalIndent(posted, "", "  ")
			// if err != nil {
			// 	return fmt.Errorf("marshaling posts to json: %w\n", err)
			// }
			// err = os.WriteFile(postStorePath, jsonBytes, 0644)
			// if err != nil {
			// 	return fmt.Errorf("writing posts to json: %w\n", err)
			// }
			// repo, err := git.PlainOpen(".")
			// if err != nil {
			// 	return fmt.Errorf("opening git repo: %w\n", err)
			// }
			// wt, err := repo.Worktree()
			// if err != nil {
			// 	return fmt.Errorf("getting worktree: %w\n", err)
			// }
			// _, err = wt.Add(postStorePath)
			// if err != nil {
			// 	return fmt.Errorf("adding post json to git: %w\n", err)
			// }
			// commit, err := wt.Commit("[skip ci] Add sent bluesky posts", &git.CommitOptions{
			// 	Author: &object.Signature{
			// 		Name:  "bluesky2x",
			// 		Email: "bluesky2x@github.com",
			// 		When:  time.Now(),
			// 	},
			// })
			// if err != nil {
			// 	return fmt.Errorf("committing post json to git: %w\n", err)
			// }
			// err = repo.Push(&git.PushOptions{
			// 	RemoteName: "origin",
			// })
			// if err != nil {
			// 	return fmt.Errorf("pushing post json to git: %w\n", err)
			// }
			// slog.Info("Pushed post json to git", "commit", commit)
			return nil
		},
	}
	return bluesky2xCmd
}
