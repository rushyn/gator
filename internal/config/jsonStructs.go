package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DB_URL string `json:"db_url"`
	Current_User_Name string `json:"current_user_name"`
}

func (c *Config)SetUser(user string) {
	c.Current_User_Name = user

	path, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Unable to get config file path, error is ::: %s \n", err)
	}
	
	config, err := json.Marshal(&c)
	if err != nil {
		log.Panicf("Unable to encode json ::: %s \n", err)
	}

	os.WriteFile(path + "/.gatorconfig.json", config, 0666)
	
}