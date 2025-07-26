package routes

import (
	"github.com/Rizzwaan/workoutVerse/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheckHandler)
	return r
}
