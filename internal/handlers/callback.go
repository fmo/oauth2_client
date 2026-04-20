package handlers

import (
	"html/template"
	"net/http"
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

	// Get claims
	claims, err := GetClaims(session)
	if err != nil {
		http.Error(w, "cant get the claims", http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("templates/callback.html")
	t.Execute(w, claims["sub"])
}
