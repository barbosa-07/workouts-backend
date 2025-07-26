package routes

import (
	"github.com/Rizzwaan/workoutVerse/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheckHandler)

	r.Get("/workout/{id}", app.WorkoutHandler.HandleWorkoutById)
	r.Post("/workout", app.WorkoutHandler.HandleCreateWorkout)

	return r
}
