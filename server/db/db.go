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
	result, err := db.Exec(schema)
	if err != nil {
		return err
	}

	if result != nil {
		id, _ := result.LastInsertId()
		n, _ := result.RowsAffected()
		if id > 0 || n > 0 {
			log.V().Dataf("Created tables and indexes LastID=%d RowsAffected=%d", id, n)
		}
	}

	return nil
}
