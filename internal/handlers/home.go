// Package handlers
package handlers

import (
	"html/template"
	"net/http"
)

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	state, err := GenerateRandomString()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	authorizeURI, err := GenerateURI("http://localhost:8080/oauth/authorize", a.RedirectURI, a.ClientID, state)
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

	t.Execute(w, authorizeURI)
}
