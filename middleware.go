package main

import (
	"context"
	"errors"

	"github.com/rushyn/gator/internal/database"
)


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error{
	return func(s *state, cmd command) error {
		user, err := s.db.CheckUser(context.Background(), s.config.Current_User_Name)
		if err != nil {
			return errors.New("loged in user not found in database")
		}
		return handler(s, cmd, user)
	}
}


