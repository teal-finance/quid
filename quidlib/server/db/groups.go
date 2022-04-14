package db

import (
	// pg import.
	"fmt"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAllGroups : get all the groups.
func SelectAllGroups() ([]server.Group, error) {
	data := []server.Group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name,namespace.name as namespace FROM grouptable "+
		"JOIN namespace ON grouptable.namespace_id = namespace.id ORDER BY grouptable.name")
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroupsForUser : get the groups for a user.
func SelectGroupsForUser(userID int64) ([]server.Group, error) {
	data := []server.Group{}
	err := db.Select(&data, "SELECT grouptable.id as id, grouptable.name as name FROM usergroup "+
		"JOIN grouptable ON usergroup.group_id = grouptable.id WHERE usergroup.user_id=$1 ORDER BY grouptable.name",
		userID)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroupsNamesForUser : get the groups for a user.
func SelectGroupsNamesForUser(userID int64) ([]string, error) {
	data := []userGroupName{}
	err := db.Select(&data, "SELECT grouptable.name as name FROM usergroup "+
		"JOIN grouptable ON usergroup.group_id = grouptable.id WHERE usergroup.user_id=$1 ORDER BY grouptable.name",
		userID)
	g := []string{}
	if err != nil {
		return g, err
	}
	for _, row := range data {
		g = append(g, row.Name)
	}
	return g, nil
}

// SelectGroupsForNamespace : get the groups for a namespace.
func SelectGroupsForNamespace(namespaceID int64) ([]server.Group, error) {
	data := []server.Group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name,namespace.name as namespace FROM grouptable "+
		"JOIN namespace ON grouptable.namespace_id = namespace.id WHERE grouptable.namespace_id=$1 ORDER BY grouptable.name",
		namespaceID)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroup : get a group.
func SelectGroup(name string, namespaceID int64) (server.Group, error) {
	data := []server.Group{}
	err := db.Select(&data, "SELECT id,name FROM grouptable WHERE(name=$1 AND namespace_id=$2)", name, namespaceID)
	if err != nil {
		return data[0], err
	}
	return data[0], nil
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
		var idi interface{}
		err := rows.Scan(&idi)
		if err != nil {
			emo.QueryError(err)
			return 0, err
		}

		return idi.(int64), nil
	}

	emo.QueryError("no group ", name)
	return 0, fmt.Errorf("no group %q", name)
}

// DeleteGroup : delete a group.
func DeleteGroup(id int64) error {
	q := "DELETE FROM grouptable WHERE id=$1"
	tx := db.MustBegin()
	tx.MustExec(q, id)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GroupExists : check if an group exists.
func GroupExists(name string, namespaceID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM grouptable WHERE(name=$1 AND namespace_id=$2)"
	var n int
	err := db.Get(&n, q, name, namespaceID)
	if err != nil {
		return false, err
	}
	exists := false
	if n == 1 {
		exists = true
	}
	return exists, nil
}

// AddUserInGroup : add a user into a group.
func AddUserInGroup(userID int64, groupID int64) error {
	q := "INSERT INTO usergroup(user_id,group_id) VALUES($1,$2)"
	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// RemoveUserFromGroup : remove a user from a group.
func RemoveUserFromGroup(userID int64, groupID int64) error {
	q := "DELETE FROM usergroup WHERE user_id=$1 AND group_id=$2"
	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// IsUserInGroup : check if a user is in a group.
func IsUserInGroup(userID int64, groupID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM usergroup WHERE(user_id=$1 AND group_id=$2)"
	var n int
	err := db.Get(&n, q, userID, groupID)
	exists := (n == 1)
	return exists, err
}
