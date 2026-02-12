package db

import "database/sql"

func createUrls(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
						id SERIAL PRIMARY KEY,
						correlation_id TEXT,
						original_url TEXT NOT NULL UNIQUE,
						short_url TEXT NOT NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						deleted_flag BOOLEAN NOT NULL DEFAULT FALSE,
						user_id INT NOT NULL)`)
	return err
}
