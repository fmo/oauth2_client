// Package handlers
package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Response struct {
	LoggedIn     bool
	AuthorizeURI string
	Username     string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	slog := slog.With(
		"handler", "HomeHandler",
	)

	slog.Info("Started")

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	username, err := a.IsUserSigned(w, r)
	if err == nil {
		resp.Username = username
		resp.LoggedIn = true
		t.Execute(w, resp)
	}

	slog.Info("generating random string will be used for state to validate response is really coming from oauth provider")
	state, err := GenerateRandomString()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	slog.Info("setting state cookie to check later when the response comes from oauth2/authorize")
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
	})

	slog.Info("generating authorize uri")
	authorizeURI, err := a.GetAuthorizeURI(state)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp.AuthorizeURI = authorizeURI

	t.Execute(w, resp)
}
