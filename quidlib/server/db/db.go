package db

import (
	"github.com/jmoiron/sqlx"

	// pg import.
	_ "github.com/lib/pq"

	emolib "github.com/teal-finance/emo"
)

var db *sqlx.DB

var emo = emolib.NewZone("db")

// Init : init the db conf.
func Init(isVerbose, isDev bool, isCmd bool) {
	if !isDev && !isCmd {
		emo.Print = isVerbose
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
