package db

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectNsAdministrators : get the admin users in a namespace.
func SelectNsAdministrators(namespaceID int64) ([]server.NsAdmin, error) {
	q := "SELECT namespaceadmin.id,namespaceadmin.user_id,namespaceadmin.namespace_id,usertable.username " +
		"FROM namespaceadmin " +
		"LEFT OUTER JOIN usertable on usertable.id=namespaceadmin.user_id " +
		"LEFT OUTER JOIN namespace on namespace.id=namespaceadmin.namespace_id " +
		"WHERE namespace.id=$1"

	var data []server.NsAdmin
	err := db.Select(&data, q, namespaceID)
	if err != nil {
		log.Error(err)
		return data, err
	}

	return data, nil
}

// SelectNonAdminUsersInNs : find non admin users in a namespace
func SelectNonAdminUsersInNs(namespaceID int64, qs string) ([]server.NonNsAdmin, error) {
	q := "SELECT usertable.id as user_id, usertable.username, namespace.id as namespace_id FROM usertable  " +
		"JOIN namespace ON usertable.namespace_id = namespace.id " +
		"WHERE (namespace.id = $1 AND usertable.username LIKE E'" + qs + "%') " +
		"AND usertable.id NOT IN ( " +
		"SELECT namespaceadmin.user_id as id " +
		"FROM namespaceadmin " +
		"LEFT OUTER JOIN usertable on usertable.id = namespaceadmin.user_id " +
		"LEFT OUTER JOIN namespace on namespace.id = namespaceadmin.namespace_id" +
		" )"
	log.Query(q, namespaceID)
	var data []server.NonNsAdmin
	err := db.Select(&data, q, namespaceID)
	if err != nil {
		log.Error(err)
		return data, err
	}

	log.Debug("Data", data)
	return data, nil
}

// CreateAdministrator : create an admin user.
func CreateAdministrator(namespaceID, userID int64) (int64, error) {
	q := "INSERT INTO namespaceadmin(namespace_id, user_id) VALUES($1,$2) RETURNING id"

	rows, err := db.Query(q, namespaceID, userID)
	if err != nil {
		log.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi any
		err := rows.Scan(&idi)
		if err != nil {
			log.QueryError(err)
			return 0, err
		}
		return idi.(int64), nil
	}

	log.QueryError("no nsAdmin for nsID=", namespaceID, "userID=", userID)
	return 0, fmt.Errorf("no namespaceadmin")
}

// IsUserAnAdmin : check if an admin user exists.
func IsUserAnAdmin(userID, namespaceID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM namespaceadmin WHERE (namespace_id=$1 AND user_id=$2)"

	var n int
	err := db.Get(&n, q, namespaceID, userID)
	exists := (n > 0)

	return exists, err
}

// DeleteAdministrator : delete an admin user for a namespace.
func DeleteAdministrator(userID, namespaceID int64) error {
	q := "DELETE FROM namespaceadmin WHERE (user_id=$1 AND namespace_id=$2)"

	log.Data(q, userID, namespaceID)

	tx := db.MustBegin()
	tx.MustExec(q, userID, namespaceID)

	return tx.Commit()
}

type UserType int

const (
	UserNoAdmin = iota
	NsAdmin
	QuidAdmin
)

// GetUserType checks if a user is
func GetUserType(nsName string, nsID, userID int64) (UserType, error) {
	if nsName == "quid" {
		// check if the user is in the quid admin group
		exists, err := IsUserInAdminGroup(userID, nsID)
		if (err != nil) || !exists {
			return UserNoAdmin, err
		}
		return QuidAdmin, nil
	}

	// check if the user is namespace administrator
	exists, err := IsUserAnAdmin(userID, nsID)
	if (err != nil) || !exists {
		return UserNoAdmin, err
	}
	return NsAdmin, nil
}
