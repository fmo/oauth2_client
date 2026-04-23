// Package handlers
package handlers

import (
	"html/template"
	"net/http"
)

type Response struct {
	LoggedIn     bool
	AuthorizeURI string
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

	_, err = r.Cookie("session_id")
	if err == nil {
		resp.LoggedIn = true
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
