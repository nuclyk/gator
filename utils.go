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

func handlerDeleteUsers(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("commands are empty")
	}

	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background())

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

	checkUser, err := s.db.GetUser(context.Background())
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
	if len(cmd.args) < 1 {
		return errors.New("You need to provide username as an argument.")
	}

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
