//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-zen-chu/switchboard"
	"github.com/spf13/cobra"
)

const (
	// numSyncPosts is the number of posts to sync from Bluesky to X.
	// Make sure not to exceed the rate limit.
	// https://developer.x.com/en/docs/x-api/lists/list-tweets/introduction
	numSyncPosts = 3 //50
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
			slog.Debug("Got posts", "bluesky", bposts)

			stor := switchboard.NewStorer()
			syncInfo, err := stor.LoadSyncInfo()
			if err != nil {
				return fmt.Errorf("loading sync info: %w\n", err)
			}

			newPosts := make([]switchboard.BlueskyPost, 0, len(bposts))
			pm := syncInfo.GetPostMap()
			for _, bpost := range bposts {
				if _, ok := pm[bpost.Cid]; ok {
					slog.Debug("Post already sent", "content", bpost.Content, "cid", bpost.Cid)
					continue
				}
				newPosts = append(newPosts, bpost)
			}

			if len(newPosts) == 0 {
				slog.Info("No new posts")
				return nil
			}

			for _, bpost := range newPosts {
				cnt := fmt.Sprintf("%s\n🤖from🦋: %s", bpost.Content, bpost.URL)

				xpost, err := xcli.Post(ctx, cnt)
				if err != nil {
					return fmt.Errorf("post tweet: %w\n", err)
				}
				slog.Debug("Posted tweet", "cid", bpost.Cid, "xpost", xpost)

				// Store sync info for when got an error while processing (& retry)
				stor.SyncInfo.Posts = append(stor.SyncInfo.Posts, switchboard.PostInfo{
					BlueskyCid:           bpost.Cid,
					TweetID:              xpost.ID,
					Content:              bpost.Content,
					BlueskyPostCreatedAt: bpost.CreatedAt,
				})
				if err := stor.StoreSyncInfo(); err != nil {
					return fmt.Errorf("storing sync info: %w\n", err)
				}
			}

			return nil
		},
	}
	return bluesky2xCmd
}
