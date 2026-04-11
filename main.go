package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

func main() {
	mux := http.NewServeMux()

	clientID := "web_client"
	clientSecret := "demo-client-secret"
	redirectURI := "http://localhost:8081/callback"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		t, err := template.ParseFiles("templates/home.html")
		if err != nil {
			http.Error(w, "cant render templates", http.StatusInternalServerError)
			return
		}

		u, err := url.Parse("http://localhost:8080/oauth/authorize")
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		state, err := generateState()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		q := u.Query()
		q.Add("response_type", "code")
		q.Add("redirect_uri", redirectURI)
		q.Add("client_id", clientID)
		q.Add("scope", "openid profile email")
		q.Add("state", state)

		u.RawQuery = q.Encode()

		http.SetCookie(w, &http.Cookie{
			Name:  "oauth_state",
			Value: state,
			Path:  "/",
		})
		t.Execute(w, u.String())
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
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
		payload.Set("client_id", clientID)
		payload.Set("client_secret", clientSecret)
		payload.Set("grant_type", "authorization_code")
		payload.Set("code", code)
		payload.Set("redirect_uri", redirectURI)

		_, _ = http.PostForm("http://localhost:8080/oauth/token", payload)

		t, _ := template.ParseFiles("templates/callback.html")
		t.Execute(w, r.URL.Query().Get("code"))
	})

	fmt.Println("Server starting on port 8081")
	http.ListenAndServe(":8081", mux)
}
