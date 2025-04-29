package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nuclyk/gator/internal/database"
)

// TODO: Add checking if the second argument is a link.
// If I type in Google Blog https://google.com then the link will be Blog
// instead of the actual link ad the Blog is the second argument. FIX IT!
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		fmt.Fprintf(os.Stderr, "You need to provide name and url.")
		os.Exit(1)
	} else if len(cmd.args) < 2 {
		fmt.Fprintf(os.Stderr, "You need to provide url.")
		os.Exit(1)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error when creating feed: %w", err)
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("Error when adding feed follow for the current user: %w", err)
	}

	fmt.Printf("ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nUrl: %v\nUser ID: %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}
