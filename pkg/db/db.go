package db

import (
	"database/sql"
)

type Db struct {
	conn *sql.DB
}

func NewDb() *Db {
	conn := sql.Open("sqlite3")

	return &Db{conn: conn}
}
