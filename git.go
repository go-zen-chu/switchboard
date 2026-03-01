//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gogithttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// DefaultGitUserEmail is the default git commit author email used when --update is set
const DefaultGitUserEmail = "action@github.com"

// DefaultGitUserName is the default git commit author name used when --update is set
const DefaultGitUserName = "GitHub Action"

// DefaultGitCommitMessage is the default git commit message used when --update is set
const DefaultGitCommitMessage = "[skip ci] Update Bluesky to X sync info"

// Gitter handles git commit and push operations
type Gitter interface {
	CommitAndPush(message string) error
}

type gitter struct {
	userEmail string
	userName  string
}

// NewGitter creates a new Gitter with the given user email and name
func NewGitter(userEmail, userName string) Gitter {
	return &gitter{
		userEmail: userEmail,
		userName:  userName,
	}
}

// CommitAndPush stages all changes, commits with the given message, and pushes to origin.
// If there are no changes to commit, it logs and returns nil.
func (g *gitter) CommitAndPush(message string) error {
	repo, err := gogit.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("opening git repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	if err := wt.AddWithOptions(&gogit.AddOptions{All: true}); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	status, err := wt.Status()
	if err != nil {
		return fmt.Errorf("git status: %w", err)
	}
	if status.IsClean() {
		slog.Info("No changes to commit, skipping git commit/push")
		return nil
	}

	_, err = wt.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  g.userName,
			Email: g.userEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	slog.Info("Committed changes", "message", message)

	pushOpts := &gogit.PushOptions{}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		pushOpts.Auth = &gogithttp.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}
	}

	if err := repo.Push(pushOpts); err != nil {
		if err == gogit.NoErrAlreadyUpToDate {
			slog.Info("Already up to date")
			return nil
		}
		return fmt.Errorf("git push: %w", err)
	}
	slog.Info("Pushed changes to remote")
	return nil
}
