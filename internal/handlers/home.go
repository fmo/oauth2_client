package handlers

import (
	"html/template"
	"net/http"
	"net/url"
)

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
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

	state, err := generateRandomString()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("redirect_uri", a.RedirectURI)
	q.Add("client_id", a.ClientID)
	q.Add("scope", "openid profile email")
	q.Add("state", state)

	u.RawQuery = q.Encode()

	http.SetCookie(w, &http.Cookie{
		Name:  "oauth_state",
		Value: state,
		Path:  "/",
	})
	t.Execute(w, u.String())
}
