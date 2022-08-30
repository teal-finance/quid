package db

import (
	"github.com/jmoiron/sqlx"

	// pg import.
	_ "github.com/lib/pq"

	"github.com/teal-finance/emo"
)

var db *sqlx.DB

var logg = emo.NewLogger("db")

// Init : init the db conf.
func Init(isVerbose, isDev, isCmd bool) {
	if !isDev && !isCmd {
		logg.Print = isVerbose
	}
}

// Connect : connect to the db.
func Connect(dataSourceName string) error {
	_db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

// ExecSchema : execute the schema.
func ExecSchema() error {
	db.MustExec(schema)
	return nil
}
