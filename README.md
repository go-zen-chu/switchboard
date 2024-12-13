# switchboard

> [!WARNING]
> Under development

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/switchboard)](http://pkg.go.dev/github.com/go-zen-chu/switchboard)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/main/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/check-pr/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/switchboard.svg)](https://github.com/go-zen-chu/switchboard/issues)

switchboard operator between sns.

## Usecase

### Sync bluesky post to x

```console
switchboard bluesky2x
# with --ai option, genai will response to your post via aictl
switchboard bluesky2x --ai 
```

### Usage

1. Post to bluesky then the post will be posted to X
2. Post to bluesky with /ai in head of post, ai will respond to your post (via aictl)
3. When someone send reply to your bluesky post, it will do nothing for privacy
4. When you delete your bluesky post, it will deleted from X too
