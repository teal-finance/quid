package db

import (
	_ "github.com/lib/pq"
)

// SelectAdministrators : get the admin users in a namespace.
func SelectAdministrators(nsID int64) ([]Administrator, error) {
	q := "SELECT administrators.id, administrators.usr_id, administrators.ns_id, users.name " +
		"FROM administrators " +
		"LEFT OUTER JOIN users on users.id=administrators.usr_id " +
		"LEFT OUTER JOIN namespaces on namespaces.id=administrators.ns_id " +
		"WHERE namespaces.id=$1"
	log.Query(q, "nsID=", nsID)

	var data []Administrator
	err := db.Select(&data, q, nsID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}

// SelectNonAdministrators : find non admin users in a namespace
func SelectNonAdministrators(nsID int64, qs string) ([]NonAdmin, error) {
	q := "SELECT users.id as usr_id, users.name as name, namespaces.id as ns_id FROM users  " +
		"JOIN namespaces ON users.ns_id = namespaces.id " +
		"WHERE (namespaces.id = $1 AND users.name LIKE E'" + qs + "%') " +
		"AND users.id NOT IN ( " +
		"SELECT administrators.usr_id as id " +
		"FROM administrators " +
		"LEFT OUTER JOIN users on users.id = administrators.usr_id " +
		"LEFT OUTER JOIN namespaces on namespaces.id = administrators.ns_id" +
		" )"
	log.Query(q, "nsID=", nsID)

	var data []NonAdmin
	err := db.Select(&data, q, nsID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debug("Data", data)
	return data, nil
}

// CreateAdministrator : create an admin user.
func CreateAdministrator(nsID, usrID int64) error {
	q := "INSERT INTO administrators(ns_id, usr_id) VALUES($1,$2)"

	_, err := db.Query(q, nsID, usrID)
	if err != nil {
		log.QueryError(err)
	}
	return err
}

// IsUserAnAdmin : check if an admin user exists.
func IsUserAnAdmin(usrID, nsID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM administrators WHERE (ns_id=$1 AND usr_id=$2)"

	var n int
	err := db.Get(&n, q, nsID, usrID)
	if err != nil {
		log.Warn(err)
		return false, err
	}

	exists := (n > 0)
	return exists, nil
}

// DeleteAdministrator : delete an admin user for a namespace.
func DeleteAdministrator(usrID, nsID int64) error {
	q := "DELETE FROM administrators WHERE (usr_id=$1 AND ns_id=$2)"

	log.Data(q, usrID, nsID)

	tx := db.MustBegin()
	tx.MustExec(q, usrID, nsID)

	err := tx.Commit()
	if err != nil {
		log.Warn(err)
	}
	return err
}

type UserType int

const (
	UserNoAdmin = iota
	NsAdmin
	QuidAdmin
)

// GetUserType checks if a user is
func GetUserType(nsName string, nsID, usrID int64) (UserType, error) {
	if nsName == "quid" {
		// check if the user is in the quid admin group
		exists, err := IsUserInAdminGroup(usrID, nsID)
		if (err != nil) || !exists {
			return UserNoAdmin, err
		}
		return QuidAdmin, nil
	}

	// check if the user is namespace administrator
	exists, err := IsUserAnAdmin(usrID, nsID)
	if (err != nil) || !exists {
		return UserNoAdmin, err
	}
	return NsAdmin, nil
}
