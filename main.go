package main

import (
	"fmt"

	"github.com/HushaamShah/gator/internal/config"
)

type state struct {
	dbconfig config.Config
}

type command struct {
	name string
	args []string
}

func main() {
	dbconfig, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	dbconfig.SetUser("Hushaam")
	config.Read()
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command requires a username")
	}
	s.dbconfig.SetUser(cmd.args[0])
	fmt.Println("User has been set.")
	return nil
}
