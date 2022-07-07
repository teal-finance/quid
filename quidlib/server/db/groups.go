package db

import (
	// pg import.
	"fmt"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAllGroups : get all the groups.
func SelectAllGroups() ([]server.Group, error) {
	q := "SELECT grouptable.id,grouptable.name,namespace.name as namespace" +
		" FROM grouptable" +
		" JOIN namespace ON grouptable.namespace_id = namespace.id" +
		" ORDER BY grouptable.name"

	data := []server.Group{}
	err := db.Select(&data, q)
	return data, err
}

// SelectGroupsForUser : get the groups for a user.
func SelectGroupsForUser(userID int64) ([]server.Group, error) {
	q := "SELECT grouptable.id as id, grouptable.name as name" +
		" FROM usergroup" +
		" JOIN grouptable ON usergroup.group_id = grouptable.id" +
		" WHERE usergroup.user_id=$1 ORDER BY grouptable.name"

	data := []server.Group{}
	err := db.Select(&data, q, userID)
	return data, err
}

// SelectGroupsNamesForUser : get the groups for a user.
func SelectGroupsNamesForUser(userID int64) ([]string, error) {
	q := "SELECT grouptable.name as name FROM usergroup" +
		" JOIN grouptable ON usergroup.group_id = grouptable.id" +
		" WHERE usergroup.user_id=$1 ORDER BY grouptable.name"

	data := []userGroupName{}
	err := db.Select(&data, q, userID)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(data))
	for _, row := range data {
		names = append(names, row.Name)
	}

	return names, nil
}

// SelectGroupsForNamespace : get the groups for a namespace.
func SelectGroupsForNamespace(namespaceID int64) ([]server.Group, error) {
	q := "SELECT grouptable.id,grouptable.name,namespace.name as namespace" +
		" FROM grouptable" +
		" JOIN namespace ON grouptable.namespace_id = namespace.id" +
		" WHERE grouptable.namespace_id=$1 ORDER BY grouptable.name"

	data := []server.Group{}
	err := db.Select(&data, q, namespaceID)
	return data, err
}

// SelectGroup : get a group.
func SelectGroup(name string, namespaceID int64) (server.Group, error) {
	data := []server.Group{}
	err := db.Select(&data, "SELECT id,name FROM grouptable WHERE(name=$1 AND namespace_id=$2)", name, namespaceID)

	if len(data) == 0 {
		return server.Group{}, err
	}

	return data[0], err
}

// CreateGroup : create a group.
func CreateGroup(name string, namespaceID int64) (int64, error) {
	q := "INSERT INTO grouptable(name,namespace_id) VALUES($1,$2) RETURNING id"

	rows, err := db.Query(q, name, namespaceID)
	if err != nil {
		emo.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi any
		err = rows.Scan(&idi)
		if err != nil {
			emo.QueryError(err)
			return 0, err
		}
		return idi.(int64), nil
	}

	err = fmt.Errorf("no group %q", name)
	if err != nil {
		emo.QueryError(err)
	}

	return 0, err
}

// DeleteGroup : delete a group.
func DeleteGroup(id int64) error {
	q := "DELETE FROM grouptable WHERE id=$1"

	tx := db.MustBegin()
	tx.MustExec(q, id)

	return tx.Commit()
}

// GroupExists : check if an group exists.
func GroupExists(name string, namespaceID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM grouptable WHERE(name=$1 AND namespace_id=$2)"

	var n int
	err := db.Get(&n, q, name, namespaceID)
	exists := (n == 1)

	return exists, err
}

// AddUserInGroup : add a user into a group.
func AddUserInGroup(userID, groupID int64) error {
	q := "INSERT INTO usergroup(user_id,group_id) VALUES($1,$2)"

	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)

	return tx.Commit()
}

// RemoveUserFromGroup : remove a user from a group.
func RemoveUserFromGroup(userID, groupID int64) error {
	q := "DELETE FROM usergroup WHERE user_id=$1 AND group_id=$2"

	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)

	return tx.Commit()
}

// IsUserInGroup : check if a user is in a group.
func IsUserInGroup(userID, groupID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM usergroup WHERE(user_id=$1 AND group_id=$2)"

	var n int
	err := db.Get(&n, q, userID, groupID)
	exists := (n == 1)

	return exists, err
}
