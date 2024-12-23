//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"testing"

	"github.com/bluesky-social/indigo/api/bsky"
)

func Test_replaceAbbreviatedURLToEmbedExternal(t *testing.T) {
	tests := []struct {
		name     string
		feedPost *bsky.FeedPost
		want     string
	}{
		{
			name: "If embedExternal is nil, return original text",
			feedPost: &bsky.FeedPost{
				Text: "This is test text",
			},
			want: "This is test text",
		},
		{
			name: "If URL is embeded in post, replace abbreviated URL to original URL",
			feedPost: &bsky.FeedPost{
				Text: "This is test text\nURL: github.com/go-zen-chu/test?test...\nsome text follows...\ntest test",
				Embed: &bsky.FeedPost_Embed{
					EmbedExternal: &bsky.EmbedExternal{
						External: &bsky.EmbedExternal_External{
							Uri: "https://github.com/go-zen-chu/test?test=test&test1=test1",
						},
					},
				},
			},
			want: "This is test text\nURL: https://github.com/go-zen-chu/test?test=test&test1=test1\nsome text follows...\ntest test",
		},
		{
			name: "If URL is embeded in post, replace abbreviated URL to original URL (case2)",
			feedPost: &bsky.FeedPost{
				Text: "This is test text\nURL: github.com/go-... some text follows... git...",
				Embed: &bsky.FeedPost_Embed{
					EmbedExternal: &bsky.EmbedExternal{
						External: &bsky.EmbedExternal_External{
							Uri: "https://github.com/go-zen-chu/test?test=test&test1=test1",
						},
					},
				},
			},
			want: "This is test text\nURL: https://github.com/go-zen-chu/test?test=test&test1=test1 some text follows... git...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceAbbreviatedURLToEmbedExternal(tt.feedPost); got != tt.want {
				t.Errorf("replaceAbbreviatedURLToEmbedExternal() = %v, want %v", got, tt.want)
			}
		})
	}
}
