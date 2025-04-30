package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/nuclyk/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	args := cmd.args
	if len(args) < 1 {
		return fmt.Errorf("You need to provide url argument")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when getting a feed with given url: %v\n", err)
		os.Exit(1)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follows, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when creating feed follow: %v", err)
		os.Exit(1)
	}

	fmt.Println(follows[0].FeedName)
	fmt.Println(follows[0].UserName)
	return nil
}
