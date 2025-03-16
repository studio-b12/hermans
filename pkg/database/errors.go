package database

import (
	"database/sql"
	"errors"

	"github.com/studio-b12/elk"
)

const (
	ErrNotFound = elk.ErrorCode("database:not-found")
	ErrDatabase = elk.ErrorCode("database:error")
)

func wrapErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return elk.Wrap(ErrNotFound, err, "no entry found in database")
	}

	return elk.Wrap(ErrDatabase, err, "database error")
}
