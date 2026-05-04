package main

import (
	"fmt"
	"net/http"

	"github.com/fmo/oauth2-client/internal"
	"github.com/fmo/oauth2-client/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	logger := internal.NewLogger(internal.Debug)

	app := handlers.NewApp(logger)

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/callback", app.CallbackHandler)

	fmt.Println("Server starting on port 8081")
	http.ListenAndServe(":8081", mux)
}
