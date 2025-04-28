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

func main() {
	// Read config file
	config, err := config.Read(".gatorconfig.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create connection with db
	db, err := sql.Open("postgres", config.DB_URL)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal(err)
	}

	// Set up new state
	current_state := state{
		db:  dbQueries,
		cfg: &config,
	}

	// Initialize commands
	cs := commands{
		mapCommand: map[string]func(*state, command) error{},
	}
	cs.register("login", handlerLogin)
	cs.register("register", handlerRegister)
	cs.register("reset", handlerReset)
	cs.register("users", handlerGetUsers)
	cs.register("agg", handlerAgg)
	cs.register("addfeed", handlerAddFeed)
	cs.register("feeds", handlerGetFeeds)
	cs.register("follow", handlerFollow)
	cs.register("following", handlerFollowing)

	args := os.Args
	// if len(args) < 2 {
	// 	log.Fatal("Program needs more arguments")
	// 	os.Exit(1)
	// }

	commandName := args[1]
	commandArgs := args[2:]
	command := command{commandName, commandArgs}

	// Run the command from cli
	err = cs.run(&current_state, command)
	if err != nil {
		fmt.Printf("Error when running the command:\n%v\n", err)
		os.Exit(1)
	}

}
