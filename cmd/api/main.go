package main

import (
	"encoding/json"
	"log"
	"net/http"

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

	repo := repository.NewRepository(database)
	if err := repo.CreateTable(); err != nil {
		log.Fatal(err)
	}

	handler := httpHandler.NewHandler(repo)

	// Establish endpoints
	http.HandleFunc("/employees", handler.Employees)
	http.HandleFunc("/employees/", handler.Employee)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "Welcome"
		json.NewEncoder(w).Encode(path)
	})

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
