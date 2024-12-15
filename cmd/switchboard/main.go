package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-zen-chu/switchboard/cmd/switchboard/cmd"
	"github.com/go-zen-chu/switchboard/internal/di"
	"github.com/spf13/cobra"
)

func main() {
	dic := di.NewContainer()
	app, err := NewApp(dic)
	if err != nil {
		slog.Error("fail init app", "error", err)
		os.Exit(1)
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("while running app", "error", err)
		os.Exit(1)
	}
}

type app struct {
	ctx     context.Context
	rootCmd *cobra.Command
}

func NewApp(req cmd.SwitchboardRequirements) (*app, error) {
	ctx := context.Background()
	return &app{
		ctx:     ctx,
		rootCmd: cmd.NewRootCmd(req),
	}, nil
}

func (a *app) Run(args []string) error {
	a.rootCmd.SetArgs(args[1:])
	if err := a.rootCmd.Execute(); err != nil {
		return fmt.Errorf("root command: %w", err)
	}
	return nil
}
