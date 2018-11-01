package api

import (
	"fmt"
	"github.com/go-pg/pg/orm"
)

func createTables() error {
	tx, err := dbHandle.Begin()
	if err != nil {
		return fmt.Errorf("Could not start table creation transaction: '%s'", err.Error())
	}
	defer tx.Rollback()
	err = dbHandle.CreateTable(interface{}((*User)(nil)), &orm.CreateTableOptions{IfNotExists: true})
	if err != nil {
		return fmt.Errorf("Error while creating User table: '%s'", err.Error())
	}

	// Admin Account Creation //
	exists, err := userExists(apiConfig.General.AdminEmail)
	if err != nil {
		return fmt.Errorf("Error check admin existance: '%s'", err.Error())
	}
	if !exists {
		err = createUser(apiConfig.General.AdminEmail, "Administrator", []string{"admin"}, apiConfig.General.AdminPassword)
		if err != nil {
			return fmt.Errorf("Could not create admin: '%s'", err.Error())
		}
	}

	return tx.Commit()
}
