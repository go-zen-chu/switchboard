# switchboard

> [!WARNING]
> Under development

[![Documentation](https://pkg.go.dev/badge/github.com/go-zen-chu/switchboard)](http://pkg.go.dev/github.com/go-zen-chu/switchboard)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/main/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![Actions Status](https://github.com/go-zen-chu/switchboard/workflows/check-pr/badge.svg)](https://github.com/go-zen-chu/switchboard/actions)
[![GitHub issues](https://img.shields.io/github/issues/go-zen-chu/switchboard.svg)](https://github.com/go-zen-chu/switchboard/issues)

Switchboard operator between sns.

## Use cases

### GitHub Actions

#### Automatically sync bluesky latest posts to X

You can do this by forking [go\-zen\-chu/bluesky2x\-workflow](https://github.com/go-zen-chu/bluesky2x-workflow). Please check its README.

### Running locally

#### Sync bluesky latest posts to X

```console
switchboard bluesky2x
# with verbose
switchboard bluesky2x -v
```

## FAQ and feature list to be implemented

[Please check issues](https://github.com/go-zen-chu/switchboard/labels/enhancement)

## Development

Using mage to make development easy

```console
# update release tag in main
git commit -am "commit something"
go run mage.go gitPushTag "release comment here"
```
