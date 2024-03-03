package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB // Declare the package variable for the database connection

func ConnectDB() (*sql.DB, error) {
	// Replace the connection string with your actual PostgreSQL connection string
	dsn := "host=localhost user=user-name password=strong-password dbname=user-name port=5432 sslmode=disable"

	// Open a connection to the database
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
