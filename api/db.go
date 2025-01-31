package main

import (
	"database/sql"
	"fmt"
	"log"
)

type DBInterface interface {
	Init() error
	Exec(query string, args ...interface{}) error
	Ping() error
}

var db DBInterface // We'll set this in main.go and in tests

// realDB implements DBInterface using *sql.DB
type realDB struct {
	db *sql.DB
}

func (r *realDB) Init() error {
	// For real usage, we expect the *sql.DB to be already constructed
	if r.db == nil {
		return fmt.Errorf("realDB has no database reference")
	}
	return nil
}

func (r *realDB) Exec(query string, args ...interface{}) error {
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *realDB) Ping() error {
	return r.db.Ping()
}

// initializeDBTable ensures the table for transactions is present.
func (r *realDB) initializeDBTable() {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			transaction_id VARCHAR(100) NOT NULL,
			amount NUMERIC(10, 2) NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`
	if err := r.Exec(createTableQuery); err != nil {
		log.Fatalf("Error creating transactions table: %v", err)
	}
}
