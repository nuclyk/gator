package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when deleting users: %v\n", err)
		os.Exit(1)
	}
	return nil
}
