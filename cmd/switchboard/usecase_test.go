package main

import (
	"context"
	"os"
	"testing"

	"github.com/go-zen-chu/switchboard"
	"github.com/go-zen-chu/switchboard/cmd/switchboard/cmd"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "If help flag given, show help",
			args:    []string{"switchboard", "--help"},
			wantErr: false,
		},
		{
			name:    "If bluesky2x subcommand used, sync bluesky post to x",
			args:    []string{"switchboard", "bluesky2x", "-v"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			bcli, _ := switchboard.NewBlueskyClient(
				ctx,
				os.Getenv("BLUESKY_IDENTIFIER"),
				os.Getenv("BLUESKY_PASSWORD"),
			)
			xcli, _ := switchboard.NewXClient(
				ctx,
				os.Getenv("X_ACCESS_TOKEN"),
				os.Getenv("X_ACCESS_SECRET"),
				os.Getenv("X_API_KEY"),
				os.Getenv("X_API_SECRET"),
			)
			app, goterr := NewApp(&cmd.SwitchboardRequirements{
				Ctx:           ctx,
				BlueskyClient: bcli,
				XClient:       xcli,
			})
			if (goterr != nil) != tt.wantErr {
				t.Errorf("NewApp error = %v, wantErr %v", goterr, tt.wantErr)
				return
			}
			goterr = app.Run(tt.args)
			if (goterr != nil) != tt.wantErr {
				t.Errorf("app.Run() error = %v, wantErr %v", goterr, tt.wantErr)
				return
			}
		})
	}
}
