package main

import (
	"context"
	"fmt"
)

func handlerFollowing(s *state, cmd command) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	for _, feed := range follows {
		fmt.Printf("%s\n%s\n", feed.FeedsName, feed.UserName)
	}
	return nil
}
