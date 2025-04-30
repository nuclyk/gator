package main

import (
	"context"
	"github.com/nuclyk/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	s.db.Unfollow(context.Background(), params)
	return nil
}
