package main

import (
	"context"
	"fmt"
)

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error when creating feed: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error when fetching the user with the feeds UserID: %v", err)
		}
		fmt.Printf("%v\n%v\n%v\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}
