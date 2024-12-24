//go:build mage
// +build mage

package main

import (
	"log/slog"
	"os"

	gbt "github.com/go-zen-chu/go-build-tools"
)

const currentVersion = "0.0.6"
const currentTagVersion = "v" + currentVersion

const imageRegistry = "amasuda"
const repository = "switchboard"

func init() {
	// by default, magefile does not output stderr
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// GitPushTag pushes current tag to remote repository
func GitPushTag(releaseComment string) error {
	return gbt.GitPushTag(currentTagVersion, releaseComment)
}
