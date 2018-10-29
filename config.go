package api

import (
	"encoding/json"
	"fmt"
	"os"
)

// DatabaseConfiguration stores the configuration for the database access
type DatabaseConfiguration struct {
	Server   string
	Database string
	User     string
	Password string
}

// Configuration stores the configuration for the API
type Configuration struct {
	Database DatabaseConfiguration
}

var apiConfig Configuration

func loadConfiguration(configFile string) error {
	file, err := os.Open(configFile)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Unable to open the config file: '%s'", configFile)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&apiConfig)
	if err != nil {
		return fmt.Errorf("Unable to read the config file: '%s'", err.Error())
	}

	return nil
}
