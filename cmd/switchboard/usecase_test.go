package main

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-zen-chu/switchboard"
	"github.com/go-zen-chu/switchboard/cmd/switchboard/cmd"

	"go.uber.org/mock/gomock"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		customizeMock func(mockBCli *switchboard.MockBlueskyClient, mockXCli *switchboard.MockXClient)
		wantErr       bool
	}{
		{
			name:    "If help flag given, show help",
			args:    []string{"switchboard", "--help"},
			wantErr: false,
		},
		{
			name: "If bluesky2x subcommand used, sync bluesky post to x",
			args: []string{"switchboard", "bluesky2x", "-v"},
			customizeMock: func(mockBCli *switchboard.MockBlueskyClient, mockXCli *switchboard.MockXClient) {
				mockBCli.EXPECT().GetMyLatestPostsCreatedAsc(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]switchboard.BlueskyPost{
						{
							Cid:       "test1test1test1test1test1test1test1test1test1test1test1test1",
							Content:   "test1",
							CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							URL:       "https://bsky.app/profile/did:plc:test1test1test1test1/post/test1test1",
							Reply:     nil,
						},
						{
							Cid:       "test2test2test2test2test2test2test2test2test2test2test2test2",
							Content:   "test2",
							CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
							URL:       "https://bsky.app/profile/did:plc:test2test2test2test2/post/test2test2",
							Reply:     nil,
						},
					}, nil)
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test1.*")).
					Return(&switchboard.XPost{
						ID: "1111111111111111111",
					}, nil)
				slog.Warn("test")
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test2.*")).
					Return(&switchboard.XPost{
						ID: "2222222222222222222",
					}, nil)
				slog.Warn("test2")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockBCli := switchboard.NewMockBlueskyClient(c)
			mockXCli := switchboard.NewMockXClient(c)

			if tt.customizeMock != nil {
				tt.customizeMock(mockBCli, mockXCli)
			}
			app, goterr := NewApp(&cmd.SwitchboardRequirements{
				Ctx:           context.Background(),
				BlueskyClient: mockBCli,
				XClient:       mockXCli,
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
			// cleanup
			if _, err := os.Stat("output"); err == nil {
				err = os.RemoveAll("output")
				if err != nil {
					t.Errorf("cleanup remove all error = %v", err)
					return
				}
			} else {
				if !os.IsNotExist(err) {
					t.Errorf("stat directory error = %v", err)
					return
				}
			}
		})
	}
}
