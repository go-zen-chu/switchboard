package switchboard

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStorer_TrimHistoryIfNeeded(t *testing.T) {
	tests := []struct {
		name                string
		maxHistorySizeBytes int
		numPosts            int
		wantErr             bool
		checkFunc           func(t *testing.T, stor *Storer, originalNumPosts int)
	}{
		{
			name:                "no trimming when under limit",
			maxHistorySizeBytes: 10 * 1024, // 10KB
			numPosts:            5,
			wantErr:             false,
			checkFunc: func(t *testing.T, stor *Storer, originalNumPosts int) {
				if len(stor.SyncInfo.Posts) != originalNumPosts {
					t.Errorf("Expected no trimming, got %d posts (original: %d)", len(stor.SyncInfo.Posts), originalNumPosts)
				}
			},
		},
		{
			name:                "trim when over limit",
			maxHistorySizeBytes: 1024, // 1KB
			numPosts:            100,
			wantErr:             false,
			checkFunc: func(t *testing.T, stor *Storer, originalNumPosts int) {
				if len(stor.SyncInfo.Posts) >= originalNumPosts {
					t.Errorf("Expected posts to be trimmed, but got %d posts (original: %d)", len(stor.SyncInfo.Posts), originalNumPosts)
				}
				if len(stor.SyncInfo.Posts) == 0 {
					t.Errorf("Expected some posts to remain after trimming, but got 0 posts")
				}
				// Verify we kept the most recent posts
				if len(stor.SyncInfo.Posts) > 0 {
					lastPost := stor.SyncInfo.Posts[len(stor.SyncInfo.Posts)-1]
					expectedLastCid := fmt.Sprintf("test-cid-%d", originalNumPosts-1)
					if lastPost.BlueskyCid != expectedLastCid {
						t.Errorf("Expected to keep most recent post with CID %s, got %s", expectedLastCid, lastPost.BlueskyCid)
					}
				}
			},
		},
		{
			name:                "no trimming when unlimited (0)",
			maxHistorySizeBytes: 0,
			numPosts:            100,
			wantErr:             false,
			checkFunc: func(t *testing.T, stor *Storer, originalNumPosts int) {
				if len(stor.SyncInfo.Posts) != originalNumPosts {
					t.Errorf("Expected no trimming with unlimited size, got %d posts (original: %d)", len(stor.SyncInfo.Posts), originalNumPosts)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stor := NewStorerWithMaxSize(tt.maxHistorySizeBytes)

			// Create test posts
			for i := 0; i < tt.numPosts; i++ {
				stor.SyncInfo.Posts = append(stor.SyncInfo.Posts, PostInfo{
					BlueskyCid:           fmt.Sprintf("test-cid-%d", i),
					TweetID:              fmt.Sprintf("test-tweet-%d", i),
					Content:              fmt.Sprintf("This is test content for post number %d", i),
					BlueskyPostCreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
				})
			}

			err := stor.trimHistoryIfNeeded()
			if (err != nil) != tt.wantErr {
				t.Errorf("trimHistoryIfNeeded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Run custom checks
			tt.checkFunc(t, stor, tt.numPosts)

			// Verify size is under limit after trimming (if not unlimited)
			if tt.maxHistorySizeBytes > 0 {
				size, err := stor.calculateSyncInfoSize()
				if err != nil {
					t.Errorf("calculateSyncInfoSize() error = %v", err)
					return
				}
				if size > tt.maxHistorySizeBytes {
					t.Errorf("Size after trimming (%d) exceeds max size (%d)", size, tt.maxHistorySizeBytes)
				}
			}
		})
	}
}

func TestStorer_CalculateSyncInfoSize(t *testing.T) {
	stor := NewStorer()

	// Empty posts should have minimal size
	emptySize, err := stor.calculateSyncInfoSize()
	if err != nil {
		t.Errorf("calculateSyncInfoSize() error = %v", err)
		return
	}

	// Add a post
	stor.SyncInfo.Posts = append(stor.SyncInfo.Posts, PostInfo{
		BlueskyCid:           "test-cid",
		TweetID:              "test-tweet",
		Content:              "Test content",
		BlueskyPostCreatedAt: time.Now(),
	})

	// Size should increase
	sizeWithPost, err := stor.calculateSyncInfoSize()
	if err != nil {
		t.Errorf("calculateSyncInfoSize() error = %v", err)
		return
	}

	if sizeWithPost <= emptySize {
		t.Errorf("Size with post (%d) should be greater than empty size (%d)", sizeWithPost, emptySize)
	}
}

func TestStorer_StoreSyncInfoWithTrimming(t *testing.T) {
	defer func() {
		// Clean up
		os.RemoveAll("./output")
	}()

	stor := NewStorerWithMaxSize(1024) // 1KB limit

	// Add many posts to exceed the limit
	for i := 0; i < 50; i++ {
		stor.SyncInfo.Posts = append(stor.SyncInfo.Posts, PostInfo{
			BlueskyCid:           fmt.Sprintf("test-cid-%d", i),
			TweetID:              fmt.Sprintf("test-tweet-%d", i),
			Content:              fmt.Sprintf("This is test content for post number %d with some additional text to increase size", i),
			BlueskyPostCreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
		})
	}

	// Store should trim automatically
	err := stor.StoreSyncInfo()
	if err != nil {
		t.Errorf("StoreSyncInfo() error = %v", err)
		return
	}

	// Verify the size is under the limit
	size, err := stor.calculateSyncInfoSize()
	if err != nil {
		t.Errorf("calculateSyncInfoSize() error = %v", err)
		return
	}

	if size > 1024 {
		t.Errorf("Stored sync info size (%d) exceeds limit (1024)", size)
	}

	// Verify fewer posts remain
	if len(stor.SyncInfo.Posts) >= 50 {
		t.Errorf("Expected posts to be trimmed, but got %d posts", len(stor.SyncInfo.Posts))
	}
}

func TestNewStorerWithMaxSize(t *testing.T) {
	tests := []struct {
		name         string
		maxSizeBytes int
	}{
		{
			name:         "default size",
			maxSizeBytes: DefaultMaxHistorySizeBytes,
		},
		{
			name:         "custom size",
			maxSizeBytes: 1024,
		},
		{
			name:         "unlimited",
			maxSizeBytes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stor := NewStorerWithMaxSize(tt.maxSizeBytes)
			if stor.MaxHistorySizeBytes != tt.maxSizeBytes {
				t.Errorf("Expected MaxHistorySizeBytes to be %d, got %d", tt.maxSizeBytes, stor.MaxHistorySizeBytes)
			}
		})
	}
}

func TestNewStorer_DefaultMaxSize(t *testing.T) {
	stor := NewStorer()
	if stor.MaxHistorySizeBytes != DefaultMaxHistorySizeBytes {
		t.Errorf("Expected default MaxHistorySizeBytes to be %d, got %d", DefaultMaxHistorySizeBytes, stor.MaxHistorySizeBytes)
	}
}
