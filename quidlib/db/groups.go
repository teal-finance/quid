package db

import (
	// pg import
	_ "github.com/lib/pq"

	"github.com/synw/quid/quidlib/models"
)

// SelectAllGroups : get all the groups
func SelectAllGroups() ([]models.Group, error) {
	data := []models.Group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name,namespace.name as namespace FROM grouptable "+
		"JOIN namespace ON grouptable.namespace_id = namespace.id ORDER BY grouptable.name")
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroupsForUser : get the groups for a user
func SelectGroupsForUser(userID int64) ([]models.Group, error) {
	data := []models.Group{}
	err := db.Select(&data, "SELECT grouptable.id as id, grouptable.name as name FROM usergroup "+
		"JOIN grouptable ON usergroup.group_id = grouptable.id WHERE usergroup.user_id=$1 ORDER BY grouptable.name",
		userID)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroupsNamesForUser : get the groups for a user
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

// SelectGroupsForNamespace : get the groups for a namespace
func SelectGroupsForNamespace(namespaceID int64) ([]models.Group, error) {
	data := []models.Group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name,namespace.name as namespace FROM grouptable "+
		"JOIN namespace ON grouptable.namespace_id = namespace.id WHERE grouptable.namespace_id=$1 ORDER BY grouptable.name",
		namespaceID)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectGroup : get a group
func SelectGroup(name string, namespaceID int64) (models.Group, error) {
	data := []models.Group{}
	err := db.Select(&data, "SELECT id,name FROM grouptable WHERE(name=$1 AND namespace_id=$2)", name, namespaceID)
	if err != nil {
		return data[0], err
	}
	return data[0], nil
}

// CreateGroup : create a group
func CreateGroup(name string, namespaceID int64) (int64, error) {
	q := "INSERT INTO grouptable(name,namespace_id) VALUES($1,$2) RETURNING id"
	rows, err := db.Query(q, name, namespaceID)
	if err != nil {
		return 0, err
	}
	var id int64
	for rows.Next() {
		var idi interface{}
		err := rows.Scan(&idi)
		if err != nil {
			return 0, err
		}
		id = idi.(int64)
	}
	return id, nil
}

// DeleteGroup : delete a group
func DeleteGroup(ID int64) error {
	q := "DELETE FROM grouptable WHERE id=$1"
	tx := db.MustBegin()
	tx.MustExec(q, ID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GroupExists : check if an group exists
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

// AddUserInGroup : add a user into a group
func AddUserInGroup(userID int64, groupID int64) error {
	q := "INSERT INTO usergroup(user_id,group_id) VALUES($1,$2)"
	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// RemoveUserFromGroup : remove a user from a group
func RemoveUserFromGroup(userID int64, groupID int64) error {
	q := "DELETE FROM usergroup WHERE user_id=$1 AND group_id=$2"
	tx := db.MustBegin()
	tx.MustExec(q, userID, groupID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// IsUserInGroup : check if a user is in a group
func IsUserInGroup(userID int64, groupID int64, namespaceID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM usergroup WHERE(user_id=$1 AND group_id=$2)"
	var n int
	err := db.Get(&n, q, userID, groupID)
	if err != nil {
		return false, err
	}
	exists := false
	if n == 1 {
		exists = true
	}
	return exists, nil
}
