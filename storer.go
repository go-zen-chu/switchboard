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
	// DefaultMaxHistorySizeBytes is the default maximum size of sync history in bytes (500KB)
	DefaultMaxHistorySizeBytes = 500 * 1024
)

type Storer struct {
	SyncInfo           *SyncInfo
	MaxHistorySizeBytes int
}

func NewStorer() *Storer {
	return &Storer{
		SyncInfo: &SyncInfo{
			Posts: make([]PostInfo, 0),
		},
		MaxHistorySizeBytes: DefaultMaxHistorySizeBytes,
	}
}

// NewStorerWithMaxSize creates a new Storer with custom max history size.
// Setting maxSizeBytes to 0 means unlimited.
func NewStorerWithMaxSize(maxSizeBytes int) *Storer {
	return &Storer{
		SyncInfo: &SyncInfo{
			Posts: make([]PostInfo, 0),
		},
		MaxHistorySizeBytes: maxSizeBytes,
	}
}

func (s *Storer) StoreSyncInfo() error {
	if s.SyncInfo == nil {
		return fmt.Errorf("sync info is nil")
	}
	
	// Trim history if needed before storing
	if err := s.trimHistoryIfNeeded(); err != nil {
		return fmt.Errorf("trimming history: %w", err)
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

// trimHistoryIfNeeded trims the sync history if it exceeds the maximum size limit.
// If MaxHistorySizeBytes is 0, no trimming is performed (unlimited).
// Removes oldest posts first to stay within the size limit.
func (s *Storer) trimHistoryIfNeeded() error {
	// If unlimited (0), skip trimming
	if s.MaxHistorySizeBytes == 0 {
		return nil
	}
	
	// Calculate current size
	currentSize, err := s.calculateSyncInfoSize()
	if err != nil {
		return fmt.Errorf("calculating sync info size: %w", err)
	}
	
	// If under the limit, no trimming needed
	if currentSize <= s.MaxHistorySizeBytes {
		return nil
	}
	
	// Remove oldest posts until we're under the limit
	// Posts are stored in chronological order (oldest first based on the sync workflow)
	for currentSize > s.MaxHistorySizeBytes && len(s.SyncInfo.Posts) > 0 {
		// Remove the oldest post (first element)
		s.SyncInfo.Posts = s.SyncInfo.Posts[1:]
		
		// Recalculate size
		currentSize, err = s.calculateSyncInfoSize()
		if err != nil {
			return fmt.Errorf("recalculating sync info size: %w", err)
		}
	}
	
	return nil
}

// calculateSyncInfoSize returns the approximate size in bytes of the SyncInfo when serialized to JSON
func (s *Storer) calculateSyncInfoSize() (int, error) {
	data, err := json.Marshal(s.SyncInfo)
	if err != nil {
		return 0, fmt.Errorf("marshaling sync info: %w", err)
	}
	return len(data), nil
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
