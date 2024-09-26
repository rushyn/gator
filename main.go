package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/rushyn/gator/internal/config"
	"github.com/rushyn/gator/internal/database"
)


func main(){
	conf := config.ReadGarorConfig()

	dbc, err := sql.Open("postgres", conf.DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(dbc)

	cmd := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmd.register("login", handlerLogin)
	cmd.register("register", handlerRegister)
	cmd.register("reset", handlerReset)
	cmd.register("users", handlerGetUsers)
	cmd.register("agg", handlerAgg)
	cmd.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmd.register("feeds", handlerFeeds)
	cmd.register("follow", middlewareLoggedIn(handlerFollow))
	cmd.register("following", handlerFollowing)
	cmd.register("unfollow", handlerUnfollow)



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
		db: db,
	}






	err = cmd.run(&state, arguments)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}


}


