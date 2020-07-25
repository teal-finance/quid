package db

import (
	"github.com/jmoiron/sqlx"

	// pg import
	_ "github.com/lib/pq"

	"github.com/synw/quid/quidlib/conf"
	emolib "github.com/synw/quid/quidlib/emo"
)

var db *sqlx.DB

var emo = emolib.Zone{
	Name:            "db",
	DeactivatePrint: true,
}

// Init : init the db conf
func Init(isDev bool) {
	emo.DeactivatePrint = !isDev
}

// Connect : connect to the db
func Connect() error {
	//fmt.Println("Connecting to database", conf.ConnStr)
	_db, err := sqlx.Connect("postgres", conf.ConnStr)
	if err != nil {
		return err
	}
	db = _db
	return nil
}

// ExecSchema : execute the schema
func ExecSchema() error {
	db.MustExec(schema)
	return nil
}
