package api

import (
	"encoding/json"
	"fmt"
	"os"
)

type generalConfiguration struct {
	Production    bool
	AdminEmail    string
	AdminPassword string
}

type databaseConfiguration struct {
	Server   string
	Database string
	User     string
	Password string
}

type httpConfiguration struct {
	Address    string
	Port       uint16
	SessionKey []byte
}

type configuration struct {
	General  generalConfiguration
	Database databaseConfiguration
	HTTP     httpConfiguration
}

var apiConfig configuration

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
