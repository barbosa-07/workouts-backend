package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rizzwaan/workoutVerse/internal/api"
	"github.com/Rizzwaan/workoutVerse/internal/store"
	"github.com/Rizzwaan/workoutVerse/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)

	// routes handler will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userHandler := api.NewUserHandler(userStore, logger)

	// Initialize the application with handlers and database connection
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status: OK\n")
}
