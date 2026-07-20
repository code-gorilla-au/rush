package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewSqLiteProvider(dbPath string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func WithTxnCtx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			return fmt.Errorf("tx rollback failed: %w", rErr)
		}

		return fmt.Errorf("txn rolled back due to error: %w", err)
	}

	return tx.Commit()
}
