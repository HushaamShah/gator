package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/HushaamShah/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command requires a username")
	}
	_, err := s.queries.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not Found")
		}
		return err
	}

	s.dbconfig.SetUser(cmd.args[0])
	fmt.Println("User has been set.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	_, err := s.queries.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}
	s.dbconfig.SetUser(cmd.args[0])
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.dbconfig.User {
			fmt.Printf("* %s (current) \n", user)
			continue
		}
		fmt.Printf("* %s \n", user)
	}
	return nil
}
