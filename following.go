package main

import (
	"context"
	"fmt"

	"github.com/nuclyk/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range follows {
		fmt.Printf("%s\n%s\n", feed.FeedsName, feed.UserName)
	}
	return nil
}
