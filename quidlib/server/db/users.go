package db

import (
	"database/sql"
	"errors"
	"fmt"

	// pg import.
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectNonDisabledUser : get a user from it's username.
func SelectNonDisabledUser(username string, namespaceID int64) (bool, server.User, error) {
	ux := server.User{}

	row := db.QueryRowx("SELECT id,username,password,is_disabled FROM usertable WHERE(username=$1 AND namespace_id=$2)", username, namespaceID)

	u := user{}
	err := row.StructScan(&u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ux, nil
		}
		return false, ux, err
	}

	if u.IsDisabled {
		return false, ux, nil
	}

	ux.Name = u.UserName
	ux.PasswordHash = u.Password
	ux.ID = u.ID
	return true, ux, nil
}

// SelectAllUsers : get the users.
func SelectAllUsers() ([]server.User, error) {
	users := []server.User{}

	data := []user{}
	err := db.Select(&data,
		"SELECT usertable.id,usertable.username,namespace.name as namespace FROM usertable "+
			"JOIN namespace ON usertable.namespace_id = namespace.id "+
			"ORDER BY usertable.username")
	if err != nil {
		return users, err
	}

	for _, u := range data {
		users = append(users, server.User{
			ID:           u.ID,
			Name:         u.UserName,
			Namespace:    u.Namespace,
			PasswordHash: "",
			Org:          "",
		})
	}

	return users, nil
}

// SelectUsersInNamespace : get the users in a namespace.
func SelectUsersInNamespace(namespaceID int64) ([]server.User, error) {
	users := []server.User{}

	data := []user{}
	err := db.Select(&data,
		"SELECT usertable.id,usertable.username,namespace.name as namespace FROM usertable "+
			"JOIN namespace ON usertable.namespace_id = namespace.id  "+
			"WHERE usertable.namespace_id=$1 ORDER BY usertable.username", namespaceID)
	if err != nil {
		return users, err
	}

	for _, u := range data {
		users = append(users, server.User{
			ID:        u.ID,
			Name:      u.UserName,
			Namespace: u.Namespace,
		})
	}

	return users, nil
}

// SearchUsersInNamespaceFromUsername : get the users in a namespace from a username.
// TODO FIXME TOFIX
/*func SearchUsersInNamespaceFromUsername(username string, namespaceID int64) ([]server.User, error) {
	data := []server.User{}
	err := db.Select(&data, "SELECT id,username FROM usertable WHERE(username LIKE $1 AND namespace_id=$2)", username+"%", namespaceID)

	return data, err
}*/

// SelectUsersInGroup : get the users in a group.
func SelectUsersInGroup(username string, namespaceID int64) (server.Group, error) {
	q := "SELECT id,username FROM grouptable" +
		" WHERE(username=$1 AND namespace_id=$2)"

	data := []server.Group{}
	err := db.Select(&data, q, username, namespaceID)
	if err != nil {
		return server.Group{}, err
	}

	return data[0], nil
}

// CreateUser : create a user.
func CreateUser(username, password string, namespaceID int64) (server.User, error) {
	user := server.User{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	uid, err := CreateUserFromNameAndPassword(username, string(hashedPassword), namespaceID)
	if err != nil {
		return user, err
	}

	user.ID = uid
	user.Name = username

	return user, nil
}

// CreateUserFromNameAndPassword : create a user.
func CreateUserFromNameAndPassword(username, passwordHash string, namespaceID int64) (int64, error) {
	q := "INSERT INTO usertable(username,password,namespace_id) VALUES($1,$2,$3) RETURNING id"
	rows, err := db.Query(q, username, passwordHash, namespaceID)
	if err != nil {
		emo.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi any
		err := rows.Scan(&idi)
		if err != nil {
			emo.QueryError(err)
			return 0, err
		}
		return idi.(int64), nil
	}

	emo.QueryError("no user ", username)
	return 0, fmt.Errorf("no user %q", username)
}

/*
// SelectGroupsForUser : get the groups for a user in a namespace
func SelectGroupsForUser(userID int64) ([]server.Group, error) {
	data := []group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name FROM usergroup "+
		"JOIN grouptable ON usergroup.group_id = grouptable.id WHERE usergroup.user_id=$1 ORDER BY grouptable.name",
		userID)
	gr := []server.Group{}
	if err != nil {
		return gr, err
	}
	for _, g := range gr {
		res := server.Group{
			ID:   g.ID,
			Name: g.Name,
		}
		gr = append(gr, res)
	}
	return gr, nil
}*/

// CountUsersInGroup : count the users in a group.
func CountUsersInGroup(groupID int64) (int, error) {
	q := "SELECT COUNT(user_id) FROM usergroup WHERE group_id=$1"

	var n int
	err := db.Get(&n, q, groupID)

	return n, err
}

// UserNameExists : check if a username exists.
func UserNameExists(username string, namespaceID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM usertable WHERE (username=$1 AND namespace_id=$2)"

	var n int
	err := db.Get(&n, q, username, namespaceID)
	exists := (n > 0)

	return exists, err
}

// DeleteUser : delete a user.
func DeleteUser(id int64) error {
	q := "DELETE FROM usertable WHERE id=$1"

	tx := db.MustBegin()
	tx.MustExec(q, id)

	return tx.Commit()
}

// IsUserInAdminGroup : check if a user is in quid admin group
func IsUserInAdminGroup(uID, nsID int64) (bool, error) {
	g, err := SelectGroup("quid_admin", nsID)
	if err != nil {
		return false, err
	}
	return IsUserInGroup(uID, g.ID)
}
