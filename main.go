package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Rizzwaan/workoutVerse/internal/app"
	"github.com/Rizzwaan/workoutVerse/internal/routes"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "Server port")
	flag.Parse()

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	application.Logger.Printf("Listening on port %d", port)

	r := routes.SetupRoutes(application)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
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
