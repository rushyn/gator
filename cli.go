package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rushyn/gator/internal/config"
	"github.com/rushyn/gator/internal/database"
)


type state struct{
	db  *database.Queries
	config *config.Config
}

type command struct {
	argumets []string
}

type commands struct {
	commands map[string]func(*state, command) error 
}

func (c *commands) register(name string, f func(*state, command) error){
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error{
	_, ok := c.commands[cmd.argumets[1]]
	if !ok{
		return errors.New("invalid command " + cmd.argumets[1])
	}

	return c.commands[cmd.argumets[1]](s, cmd)
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.argumets) < 3 {
		return errors.New("no username given")
	}

	ctx := context.Background()
	user, err := s.db.CheckUser(ctx, cmd.argumets[2])
	if err != nil {
		return errors.New("user douse not exists")
	}

	s.config.SetUser(user.Name)
	fmt.Printf("The loged in use is :>%s<: !!!\n", cmd.argumets[2])
	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.argumets) < 3 {
		return errors.New("no username given")
	}
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	ctx := context.Background()
	newUser := database.CreateUserParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.argumets[2],
	}
	user, err :=s.db.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}
	fmt.Printf("New user %s was created!!!\n", user.Name)
	fmt.Println(user.Name)
	fmt.Println(user.CreatedAt)
	fmt.Println(user.UpdatedAt)
	fmt.Println(user.ID)
	s.config.SetUser(user.Name)

	return nil
}


func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("All users deleted !!!")
	return nil
}


func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.ShowAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users{
		if user.Name == s.config.Current_User_Name {
			fmt.Printf("%s (current)\n", user.Name)
		}else{
			fmt.Printf("%s\n", user.Name)
		}

	}
	return nil
}