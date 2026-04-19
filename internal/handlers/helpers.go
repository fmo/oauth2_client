package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
)

var ErrWrongState = errors.New("bad state")

func ValidateState(r *http.Request, state string) error {
	c, err := r.Cookie("oauth_state")
	if err != nil {
		return errors.New("cant get the cookie")
	}

	if c.Value != state {
		return ErrWrongState
	}

	return nil
}

func GenerateRandomString() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GenerateTokenExchangePayload(clientID, clientSecret, code, redirectURI string) url.Values {
	payload := url.Values{}
	payload.Set("client_id", clientID)
	payload.Set("client_secret", clientSecret)
	payload.Set("grant_type", "authorization_code")
	payload.Set("code", code)
	payload.Set("rediret_uri", redirectURI)

	return payload
}

func GenerateURI(base, redirectURI, clientID, state string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("response_type", "code")
	q.Add("redirect_uri", redirectURI)
	q.Add("client_id", clientID)
	q.Add("scope", "openid profile email")
	q.Add("state", state)

	u.RawQuery = q.Encode()

	return u.String(), nil
}
