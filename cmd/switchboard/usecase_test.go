package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-zen-chu/switchboard"
	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/resources"

	"go.uber.org/mock/gomock"
)

type switchboardRequirementsForTest struct {
	ctx  context.Context
	bcli switchboard.BlueskyClient
	xcli switchboard.XClient
}

func (s *switchboardRequirementsForTest) Context() context.Context {
	return s.ctx
}

func (s *switchboardRequirementsForTest) BlueskyClient() switchboard.BlueskyClient {
	return s.bcli
}

func (s *switchboardRequirementsForTest) XClient() switchboard.XClient {
	return s.xcli
}

func TestMain(t *testing.T) {

	cleanupOutputDir := func(t *testing.T) {
		if _, err := os.Stat("output"); err == nil {
			err = os.RemoveAll("output")
			if err != nil {
				t.Errorf("cleanup remove ./output error = %v", err)
				return
			}
		} else {
			if !os.IsNotExist(err) {
				t.Errorf("stat directory error = %v", err)
				return
			}
		}
	}
	cleanupGitHubDir := func(t *testing.T) {
		if _, err := os.Stat(".github"); err == nil {
			err = os.RemoveAll(".github")
			if err != nil {
				t.Errorf("cleanup remove ./.github error = %v", err)
				return
			}
		} else {
			if !os.IsNotExist(err) {
				t.Errorf("stat directory error = %v", err)
				return
			}
		}
	}

	tests := []struct {
		name          string
		args          []string
		customizeMock func(mockBCli *switchboard.MockBlueskyClient, mockXCli *switchboard.MockXClient)
		wantErr       bool
		cleanup       func(t *testing.T)
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
						},
						{
							Cid:       "test2test2test2test2test2test2test2test2test2test2test2test2",
							Content:   "test2",
							CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
							URL:       "https://bsky.app/profile/did:plc:test2test2test2test2/post/test2test2",
						},
					}, nil)
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test1.*")).
					Return(&switchboard.XPost{
						ID: "1111111111111111111",
					}, nil)
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test2.*")).
					Return(&switchboard.XPost{
						ID: "2222222222222222222",
					}, nil)
			},
			wantErr: false,
			cleanup: cleanupOutputDir,
		},
		{
			name: "If bluesky2x subcommand got forbidden error from X API, warn the error but continue",
			args: []string{"switchboard", "bluesky2x"},
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
					Return(nil, errors.New("403 Forbidden error from X"))
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test2.*")).
					Return(&switchboard.XPost{
						ID: "2222222222222222222",
					}, nil)
			},
			wantErr: false,
			cleanup: cleanupOutputDir,
		},
		{
			name: "If bluesky2x subcommand got forbidden (duplicate post) error from X API, warn the error but continue",
			args: []string{"switchboard", "bluesky2x"},
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
					}, nil)
				mockXCli.EXPECT().Post(gomock.Any(), gomock.Regex("test1.*")).
					Return(nil, &switchboard.ErrXDuplicatePost{
						GoTwiError: &gotwi.GotwiError{
							Non2XXError: resources.Non2XXError{
								StatusCode: 403,
								Title:      "Forbidden",
								Detail:     "You are not allowed to create a Tweet with duplicate content",
							},
						},
					})
			},
			wantErr: false,
			cleanup: cleanupOutputDir,
		},
		{
			name: "If bluesky2x subcommand got bluesky post longer than approx 280, truncate post",
			args: []string{"switchboard", "bluesky2x"},
			customizeMock: func(mockBCli *switchboard.MockBlueskyClient, mockXCli *switchboard.MockXClient) {
				test1URL := "https://bsky.app/profile/did:plc:test1test1test1test1/post/test1test1"
				test2URL := "https://bsky.app/profile/did:plc:test2test2test2test2/post/test2test2"
				mockBCli.EXPECT().GetMyLatestPostsCreatedAsc(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]switchboard.BlueskyPost{
						{
							Cid:       "test1test1test1test1test1test1test1test1test1test1test1test1",
							Content:   strings.Repeat("x", 300),
							CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							URL:       test1URL,
						},
						{
							Cid:       "test2test2test2test2test2test2test2test2test2test2test2test2",
							Content:   strings.Repeat("あ", 150),
							CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
							URL:       test2URL,
							Reply:     nil,
						},
					}, nil)
				truncatedText1 := strings.Repeat("x", 202) + "...\n🤖from🦋:" + test1URL
				// 280 - 40(offset) - 34 (suffixLength) - 3 (ellipsis) = 202 / 2(CJK) = 101
				truncatedText2 := strings.Repeat("あ", 101) + "...\n🤖from🦋:" + test2URL
				gomock.InOrder(
					mockXCli.EXPECT().Post(gomock.Any(), truncatedText1).
						Return(&switchboard.XPost{
							ID: "1111111111111111111",
						}, nil),
					mockXCli.EXPECT().Post(gomock.Any(), truncatedText2).
						Return(&switchboard.XPost{
							ID: "2222222222222222222",
						}, nil),
				)
			},
			wantErr: false,
			cleanup: cleanupOutputDir,
		},
		{
			name:    "If bluesky2x --gen-workflow-file subcommand used, generate workflow files",
			args:    []string{"switchboard", "bluesky2x", "--gen-workflow-file"},
			wantErr: false,
			cleanup: cleanupGitHubDir,
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

			app, goterr := NewApp(&switchboardRequirementsForTest{
				ctx:  context.Background(),
				bcli: mockBCli,
				xcli: mockXCli,
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

			if tt.cleanup != nil {
				tt.cleanup(t)
			}
		})
	}
}
