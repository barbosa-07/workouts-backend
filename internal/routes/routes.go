package routes

import (
	"github.com/Rizzwaan/workoutVerse/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheckHandler)

	r.Get("/workouts/{id}", app.WorkoutHandler.HandleWorkoutById)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)

	return r
}
