package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"example/goserver/internal/db"
	httpHandler "example/goserver/internal/http"
	"example/goserver/internal/repository"
)

func main() {
	var err error
	database, err := db.Open("employees.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := repository.NewRepository(database)
	if err := repo.CreateTable(ctx); err != nil {
		log.Fatal(err)
	}

	handler := httpHandler.NewHandler(repo)

	// Establish endpoints
	http.HandleFunc("GET 	/employees", handler.ListEmployees)
	http.HandleFunc("POST 	/employees", handler.CreateEmployee)
	http.HandleFunc("GET 	/employees/", handler.GetEmployee)
	http.HandleFunc("DELETE /employees/", handler.DeleteEmployee)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "Welcome"
		json.NewEncoder(w).Encode(path)
	})

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(srv.ListenAndServe())
}
