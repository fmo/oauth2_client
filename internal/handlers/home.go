// Package handlers
package handlers

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HomeViewData struct {
	SignedIn  bool
	SigninURI string
	Username  string
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeViewData := &HomeViewData{}

	a.Logger.Info("===== HomeHandler Start =====\n")
	a.Logger.WithField("client_id", a.ClientID).Info("Client for oauth sign-in flow")

	username := a.IsUserSigned(w, r)
	if username != "" {
		a.Logger.Info("User is alredy signed in")
		homeViewData.Username = username
		homeViewData.SignedIn = true
	} else {
		a.Logger.WithField("state", "").Info("Generating random string for state")
		state, err := GenerateRandomString()
		if err != nil {
			a.Logger.WithError(err).Error("Cant generate random string for state")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		a.Logger.WithField("state", state).Debug("State is created")

		a.Logger.Info("Setting state cookie with the random string")
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
		})
		a.Logger.WithFields(logrus.Fields{
			"cookie_name":  "auth_state",
			"cookie_value": state,
		}).Debug("Cookie Values")

		a.Logger.Info("Generating authorize uri")
		signinURI, err := a.GetAuthorizeURI(state)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		homeViewData.SigninURI = signinURI
		a.Logger.WithField("authorize_uri", homeViewData.SigninURI).Debug("Authorize URI")
	}

	a.Logger.Info("Parsing and executing template, ready to go to oauth provider")
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "cant render templates", http.StatusInternalServerError)
		return
	}

	a.Logger.Info("===== HomeHandler End =====")
	t.Execute(w, homeViewData)
}
