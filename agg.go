package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedUrl)

	if err != nil {
		return fmt.Errorf("Error when fetching rss feed: %w", err)
	}

	fmt.Print(rssFeed)
	return nil
}
