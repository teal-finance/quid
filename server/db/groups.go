package db

import (
	// pg import.

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/server"
)

// SelectAllGroups : get all the groups.
func SelectAllGroups() ([]server.Group, error) {
	q := "SELECT grouptable.id,grouptable.name,namespace.name as namespace" +
		" FROM grouptable" +
		" JOIN namespace ON grouptable.namespace_id = namespace.id" +
		" ORDER BY grouptable.name"

	var data []server.Group
	err := db.Select(&data, q)
	return data, err
}

// SelectGroupsForUser : get the groups for a user.
func SelectGroupsForUser(userID int64) ([]server.Group, error) {
	q := "SELECT grouptable.id as id, grouptable.name as name" +
		" FROM usergroup" +
		" JOIN grouptable ON usergroup.group_id = grouptable.id" +
		" WHERE usergroup.user_id=$1 ORDER BY grouptable.name"

	var data []server.Group
	err := db.Select(&data, q, userID)
	return data, err
}

// SelectGroupsNamesForUser : get the groups for a user.
func SelectGroupsNamesForUser(userID int64) ([]string, error) {
	q := "SELECT grouptable.name as name FROM usergroup" +
		" JOIN grouptable ON usergroup.group_id = grouptable.id" +
		" WHERE usergroup.user_id=$1 ORDER BY grouptable.name"

	var data []userGroupName
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

// SelectNsGroups : get the groups for a namespace.
func SelectNsGroups(namespaceID int64) ([]server.Group, error) {
	q := "SELECT grouptable.id,grouptable.name,namespace.name as namespace" +
		" FROM grouptable" +
		" JOIN namespace ON grouptable.namespace_id = namespace.id" +
		" WHERE grouptable.namespace_id=$1 ORDER BY grouptable.name"

	var data []server.Group
	err := db.Select(&data, q, namespaceID)
	return data, err
}

// SelectGroup : get a group.
func SelectGroup(name string, namespaceID int64) (server.Group, error) {
	var data []server.Group
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
		log.QueryError(err)
		return 0, err
	}

	return getFirstID(name, rows)
}

// createGroup : create a group.
func CreateGroupIfExist(name string, namespaceID int64) (server.Group, bool, error) {
	var grp server.Group

	exists, err := GroupExists(name, namespaceID)
	if err != nil {
		log.QueryError("createGroup GroupExists:", err)
		return grp, false, err
	}
	if exists {
		log.QueryError("createGroup: group already exists")
		return grp, true, nil
	}

	gid, err := CreateGroup(name, namespaceID)
	if err != nil {
		log.QueryError("createGroup:", err)
		return grp, false, err
	}

	grp = server.Group{
		Name:      name,
		Namespace: "",
		ID:        gid,
	}
	return grp, false, nil
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
	log.Query(q, userID, groupID)

	var n int
	err := db.Get(&n, q, userID, groupID)

	exists := (n == 1)
	return exists, err
}
