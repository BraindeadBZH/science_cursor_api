package api

import (
	"github.com/pengux/check"
)

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var loginFormValidation = check.Struct{
	"Email": check.Composite{
		check.NonEmpty{},
		check.Email{},
	},
	"Password": check.Composite{
		check.NonEmpty{},
	},
}
