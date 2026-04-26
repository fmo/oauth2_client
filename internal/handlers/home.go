// Package handlers
package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type Response struct {
	LoggedIn     bool
	AuthorizeURI string
	Username     string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	log.Println("")
	log.Println("[DEBUG] HomeHandler")

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] checking session cookie to see if the user already logged in")
	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		log.Println("[DEBUG] cookie exists and the session id in there is")

		session, ok := a.Sessions[sessionCookie.Value]
		if !ok {
			log.Println("[DEBUG] session id is not recorded so i'm doing to delete session cookie")
			UnsetCookie(w, "session_id")
		} else {
			log.Println("[DEBUG] session id in the cookie also in the system so i will get the claims")
			claims, err := GetClaims(session)
			if err != nil {
				http.Error(w, "cant get claims", http.StatusInternalServerError)
				return
			}

			resp.Username = claims["sub"].(string)

			resp.LoggedIn = true

			t.Execute(w, resp)
			return
		}
	} else {
		log.Println("[DEBUG] session cookie does not exist, so user is not logged in")
	}

	log.Println("[DEBUG] generating random string will be used for state to validate response is really coming from oauth provider")
	state, err := GenerateRandomString()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] generating authorize uri")
	authorizeURI, err := a.GetAuthorizeURI(state)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Println("[DEBUG] setting state cookie to check later when the response comes from oauth2/authorize")
	// Set state cookie to check later in callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
	})

	resp.AuthorizeURI = authorizeURI

	t.Execute(w, resp)
}
