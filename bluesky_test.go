//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package switchboard

import (
	"testing"

	"github.com/bluesky-social/indigo/api/bsky"
)

func Test_replaceAbbreviatedURLToOriginal(t *testing.T) {
	tests := []struct {
		name     string
		feedPost *bsky.FeedPost
		want     string
	}{
		{
			name: "If facets is nil, return original text",
			feedPost: &bsky.FeedPost{
				Text: "This is test text",
			},
			want: "This is test text",
		},
		{
			name: "If URL is in facets of the post, replace abbreviated URL to original URL",
			feedPost: &bsky.FeedPost{
				Text: "This is test text\nURL: github.com/go-zen-chu/test?test...\nsome text follows...\ntest test",
				Facets: []*bsky.RichtextFacet{
					{
						Features: []*bsky.RichtextFacet_Features_Elem{
							{
								RichtextFacet_Link: &bsky.RichtextFacet_Link{
									Uri: "https://github.com/go-zen-chu/test?test=test&test1=test1",
								},
							},
						},
					},
				},
			},
			want: "This is test text\nURL: https://github.com/go-zen-chu/test?test=test&test1=test1\nsome text follows...\ntest test",
		},
		{
			name: "If URL is in facets of the post post, replace abbreviated URL to original URL in order",
			feedPost: &bsky.FeedPost{
				Text: "This is test text\nURL: github.com/go-... github.com/g... github.com/go-zen... other.com/test...",
				Facets: []*bsky.RichtextFacet{
					{
						Features: []*bsky.RichtextFacet_Features_Elem{
							{
								RichtextFacet_Link: &bsky.RichtextFacet_Link{
									Uri: "https://github.com/go-zen-chu/test?test=test&test1=test1",
								},
							},
							{
								RichtextFacet_Link: &bsky.RichtextFacet_Link{
									Uri: "https://github.com/go-zen-chu/test2",
								},
							},
							{
								RichtextFacet_Link: &bsky.RichtextFacet_Link{
									Uri: "https://github.com/go-zen-chu/test3",
								},
							},
						},
					},
					{
						Features: []*bsky.RichtextFacet_Features_Elem{
							{
								RichtextFacet_Link: &bsky.RichtextFacet_Link{
									Uri: "https://other.com/test",
								},
							},
						},
					},
				},
			},
			want: "This is test text\nURL: https://github.com/go-zen-chu/test?test=test&test1=test1 https://github.com/go-zen-chu/test2 https://github.com/go-zen-chu/test3 https://other.com/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceAbbreviatedURLToOriginal(tt.feedPost); got != tt.want {
				t.Errorf("replaceAbbreviatedURLToOriginal() = %v, want %v", got, tt.want)
			}
		})
	}
}
