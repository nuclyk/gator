package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
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

func initCommands() commands {
	cmds := commands{
		mapCommand: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	return cmds
}

func initConfig() config.Config {
	config, err := config.Read(".gatorconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func initDBandState(cfg config.Config) state {
	db, err := sql.Open("postgres", cfg.DB_URL)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal(err)
	}

	currentState := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	return currentState
}

func StartCli() {
	args := os.Args
	commandName := args[1]
	commandArgs := args[2:]
	command := command{commandName, commandArgs}

	config := initConfig()
	state := initDBandState(config)
	commands := initCommands()

	err := commands.run(&state, command)
	if err != nil {
		fmt.Printf("Error when running the command:\n%v\n", err)
		os.Exit(1)
	}
}

func (cmds commands) run(s *state, cmd command) error {
	if _, ok := cmds.mapCommand[cmd.name]; ok {
		cmds.mapCommand[cmd.name](s, cmd)
	} else {
		return fmt.Errorf("Command %s not found in the register", cmd.name)
	}
	return nil
}

func (cs commands) register(name string, f func(*state, command) error) {
	cs.mapCommand[name] = f
}
