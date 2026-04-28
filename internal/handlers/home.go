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

	logger := slog.With(
		"handler", "HomeHandler",
	)

	slog.SetDefault(logger)

	slog.Info("Started")

	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	slog.Info("checking session cookie to see if the user already logged in")
	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		slog.Info("cookie exists and the session id in there is")
		session, ok := a.Sessions[sessionCookie.Value]
		if !ok {
			slog.Info("session id is not recorded so i'm doing to delete session cookie")
			UnsetCookie(w, "session_id")
		} else {
			slog.Info("session id in the cookie also in the system so i will get the claims")
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
		slog.Info("session cookie does not exist, so user is not logged in")
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
