package main

import (
	"errors"
	"fmt"

	"github.com/rushyn/gator/internal/config"
)


type state struct{
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
	s.config.SetUser(cmd.argumets[2])
	fmt.Printf("The loged in use is :>%s<: !!!\n", cmd.argumets[2])
	return nil
}