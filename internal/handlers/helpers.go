package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
)

func GenerateRandomString() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
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
