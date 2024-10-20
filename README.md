# switchboard

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/switchboard)](http://pkg.go.dev/github.com/go-zen-chu/switchboard)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/main/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/check-pr/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/switchboard.svg)](https://github.com/go-zen-chu/switchboard/issues)

switchboard operator between sns.

## Usecase

1. Sync bluesky post to x

    ```console
    switchboard bluesky2x
    ```

2. Reply to bluesky post via openai

    ```console
    switchboard ai-reply -sns=bluesky
    ```
