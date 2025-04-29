package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
)

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
