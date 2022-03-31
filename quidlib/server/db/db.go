package db

import (
	"github.com/jmoiron/sqlx"

	// pg import
	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/conf"
	emolib "github.com/teal-finance/quid/quidlib/emo"
)

var db *sqlx.DB

var emo = emolib.Zone{
	Name:    "db",
	NoPrint: true,
}

// Init : init the db conf
func Init(isVerbose bool) {
	emo.NoPrint = !isVerbose
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
