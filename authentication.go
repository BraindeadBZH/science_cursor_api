package api

import (
	"errors"
	"log"
	"net/http"
)

func authenticationRoutes() {
	router.HandleFunc("/login", handlerWithSession(login)).Methods("POST")
	router.HandleFunc("/logout", handlerWithSession(logout)).Methods("POST")
}

func login(req *http.Request, session *sessionData) (interface{}, error) {
	log.Println("Call - login")
	if !session.Auth {
		var form loginForm
		valid := retrieveForm(&form, req, loginFormValidation)
		if valid.HasErrors() {
			log.Println("Login form is invalid:", valid.Error())
			return nil, errors.New("Invalid login form")
		}

		user, _ := getUserByEmail(form.Email)
		if user == nil {
			log.Println("Unable to find user:", form.Email)
			return nil, errors.New("Unable to find user")
		}

		err := user.authenticate(form.Password)
		if err != nil {
			log.Println("Password check failed")
			return nil, errors.New("Password check failed")
		}

		session.Auth = true
		session.Email = user.Email
		session.Name = user.Name
		session.Roles = user.Roles

		return nil, nil
	}

	return nil, errors.New("Already logged in")
}

func logout(req *http.Request, session *sessionData) (interface{}, error) {
	log.Println("Call - logout")
	session.MarkedForDeletion = true
	return nil, nil
}
