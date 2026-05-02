// Package handlers
package handlers

import (
	"html/template"
	"net/http"
)

type Response struct {
	SignedIn  bool
	SigninURI string
	Username  string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	a.Logger.Info("===== HomeHandler =====")

	username := a.IsUserSigned(w, r)
	if username != "" {
		a.Logger.Info("User is alredy signed in")
		resp.Username = username
		resp.SignedIn = true
	} else {
		a.Logger.Info("Generating random string for state")
		state, err := GenerateRandomString()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		a.Logger.Info("Setting state cookie with the random string")
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
		})

		a.Logger.Info("Generating signin uri")
		signinURI, err := a.GetAuthorizeURI(state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp.SigninURI = signinURI
	}

	a.Logger.Info("Parsing and executing template")
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	t.Execute(w, resp)
}
