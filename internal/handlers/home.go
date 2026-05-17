// Package handlers
package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type HomeViewData struct {
	SignedIn  bool
	SigninURI string
	Username  string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeViewData := &HomeViewData{}

	a.Logger.Info("Client for oauth sign-in flow", "client_id", a.ClientID)

	username := a.IsUserSigned(w, r)
	if username != "" {
		userGroup := slog.Group("user", "username", username)
		a.Logger.Info("User is alredy signed in", userGroup)
		homeViewData.Username = username
		homeViewData.SignedIn = true
	} else {
		a.Logger.Info("Generating random string for state")
		state, err := GenerateRandomString()
		if err != nil {
			a.Logger.Error("Cant generate random string for state", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		a.Logger.Debug("State is created", "state", state)

		a.Logger.Info("Setting state cookie with the random string")
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
		})
		a.Logger.Debug("Cookie Values", "cookie_name", "auth_state", "cookie_value", state)

		a.Logger.Info("Generating authorize uri")
		signinURI, err := a.GetAuthorizeURI(state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		homeViewData.SigninURI = signinURI
		a.Logger.Debug("Authorize URI", "authorize_uri", homeViewData.SigninURI)
	}

	a.Logger.Info("Parsing and executing template, ready to go to oauth provider")
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	t.Execute(w, homeViewData)
}
