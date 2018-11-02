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

type funcWithSession func(r *http.Request, s *sessionData) (interface{}, error)

func httpListen() {
	sessionStore = sessions.NewCookieStore(apiConfig.HTTP.SessionKey)
	http.ListenAndServe(fmt.Sprintf("%s:%d", apiConfig.HTTP.Address, apiConfig.HTTP.Port), router)
}

func handlerWithSession(handler funcWithSession) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		// Retrieves client session from cookie
		clientSession, err := sessionStore.Get(req, storeName)
		if err != nil {
			fmt.Println("Could not retrieve the client session")
			return
		}
		var serverSession *sessionModel
		var needCreation = false
		if clientSession.IsNew {
			fmt.Println("Client session is new")
			needCreation = true
		} else {
			key, exists := clientSession.Values["key"]
			if !exists {
				fmt.Println("Client session do not contain key")
				needCreation = true
			} else {
				exists, _ = sessionExists(key.([]byte))
				if !exists {
					fmt.Println("Key is obsolete")
					needCreation = true
				}
			}
		}

		if needCreation {
			serverSession, err = createSession()
			if err != nil {
				fmt.Println("Could not create server session")
				return
			}
			clientSession.Values["key"] = serverSession.Key
		} else {
			key := clientSession.Values["key"]
			serverSession, err = getSession(key.([]byte))
			if err != nil {
				fmt.Println("Error while retrieving server session")
				return
			}
		}

		answer, err := handler(req, &serverSession.Data)

		// Response handling //
		if err != nil {
			if apiConfig.General.Production {
				http.Error(rw, "500 - Unexpected error", http.StatusInternalServerError)
			} else {
				http.Error(rw, fmt.Sprintf("500 - %s", err.Error()), http.StatusInternalServerError)
			}
		} else {
			err = clientSession.Save(req, rw)
			if err != nil {
				fmt.Println("Could not save the client session:", err.Error())
			}

			if serverSession.Data.MarkedForDeletion {
				err = serverSession.destroy()
				if err != nil {
					fmt.Println("Could not destroy the server session:", err.Error())
				}
			} else {
				err = serverSession.save()
				if err != nil {
					fmt.Println("Could not save the server session:", err.Error())
				}
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
