package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nuclyk/gator/internal/config"
	"github.com/nuclyk/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	mapCommand map[string]func(*state, command) error
}

func handlerFollowing(s *state, cmd command) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when fetching feed follows: %w", err)
		os.Exit(1)
	}

	for _, feed := range follows {
		fmt.Printf("%s\n%s\n", feed.FeedsName, feed.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	args := cmd.args
	if len(args) < 1 {
		return fmt.Errorf("You need to provide url argument")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when getting a user: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error when creating feed follow: %w", err)
		os.Exit(1)
	}

	fmt.Println(follows[0].FeedName)
	fmt.Println(follows[0].UserName)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when deleting users: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedUrl)

	if err != nil {
		return fmt.Errorf("Error when fetching rss feed: %w", err)
	}

	fmt.Print(rssFeed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	username := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when fetching the user: %v", err)
		os.Exit(1)

	}

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
		UserID:    currentUser.ID,
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
		UserID:    currentUser.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("Error when adding feed follow for the current user: %w", err)
	}

	fmt.Printf("ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nUrl: %v\nUser ID: %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)

	return nil
}

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

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when retreiving users: %v\n", err)
		os.Exit(1)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("commands are empty")
	}

	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)

	if err != nil || user.Name != username {
		log.Fatal("Can't login if the user doesn't exist.")
		os.Exit(1)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Error when setting up the user: %v", err)
	}

	fmt.Printf("User %s set successfuly :)", cmd.args)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	username := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	checkUser, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User with the name %s doesn't exist.\n", username)
	}

	if checkUser.Name == username {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	s.cfg.SetUser(username)
	user, err := s.db.CreateUser(context.Background(), params)

	if err != nil {
		return fmt.Errorf("Error when creating a user: %v\n", err)
	}

	fmt.Printf("New user was created in the database: %s", user.Name)

	return nil

}

func (cs commands) run(s *state, cmd command) error {
	if _, ok := cs.mapCommand[cmd.name]; ok {
		cs.mapCommand[cmd.name](s, cmd)
	} else {
		return fmt.Errorf("Command %s not found in the register", cmd.name)
	}
	return nil
}

func (cs commands) register(name string, f func(*state, command) error) {
	cs.mapCommand[name] = f
}
