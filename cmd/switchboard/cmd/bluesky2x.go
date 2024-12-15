//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"text/template"

	_ "embed"

	"github.com/go-zen-chu/switchboard"
	"github.com/spf13/cobra"
)

type Bluesky2XCmdRequirements interface {
	Context() context.Context
	BlueskyClient() switchboard.BlueskyClient
	XClient() switchboard.XClient
}

func NewBluesky2XCmd(req Bluesky2XCmdRequirements) *cobra.Command {
	const defaultGenWorkflowFile = false
	var genWorkflowFile bool

	const defaultNumSyncLatestPosts = 10
	var numSyncLatestPosts int

	// bluesky2xCmd represents the bluesky2x command
	var bluesky2xCmd = &cobra.Command{
		Use:   "bluesky2x",
		Short: "Send bluesky post to x",
		Long:  `Send bluesky post to x`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if genWorkflowFile {
				if err := generateWorkflowFile(); err != nil {
					return fmt.Errorf("generating workflow file: %w", err)
				}
				return nil
			}
			ctx := req.Context()
			bcli := req.BlueskyClient()
			xcli := req.XClient()
			if err := syncBlueskyLatestPosts2X(ctx, bcli, xcli, numSyncLatestPosts); err != nil {
				return fmt.Errorf("syncing bluesky latest posts to x: %w", err)
			}
			return nil
		},
	}
	bluesky2xCmd.Flags().BoolVar(&genWorkflowFile, "gen-workflow-file", defaultGenWorkflowFile, "Generate workflow files to .github/workflows folder")
	bluesky2xCmd.Flags().IntVarP(&numSyncLatestPosts, "num-sync-posts", "n", defaultNumSyncLatestPosts, `Number of latest posts to sync from Bluesky to X per run.
Make sure not to exceed the rate limit (Especially, if you are using Free plan).
https://developer.x.com/en/docs/x-api/lists/list-tweets/introduction
`)
	return bluesky2xCmd
}

func syncBlueskyLatestPosts2X(ctx context.Context, bcli switchboard.BlueskyClient, xcli switchboard.XClient, numSyncLatestPosts int) error {
	slog.Info("Start syncing latest posts from bluesky to X")
	bposts, err := bcli.GetMyLatestPostsCreatedAsc(ctx, numSyncLatestPosts)
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
		slog.Info("No new posts. Finished.")
		return nil
	}

	for _, bpost := range newPosts {
		cnt := fmt.Sprintf("%s\n🤖from🦋: %s", bpost.Content, bpost.URL)
		xpost, err := xcli.Post(ctx, cnt)
		if err != nil {
			var errXDup *switchboard.ErrXDuplicatePost
			if errors.As(err, &errXDup) {
				slog.Warn("Found duplicate tweet in X", "content", cnt)
				continue
			}
			return fmt.Errorf("post tweet: %w\n", err)
		}
		slog.Debug("Posted tweet", "cid", bpost.Cid, "tweet id", xpost.ID, "content", cnt)

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
		slog.Debug("Updated sync info")
	}
	slog.Info("Finished syncing from bluesky to X")
	return nil
}

type WorkflowInfo struct {
	WorkflowName string
	WorkflowOn   string
}

//go:embed template/workflow-template.yml
var wfTmplStr string

func generateWorkflowFile() error {
	tmpl, err := template.New("workflow").Parse(wfTmplStr)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	wfCronName := "bluesky2x-cron"
	wfCronFile, err := os.Create(".github/workflows/" + wfCronName + ".yml")
	if err != nil {
		return fmt.Errorf("creating %s file: %w", wfCronName, err)
	}
	defer wfCronFile.Close()
	wfCronInfo := WorkflowInfo{
		WorkflowName: wfCronName,
		WorkflowOn: `
on:
  # NOTE: Configuring less than 1 hour may cause 429 Too Many Request from twitter api
  schedule:
    - cron: "15 * * * *"`,
	}
	if err := tmpl.Execute(wfCronFile, wfCronInfo); err != nil {
		return fmt.Errorf("executing %s template: %w", wfCronName, err)
	}

	wfManualName := "bluesky2x-manual"
	wfManualFile, err := os.Create(".github/workflows/" + wfManualName + ".yml")
	if err != nil {
		return fmt.Errorf("creating %s file: %w", wfManualName, err)
	}
	defer wfManualFile.Close()
	wfManualInfo := WorkflowInfo{
		WorkflowName: wfManualName,
		WorkflowOn: `
on:
  workflow_dispatch`,
	}
	if err := tmpl.Execute(wfManualFile, wfManualInfo); err != nil {
		return fmt.Errorf("executing %s template: %w", wfManualName, err)
	}

	return nil
}
