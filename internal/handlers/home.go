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

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	session, err := r.Cookie("session_id")
	if err == nil {
		resp.LoggedIn = true
		log.Println("[DEBUG] session_id - ", session.Value)
		claims, err := GetClaims(a.Sessions[session.Value])
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp.Username = claims["sub"].(string)

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
