package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to reach database: %w", err)
	}

	return db, nil
}
