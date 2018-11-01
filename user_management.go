package api

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func userExists(email string) (bool, error) {
	count, err := dbHandle.Model((*User)(nil)).Where("email = ?", email).Count()
	if err != nil {
		return false, fmt.Errorf("Could not check user existance: '%s'", err.Error())
	}
	return count > 0, nil
}

func getUserByEmail(email string) (*User, error) {
	user := &User{}
	err := dbHandle.Model((*User)(nil)).Where("email = ?", email).Select(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func createUser(email string, name string, roles []string, password string) error {
	user := &User{}
	user.Email = email
	user.Name = name
	user.Roles = roles

	user.Salt = make([]byte, 32)
	_, err := rand.Read(user.Salt)
	if err != nil {
		return fmt.Errorf("Error while creating random salt: '%s'", err.Error())
	}

	user.Password, err = bcrypt.GenerateFromPassword(append([]byte(apiConfig.General.AdminPassword), user.Salt...), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error while hashing password: '%s'", err.Error())
	}

	err = dbHandle.Insert(user)
	if err != nil {
		return fmt.Errorf("Error while creating user: '%s'", err.Error())
	}
	return nil
}
