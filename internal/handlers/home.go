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

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		log.Println("[DEBUG] session_id - ", sessionCookie.Value)

		session, ok := a.Sessions[sessionCookie.Value]
		if !ok {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}

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

	state, err := GenerateRandomString()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	authorizeURI, err := a.GetAuthorizeURI(state)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Set state cookie to check later in callback
	http.SetCookie(w, &http.Cookie{
		Name:  "oauth_state",
		Value: state,
		Path:  "/",
	})

	resp.AuthorizeURI = authorizeURI

	t.Execute(w, resp)
}
