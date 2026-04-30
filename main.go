package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/fmo/oauth2-client/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	app := handlers.NewApp()

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/callback", app.CallbackHandler)

	fmt.Println("Server starting on port 8081")
	http.ListenAndServe(":8081", mux)
}
