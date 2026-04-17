package handlers

type App struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Sessions     map[string]Session
}

type Session struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	UserID      string `json:"user_id"`
}

func NewApp() *App {
	return &App{
		ClientID:     "web_client",
		ClientSecret: "demo-client-secret",
		RedirectURI:  "http://localhost:8081/callback",
		Sessions:     make(map[string]Session),
	}
}
