package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

func (a *App) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	c, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "cant get cookie", http.StatusInternalServerError)
		return
	}

	if c.Value != state {
		log.Println("cookie val: ", c.Value, "state: ", state)
		http.Error(w, "bad state", http.StatusBadRequest)
		return
	}

	payload := url.Values{}
	payload.Set("client_id", a.ClientID)
	payload.Set("client_secret", a.ClientSecret)
	payload.Set("grant_type", "authorization_code")
	payload.Set("code", code)
	payload.Set("redirect_uri", a.RedirectURI)

	resp, _ := http.PostForm("http://localhost:8080/oauth/token", payload)

	var session Session

	json.NewDecoder(resp.Body).Decode(&session)

	sessionID, err := generateRandomString()
	if err != nil {
		http.Error(w, "cant generate session id", http.StatusInternalServerError)
		return
	}

	a.Sessions[sessionID] = session

	idToken, err := jwt.Parse(session.IDToken, func(token *jwt.Token) (any, error) {
		return []byte("my-secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Println(err)
		http.Error(w, "cant parse jwt", http.StatusInternalServerError)
		return
	}

	claims := make(map[string]any)

	if c, ok := idToken.Claims.(jwt.MapClaims); ok {
		claims = c
	}

	t, _ := template.ParseFiles("templates/callback.html")
	t.Execute(w, claims["sub"])
}
