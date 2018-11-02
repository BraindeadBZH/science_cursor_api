package api

import (
	"encoding/json"
	"fmt"
	"github.com/pengux/check"
	"net/http"
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

func retrieveForm(form interface{}, req *http.Request, validation check.Struct) check.StructError {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(form)
	if err != nil {
		fmt.Println("Cannot retrieve form:", err.Error())
		return check.StructError{"form": make([]check.Error, 0)}
	}
	return validation.Validate(form)
}
