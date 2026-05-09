// Package handlers
package handlers

import (
	"html/template"
	"net/http"
)

type HomeViewData struct {
	SignedIn  bool
	SigninURI string
	Username  string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeViewData := &HomeViewData{}

	a.Logger.Info("===== HomeHandler =====\n")

	a.Logger.Info("Client Id: %s", a.ClientID)

	username := a.IsUserSigned(w, r)
	if username != "" {
		a.Logger.Info("User is alredy signed in")
		homeViewData.Username = username
		homeViewData.SignedIn = true
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
		a.Logger.Debug("Cookie name: %s, cookie value: %s", "auth_state", state)

		a.Logger.Info("Generating authorize uri")
		signinURI, err := a.GetAuthorizeURI(state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		homeViewData.SigninURI = signinURI
		a.Logger.Debug("Authorize URI: %s", homeViewData.SigninURI)
	}

	a.Logger.Info("Parsing and executing template")
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	a.Logger.Info("\n")
	t.Execute(w, homeViewData)
}
