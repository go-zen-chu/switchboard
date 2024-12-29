//go:build mage
// +build mage

package main

import (
	"fmt"
	"log/slog"
	"os"

	gbt "github.com/go-zen-chu/go-build-tools"
)

const currentVersion = "0.0.10"
const currentTagVersion = "v" + currentVersion

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// GitPushTag pushes current tag to remote repository
func GitPushTag(releaseComment string) error {
	err := gbt.GitPushTag(currentTagVersion, releaseComment)
	if err != nil {
		return fmt.Errorf("git push tag: %w", err)
	}
	slog.Info("successfully git push tag", "tag", currentTagVersion)
	return nil
}
