package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func RunMigrations(db *sql.DB) error {
	files, err := os.ReadDir("./auth-service/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	for _, file := range files {
		if file.IsDir() || file.Name() == "migrations.go" {
			continue
		}
		if err := applyMigration(db, "./auth-service/migrations/"+file.Name()); err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		log.Println("Migration applied:" + file.Name())

	}
	return nil
}

func applyMigration(db *sql.DB, filePath string) error {
	migration, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	return nil
}
