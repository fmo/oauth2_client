package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (a *App) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Validate state
	err := ValidateState(r, state)
	if err != nil {
		http.Error(w, "bad state", http.StatusBadRequest)
		return
	}

	// Auth Token exchange payload
	payload := GenerateTokenExchangePayload(a.ClientID, a.ClientSecret, code, a.RedirectURI)
	resp, _ := http.PostForm("http://localhost:8080/oauth/token", payload)

	if resp.StatusCode == http.StatusUnauthorized {
		http.Error(w, "wrong code has been sent", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := a.SaveSession(resp)
	if err != nil {
		http.Error(w, "cant create session", http.StatusInternalServerError)
		return
	}

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
