package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.argumets) < 3 {
		return errors.New("no pooling inteval was givne, interval needs to be 1s, 1m, 1h")
	}

	if cmd.argumets[2] == "1s" || cmd.argumets[2] == "1m" || cmd.argumets[2] == "1h" {
	}else{
		return errors.New("invalid polling inteveal uspplied, should be 1s, 1m, 1h")
	}

	var timeBetweenRequests time.Duration

	if cmd.argumets[2] == "1s" {
		timeBetweenRequests = time.Duration(1 * time.Second)
	}else if cmd.argumets[2] == "1m" {
		timeBetweenRequests = time.Duration(1 * time.Minute)
	}else if cmd.argumets[2] == "1h"{
		timeBetweenRequests = time.Duration(1 * time.Hour)
	}

	ticket := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticket.C{
		scrapeFeeds(s)
	}
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.argumets) < 4 {
		return errors.New("no url given")
	}

	if len(cmd.argumets) < 3 {
		return errors.New("no feed name given")
	}

	ctx := context.Background()

	newUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	feedAdd := database.CreateFeedParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.argumets[2],
		Url:       cmd.argumets[3],
		UserID:    user.ID,
	}
	
	feed, err := s.db.CreateFeed(ctx, feedAdd)
	if err != nil {
		return errors.New("unable to create feed")
	}

	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	
	for _, feed := range rss.Channel.Item{
		fmt.Println(feed)
	}

	c := command{
		argumets: []string{"", "", cmd.argumets[3]},
	}

	return handlerFollow(s, c, user)

}


func handlerFeeds (s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.ReturnAllFeeds(ctx)
	if err != nil {
		return errors.New("unable to get feeds")
	}

	for _, feed := range feeds{
		fmt.Printf("Feed %s URL is %s and was created by %s.\n", feed.Feedname, feed.Url, feed.Username)
	}
	return nil
}

func handlerFollow (s *state, cmd command, user database.User) error {
	if len(cmd.argumets) < 3 {
		return errors.New("no url to follow was given")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.argumets[2])
	if err != nil {
		return errors.New("unable to get find feed")
	}
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	addFollow := database.CreateFeedFollowParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	
	feedRow, err := s.db.CreateFeedFollow(ctx, addFollow)
	if err != nil {
		return errors.New("unable to create feed")
	}
	fmt.Printf("User: %s has subscribed to Feed: %s !!!\n", feedRow.UserName, feedRow.FeedName)


	return nil
}



func handlerFollowing (s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, s.config.Current_User_Name)
	if err != nil {
		return errors.New("unable to get feeds")
	}

	fmt.Println("You are fallwoing this feeds.")
	for _, feed := range feeds{
		fmt.Printf("%s \n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow (s *state, cmd command) error{
	if len(cmd.argumets) < 3 {
		return errors.New("no feed url given")
	}
	dbU := database.UnfollowParams{
		Name: s.config.Current_User_Name,
		Url: cmd.argumets[2],
	}
	err := s.db.Unfollow(context.Background(), dbU)
	if err != nil {
		fmt.Println(err)
		return errors.New("unable to unfollow")
	}

	return nil
}


func handlerBrowse (s *state, cmd command) error{
	var limit int
	if len(cmd.argumets) < 3 {
		limit = 2
	}else{
		var err error
		limit, err = strconv.Atoi(cmd.argumets[2])
		if err != nil {
			return err
		}		
	}

	feed, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Name:  s.config.Current_User_Name,
		Limit: int32(limit),
	})
	if err != nil{
		return err
	}

	for i, item := range feed{
		fmt.Printf("\n this is the item number %v \n", i + 1)
		fmt.Printf("Post iD is %s \n", item.ID)
		fmt.Printf("This post was created at %s \n", item.CreatedAt)
		fmt.Printf("This post was update at %s \n", item.UpdatedAt)
		fmt.Printf("The time is --- %s \n", item.Title)
		fmt.Printf("The URL is %s \n", item.Url)
		fmt.Printf("\nDescription \n %s \n\n", item.Desription)
		fmt.Printf("It was published on %s \n", item.PublishedAt)
		fmt.Printf("And its part of feed %s \n", item.FeedID)
	}

	


	return nil
}