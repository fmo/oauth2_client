package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fmo/oauth2-client/internal"
)

type App struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Sessions     map[string]*Session
	AuthServer   string
	Logger       *internal.Logger
}

type Session struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	UserID      string `json:"user_id"`
}

func NewApp(l *internal.Logger) *App {
	return &App{
		ClientID:     "web_client",
		ClientSecret: "demo-client-secret",
		RedirectURI:  "http://localhost:8081/callback",
		Sessions:     make(map[string]*Session),
		AuthServer:   "http://localhost:8080",
		Logger:       l,
	}
}

func (a *App) GetAuthorizeURI(state string) (string, error) {
	base := fmt.Sprintf("%s/oauth/authorize", a.AuthServer)

	return GenerateURI(base, a.RedirectURI, a.ClientID, state)
}

func (a *App) SaveSession(resp *http.Response, w http.ResponseWriter) (*Session, error) {
	var session *Session
	json.NewDecoder(resp.Body).Decode(&session)

	sessionID, err := GenerateRandomString()
	if err != nil {
		return nil, errors.New("cant generate session")
	}

	a.Sessions[sessionID] = session

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Path:     "/",
		Value:    sessionID,
		Expires:  time.Now().Add(60 * time.Minute),
		HttpOnly: true,
	})

	return session, nil
}

func (a *App) IsUserSigned(w http.ResponseWriter, r *http.Request) string {
	a.Logger.Info("Checking session cookie if it exists")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		a.Logger.Info("Session cookie does not exist, so user is not logged in")
		return ""
	}

	a.Logger.Info("Session cookie exists")
	session, ok := a.Sessions[sessionCookie.Value]
	if !ok {
		a.Logger.Info("Session id is not recored so deleting session cookie")
		UnsetCookie(w, "session_id")
		return ""
	}

	a.Logger.Info("Session id in the cookie also in the system so getting the claims")
	claims, err := GetClaims(session)
	if err != nil {
		return ""
	}

	return claims["sub"].(string)
}
