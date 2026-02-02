package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	return createSchema()
}

func createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL, -- Generic hash for app login
		encryption_salt BLOB NOT NULL -- Salt for deriving key encryption key
	);

	CREATE TABLE IF NOT EXISTS hosts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		alias TEXT NOT NULL,
		hostname TEXT NOT NULL,
		port INTEGER NOT NULL DEFAULT 22,
		username TEXT NOT NULL,
		encrypted_password BLOB, -- Encrypted with user's derived key
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err := DB.Exec(schema)
	if err != nil {
		log.Printf("Schema error: %v", err)
	}
	return err
}
