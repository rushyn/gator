package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rushyn/gator/internal/config"
)


func main(){
	conf := config.ReadGarorConfig()

	cmd := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmd.register("login", handlerLogin)

	arguments := command{
		argumets: os.Args,
	} 

	if len(arguments.argumets) < 2 {
		log.Println("No arguments supplied.")
		os.Exit(1)
		return
	}

	state := state{
		config: &conf,
	}

	err := cmd.run(&state, arguments)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}


}