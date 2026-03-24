// Package main adalah entry point untuk tool migrasi database.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/flashlab/vernon-license/migrations"
)

func main() {
	// Load .env jika ada
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://vernon:secret@localhost:5432/vernon_license?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	src := migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations.FS,
		Root:       ".",
	}

	direction := migrate.Up
	if len(os.Args) > 1 && os.Args[1] == "down" {
		direction = migrate.Down
	}

	n, err := migrate.Exec(db, "postgres", src, direction)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("Applied %d migration(s)\n", n)
}
