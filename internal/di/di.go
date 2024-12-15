package di

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-zen-chu/switchboard"
)

type Container struct {
	cache map[string]any
}

func NewContainer() *Container {
	return &Container{
		cache: map[string]any{},
	}
}

func initOnce[T any](c *Container, component string, fn func() (T, error)) T {
	if v, ok := c.cache[component]; ok {
		return v.(T)
	}
	var err error
	v, err := fn()
	if err != nil {
		slog.Error("failed to set up "+component, "error", err)
		os.Exit(1)
	}
	c.cache[component] = v
	return v
}

func (c *Container) Context() context.Context {
	return initOnce(c, "Context", func() (context.Context, error) {
		return context.Background(), nil
	})
}

func (c *Container) BlueskyClient() switchboard.BlueskyClient {
	return initOnce(c, "BlueskyClient", func() (switchboard.BlueskyClient, error) {
		bcli, err := switchboard.NewBlueskyClient(
			c.Context(),
			os.Getenv("BLUESKY_IDENTIFIER"),
			os.Getenv("BLUESKY_PASSWORD"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set up BlueskyClient: %w", err)
		}
		return bcli, nil
	})
}

func (c *Container) XClient() switchboard.XClient {
	return initOnce(c, "XClient", func() (switchboard.XClient, error) {
		xcli, err := switchboard.NewXClient(
			c.Context(),
			os.Getenv("X_ACCESS_TOKEN"),
			os.Getenv("X_ACCESS_SECRET"),
			os.Getenv("X_API_KEY"),
			os.Getenv("X_API_SECRET"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set up XClient: %w", err)
		}
		return xcli, nil
	})
}
