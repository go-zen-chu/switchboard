import os
from atproto import Client as BlueskyClient
from atproto import AtUri
import tweepy


class BlueskyClient():
    def __init__(self, identifier: str, password: str):
        self.client = BlueskyClient()
        self.client.login(identifier, password)
    
    def get_latest_post(self, num_posts=3):
        
    bluesky = BlueskyClient()
    bluesky.login(os.environ['BLUESKY_IDENTIFIER'], os.environ['BLUESKY_PASSWORD'])


class XClient():
    def __init__(self, access_token: str, access_secret: str, api_key: str, api_secret: str, bearer_token: str):
        self.client = tweepy.Client(
            consumer_key=api_key,
            consumer_secret=api_secret,
            access_token=access_token,
            access_token_secret=access_secret
        )
if __name__ == "__main__":
    bcli = BlueskyClient(os.environ['BLUESKY_IDENTIFIER'], os.environ['BLUESKY_PASSWORD'])


client = tweepy.Client(
    consumer_key=os.environ['X_API_KEY'],
    consumer_secret=os.environ['X_API_SECRET'],
    access_token=os.environ['X_ACCESS_TOKEN'],
    access_token_secret=os.environ['X_ACCESS_SECRET']
)

profile = bluesky.get_profile(os.environ['BLUESKY_IDENTIFIER'])
posts = bluesky.get_author_feed(profile.did)

if posts.feed:
    latest_post = posts.feed[0]
    post_text = latest_post.post.record.text

    parsed_uri = AtUri.from_str(latest_post.post.uri)
    post_uri = f"https://bsky.app/profile/{parsed_uri.hostname}/post/{parsed_uri.rkey}"

    txt = f"{post_text}\n\n🤖from🦋{post_uri}"

    client.create_tweet(text=txt)
    print(f"Posted to X:\n{txt}")
else:
    print("No recent Bluesky posts found")

