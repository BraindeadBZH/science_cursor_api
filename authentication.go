package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func authenticationRoutes() {
	router.HandleFunc("/login", handlerWithSession(login)).Methods("POST")
	router.HandleFunc("/logout", handlerWithSession(logout)).Methods("POST")
}

func login(req *http.Request, session *sessions.Session) (interface{}, error) {
	log.Println("Call - login")
	loggedIn, exists := session.Values["loggedIn"]
	if !exists || !loggedIn.(bool) {
		decoder := json.NewDecoder(req.Body)
		var form loginForm
		err := decoder.Decode(&form)
		if err != nil {
			log.Println("Unable to read form:", err.Error())
			return nil, errors.New("Unable to read login form")
		}

		valid := loginFormValidation.Validate(form)

		if valid.HasErrors() {
			log.Println("Login form is invalid:", valid.Error())
			return nil, errors.New("Invalid login form")
		}

		user, _ := getUserByEmail(form.Email)
		if user == nil {
			log.Println("Unable to find user:", form.Email)
			return nil, errors.New("Unable to find user")
		}

		err = bcrypt.CompareHashAndPassword(user.Password, append([]byte(form.Password), user.Salt...))
		if err != nil {
			log.Println("Password check failed")
			return nil, errors.New("Password check failed")
		}

		session.Values["loggedIn"] = true

		return nil, nil
	}

	return nil, errors.New("Already logged in")
}

func logout(req *http.Request, session *sessions.Session) (interface{}, error) {
	log.Println("Call - logout")
	// Clear the session
	session.Values = make(map[interface{}]interface{})

	return nil, nil
}
