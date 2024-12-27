//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"testing"
)

func TestCountTweetCharacters(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{
			name:    "only ascii characters",
			content: "This is test text 1234567890-=~+*...",
			want:    36,
		},
		{
			name:    "emoji and CJK characters",
			content: "こんにちは 你好 안녕하세요😊💕🕖",
			want:    32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountTweetCharacters(tt.content); got != tt.want {
				t.Errorf("CountTweetCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}
