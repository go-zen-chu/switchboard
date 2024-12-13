package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-zen-chu/switchboard"
	"github.com/go-zen-chu/switchboard/cmd/switchboard/cmd"
	"github.com/spf13/cobra"
)

func main() {
	swq, err := newSwitchboardRequirements()
	if err != nil {
		slog.Error("fail init requirements", "error", err)
		os.Exit(1)
	}
	app, err := NewApp(swq)
	if err != nil {
		slog.Error("fail init app", "error", err)
		os.Exit(1)
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("while running app", "error", err)
		os.Exit(1)
	}
}

func newSwitchboardRequirements() (*cmd.SwitchboardRequirements, error) {
	ctx := context.Background()
	bcli, err := switchboard.NewBlueskyClient(
		ctx,
		os.Getenv("BLUESKY_IDENTIFIER"),
		os.Getenv("BLUESKY_PASSWORD"),
	)
	if err != nil {
		return nil, fmt.Errorf("init bluesky client: %w", err)
	}
	xcli, err := switchboard.NewXClient(
		ctx,
		os.Getenv("X_ID"),
		os.Getenv("X_ACCESS_TOKEN"),
		os.Getenv("X_ACCESS_SECRET"),
		os.Getenv("X_API_KEY"),
		os.Getenv("X_API_SECRET"),
		os.Getenv("X_BEARER_TOKEN"),
	)
	if err != nil {
		return nil, fmt.Errorf("init x client: %w", err)
	}
	return &cmd.SwitchboardRequirements{
		Ctx:           ctx,
		BlueskyClient: bcli,
		XClient:       xcli,
	}, nil
}

type app struct {
	ctx     context.Context
	rootCmd *cobra.Command
}

func NewApp(switchboardReq *cmd.SwitchboardRequirements) (*app, error) {
	ctx := context.Background()
	return &app{
		ctx:     ctx,
		rootCmd: cmd.NewRootCmd(switchboardReq),
	}, nil
}

func (a *app) Run(args []string) error {
	a.rootCmd.SetArgs(args[1:])
	if err := a.rootCmd.Execute(); err != nil {
		return fmt.Errorf("root command: %w", err)
	}
	return nil
}
