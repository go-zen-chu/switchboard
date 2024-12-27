//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"strings"
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

func TestTruncateTweet(t *testing.T) {
	type args struct {
		content      string
		suffixLength int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "If only ascii characters with less than maxLength, return original text",
			args: args{
				content:      "This is test text 1234567890-=~+*",
				suffixLength: 34,
			},
			want: "This is test text 1234567890-=~+*",
		},
		{
			name: "If emoji and CJK characters with less than maxLength, return original text",
			args: args{
				content:      "こんにちは 你好 안녕하세요😊💕🕖",
				suffixLength: 34,
			},
			want: "こんにちは 你好 안녕하세요😊💕🕖",
		},
		{
			name: "If only ascii characters with more than maxLength, return truncated text",
			args: args{
				content:      strings.Repeat("x", 300),
				suffixLength: 34,
			},
			want: strings.Repeat("x", 243) + "...",
		},
		// 要修正。日本語だと 140 - 34 - 3 = 103 の長さになる
		{
			name: "If emoji and CJK characters with more than maxLength, return truncated text",
			args: args{
				content:      strings.Repeat("あ", 150),
				suffixLength: 34,
			},
			want: strings.Repeat("あ", 103) + "...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateTweet(tt.args.content, tt.args.suffixLength); got != tt.want {
				t.Errorf("TruncateTweet() = %v, want %v", got, tt.want)
			}
		})
	}
}
