package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewSqLiteProvider(dbPath string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
