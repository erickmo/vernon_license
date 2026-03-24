// Package database menyediakan implementasi repository berbasis PostgreSQL menggunakan sqlx.
package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// New membuka koneksi ke PostgreSQL menggunakan databaseURL yang diberikan,
// mengkonfigurasi connection pool, dan memverifikasi koneksi dengan Ping.
func New(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database.New: open: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database.New: ping: %w", err)
	}

	return db, nil
}
