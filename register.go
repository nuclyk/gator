package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nuclyk/gator/internal/database"
)

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
