package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/fmo/oauth2-client/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Set logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Initiate app
	app := handlers.NewApp(logger)

	// Router setup
	r := chi.NewRouter()

	r.Get("/", app.HomeHandler)
	r.Get("/callback", app.CallbackHandler)

	// Start server
	logger.Info("Server starting", "port", "8081")
	http.ListenAndServe(":8081", r)
}
