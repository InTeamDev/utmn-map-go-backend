package database

import (
	"database/sql"
	"fmt"
	"log"
)

func Init(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}
