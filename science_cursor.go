package api

import (
	"fmt"
	"github.com/go-pg/pg"
)

// DbHandle gives access to the database
var dbHandle *pg.DB

// Run starts the API using the given config file
func Run(configFile string) error {
	fmt.Println("Starting API using config file:", configFile)

	err := loadConfiguration(configFile)
	if err != nil {
		return fmt.Errorf("Unable to load configuration: '%s'", err.Error())
	}

	dbHandle = pg.Connect(
		&pg.Options{Addr: apiConfig.Database.Server,
			Database: apiConfig.Database.Database,
			User:     apiConfig.Database.User,
			Password: apiConfig.Database.Password})

	err = createTables()
	if err != nil {
		return fmt.Errorf("Error while creating tables: '%s'", err.Error())
	}

	go sessionCleanupRoutine()

	authenticationRoutes()

	httpListen()

	return nil
}
