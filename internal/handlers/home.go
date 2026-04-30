// Package handlers
package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Response struct {
	SignedIn  bool
	SigninURI string
	Username  string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	slog := slog.With(
		"handler", "HomeHandler",
	)

	slog.Info("Started")

	username := a.IsUserSigned(w, r)
	if username != "" {
		slog.Info("User is alredy signed in")
		resp.Username = username
		resp.SignedIn = true
	} else {
		slog.Info("Generating random string for state")
		state, err := GenerateRandomString()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		slog.Info("Setting state cookie with the random string")
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
		})

		slog.Info("Generating signin uri")
		signinURI, err := a.GetAuthorizeURI(state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		resp.SigninURI = signinURI
	}

	slog.Info("Parsing and execute templates")
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	t.Execute(w, resp)
}
