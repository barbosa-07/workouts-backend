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

	// routes handler will go here
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status: OK\n")
}
