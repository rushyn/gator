package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

func ReadGarorConfig () Config {

	path, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Unable to get config file path, error is ::: %s \n", err)
	}
	
	configFile, err := os.ReadFile(path + "/.gatorconfig.json")
	if err != nil {
		log.Panicf("Unable to Read configfile .gatorconfig.json ::: %s \n", err)
	}
	
	ioReader := bytes.NewReader(configFile)

	Config := Config{}
	decoder := json.NewDecoder(ioReader)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Panicf("Unable unable to decode json ::: %s \n", err)
	}
	return Config
}