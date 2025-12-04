package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() {
	db, err := sql.Open("sqlite3", "people.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
