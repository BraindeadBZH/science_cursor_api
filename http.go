package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

const storeName = "science_cursor"

var router = mux.NewRouter()
var sessionStore *sessions.CookieStore

type funcWithSession func(r *http.Request, s *sessions.Session) (interface{}, error)

func httpListen() {
	sessionStore = sessions.NewCookieStore(apiConfig.HTTP.SessionKey)
	http.ListenAndServe(fmt.Sprintf("%s:%d", apiConfig.HTTP.Address, apiConfig.HTTP.Port), router)
}

func handlerWithSession(handler funcWithSession) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		session, err := sessionStore.Get(req, storeName)
		if err != nil {
			fmt.Println("Could not retrieve the session")
			return
		}

		answer, err := handler(req, session)

		// Response handling //
		if err != nil {
			if apiConfig.General.Production {
				http.Error(rw, "500 - Unexpected error", http.StatusInternalServerError)
			} else {
				http.Error(rw, fmt.Sprintf("500 - %s", err.Error()), http.StatusInternalServerError)
			}
		} else {
			err = session.Save(req, rw)
			if err != nil {
				fmt.Println("Could not save the session:", err.Error())
			}

			if answer == nil {
				rw.WriteHeader(http.StatusNoContent)
			} else {
				rw.WriteHeader(http.StatusOK)
				data, err := json.Marshal(answer)

				if err != nil {
					if apiConfig.General.Production {
						http.Error(rw, "500 - Unexpected error", http.StatusInternalServerError)
					} else {
						http.Error(rw, fmt.Sprintf("500 - %s", err.Error()), http.StatusInternalServerError)
					}
				}

				rw.Write(data)
			}
		}
	}
}
