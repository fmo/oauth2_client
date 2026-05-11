package main

import (
	"fmt"
	"net/http"

	"github.com/fmo/oauth2-client/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)

	app := handlers.NewApp(logger)

	r := chi.NewRouter()

	r.Get("/", app.HomeHandler)
	r.Get("/callback", app.CallbackHandler)

	fmt.Println("Server starting on port 8081")
	http.ListenAndServe(":8081", r)
}
