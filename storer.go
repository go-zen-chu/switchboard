package switchboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	syncInfoPath = "./output/sync_info.json"
)

type Storer struct {
	SyncInfo *SyncInfo
}

func NewStorer() *Storer {
	return &Storer{
		SyncInfo: &SyncInfo{
			Posts: make([]PostInfo, 0),
		},
	}
}

func (s *Storer) StoreSyncInfo() error {
	if s.SyncInfo == nil {
		return fmt.Errorf("sync info is nil")
	}
	dirPath := filepath.Dir(syncInfoPath)
	if _, err := os.Stat(dirPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat directory: %w", err)
		}
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
	}
	f, err := os.Create(syncInfoPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(s.SyncInfo)
	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}

// LoadSyncInfo loads sync info from file. If file does not exist, return empty SyncInfo
func (s *Storer) LoadSyncInfo() (*SyncInfo, error) {
	if _, err := os.Stat(syncInfoPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat file: %w", err)
		}
		// if file path does not exist, return empty SyncInfo
		return s.SyncInfo, nil
	}
	f, err := os.Open(syncInfoPath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(s.SyncInfo)
	if err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return s.SyncInfo, nil
}

type SyncInfo struct {
	Posts []PostInfo `json:"posts"`
}

type PostInfo struct {
	BlueskyCid           string    `json:"bluesky_cid"`
	BlueskyPostCreatedAt time.Time `json:"bluesky_post_created_at"`
	TweetID              string    `json:"tweet_id"`
	Content              string    `json:"content"`
}

func (s *SyncInfo) GetPostMap() map[string]PostInfo {
	m := make(map[string]PostInfo)
	for _, p := range s.Posts {
		m[p.BlueskyCid] = p
	}
	return m
}
