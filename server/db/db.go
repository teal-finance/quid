package db

import (
	"strings"

	"github.com/jmoiron/sqlx"

	// pg import.
	_ "github.com/lib/pq"

	"github.com/teal-finance/emo"
)

var db *sqlx.DB

var log = emo.NewZone("db")

// Connect : connect to the db.
func Connect(url string) error {
	url = strings.Replace(url, "postgresql://", "postgres://", 1)

	_db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

// ExecSchema : execute the schema.
func ExecSchema() error {
	_, err := db.Exec(schema)
	return err
}
