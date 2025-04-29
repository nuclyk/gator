package main

import (
	"context"
	"fmt"
	"os"
)

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
