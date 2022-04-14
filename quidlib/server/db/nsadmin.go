package db

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAdministratorsInNamespace : get the admin users in a namespace.
func SelectAdministratorsInNamespace(namespaceID int64) ([]server.NsAdmin, error) {
	data := []server.NsAdmin{}
	err := db.Select(&data,
		"SELECT namespaceadmin.id,namespaceadmin.user_id,namespaceadmin.namespace_id,usertable.username "+
			"FROM namespaceadmin "+
			"LEFT OUTER JOIN usertable on usertable.id=namespaceadmin.user_id "+
			"LEFT OUTER JOIN namespace on namespace.id=namespaceadmin.namespace_id "+
			"WHERE namespace.id=$1", namespaceID)
	if err != nil {
		fmt.Println("ERR", err)
		return data, err
	}

	return data, nil
}

// CreateAdministrator : create an admin user.
func CreateAdministrator(namespaceID int64, userID int64) (int64, error) {
	q := "INSERT INTO namespaceadmin(namespace_id, user_id) VALUES($1,$2) RETURNING id"
	rows, err := db.Query(q, namespaceID, userID)
	if err != nil {
		emo.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi interface{}
		err := rows.Scan(&idi)
		if err != nil {
			emo.QueryError(err)
			return 0, err
		}

		return idi.(int64), nil
	}

	emo.QueryError("no namespaceadmin for namespaceID=", namespaceID, " userID=", userID)
	return 0, fmt.Errorf("no namespaceadmin")

}

// AdministratorExists : check if an admin user exists.
func AdministratorExists(userID int64, namespaceID int64) (bool, error) {
	var n int
	q := "SELECT COUNT(id) FROM namespaceadmin WHERE (namespace_id=$1 AND user_id=$2)"
	err := db.Get(&n, q, namespaceID, userID)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}

// DeleteAdministrator : delete an admin user for a namespace.
func DeleteAdministrator(userID int64, namespaceID int64) error {
	q := "DELETE FROM namespaceadmin WHERE (user_id=$1 AND namespace_id=$2)"
	fmt.Println(q, userID, namespaceID)
	tx := db.MustBegin()
	tx.MustExec(q, userID, namespaceID)
	return tx.Commit()
}
