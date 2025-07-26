package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Rizzwaan/workoutVerse/internal/app"
)

func main() {
	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	application.Logger.Println("Application started successfully")

	http.HandleFunc("/health", HealthCheckHandler)

	server := &http.Server{
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		application.Logger.Fatalf("Server failed to start: %v", err)
	} else {
		application.Logger.Println("Server is running on port 8080")
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status: OK\n")
}
