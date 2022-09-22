package db

import (
	"database/sql"
	"errors"

	// pg import.
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/teal-finance/quid/server"
)

// SelectEnabledUsrID : get a user id from it's username.
func SelectEnabledUsrID(username string) (bool, int64, error) {
	row := db.QueryRowx("SELECT id,username,password,enabled FROM users WHERE(username=$1)", username)
	var u user
	err := row.StructScan(&u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.NotFound("User", username, "not found")
			return false, 0, nil
		}
		return false, 0, err
	}
	// emo.Found("BASE USER", u.ID, u.UserName)
	if !u.Enabled {
		return false, 0, nil
	}
	// emo.Found("USER", u.ID)
	return true, u.ID, nil
}

// SelectEnabledUser : get a user from it's username.
func SelectEnabledUser(username string, nsID int64) (bool, server.User, error) {
	var usr server.User

	row := db.QueryRowx("SELECT id,username,password,enabled FROM users WHERE(username=$1 AND ns_id=$2)", username, nsID)

	var u user
	err := row.StructScan(&u)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.NotFound("User", username, "not found for", nsID)
			return false, usr, nil
		}
		return false, usr, err
	}

	if !u.Enabled {
		return false, usr, nil
	}

	usr.Name = u.UserName
	usr.PasswordHash = u.Password
	usr.Namespace = u.Namespace
	usr.ID = u.ID
	return true, usr, nil
}

// SelectAllUsers : get the users.
func SelectAllUsers() ([]server.User, error) {
	var data []user
	err := db.Select(&data,
		"SELECT users.id,users.username,namespaces.name as namespaces FROM users "+
			"JOIN namespaces ON users.ns_id = namespaces.id "+
			"ORDER BY users.username")
	if err != nil {
		log.S().Warning(err)
		return nil, err
	}

	users := make([]server.User, 0, len(data))
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

// SelectNsUsers : get the users in a namespace.
func SelectNsUsers(nsID int64) ([]server.User, error) {
	var data []user
	err := db.Select(&data,
		"SELECT users.id,users.username,namespaces.name as namespaces FROM users "+
			"JOIN namespaces ON users.ns_id = namespaces.id  "+
			"WHERE users.ns_id=$1 ORDER BY users.username", nsID)
	if err != nil {
		log.S().Warning(err)
		return nil, err
	}

	users := make([]server.User, 0, len(data))
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
// TODO FIXME
/*func SearchUsersInNamespaceFromUsername(username string, nsID int64) ([]server.User, error) {
	var data []server.User
	err := db.Select(&data, "SELECT id,username FROM users WHERE(username LIKE $1 AND ns_id=$2)", username+"%", nsID)
	if err != nil {
		log.S().Warning(err)
		return nil, err
	}
	return data, nil
}*/

// SelectUsersInGroup : get the users in a group.
func SelectUsersInGroup(username string, nsID int64) (server.Group, error) {
	q := "SELECT id,username FROM groups" +
		" WHERE(username=$1 AND ns_id=$2)"

	var data []server.Group
	err := db.Select(&data, q, username, nsID)
	if err != nil {
		log.S().Warning(err)
		return server.Group{}, err
	}
	if len(data) == 0 {
		return server.Group{}, log.Warn("SelectUsersInGroup is empty").Err()
	}

	return data[0], nil
}

// CreateUser : create a user.
func CreateUser(username, password string, nsID int64) (server.User, error) {
	var user server.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.S().Warning(err)
		return user, err
	}

	uid, err := CreateUserFromNameAndPassword(username, string(hashedPassword), nsID)
	if err != nil {
		log.S().Warning(err)
		return user, err
	}

	user.ID = uid
	user.Name = username

	return user, nil
}

// CreateUserFromNameAndPassword : create a user.
func CreateUserFromNameAndPassword(username, passwordHash string, nsID int64) (int64, error) {
	q := "INSERT INTO users(username,password,ns_id) VALUES($1,$2,$3) RETURNING id"
	rows, err := db.Query(q, username, passwordHash, nsID)
	if err != nil {
		log.QueryError(err)
		return 0, err
	}

	return getFirstID(username, rows)
}

/*
// SelectGroupsForUser : get the groups for a user in a namespace
func SelectGroupsForUser(usrID int64) ([]server.Group, error) {
	var data []group
	err := db.Select(&data, "SELECT groups.id,groups.name FROM user_groups "+
		"JOIN groups ON user_groups.grp_id = groups.id WHERE user_groups.usr_id=$1 ORDER BY groups.name",
		usrID)
	var gr []server.Group
	if err != nil {
		log.S().Warning(err)
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
func CountUsersInGroup(grpID int64) (int, error) {
	q := "SELECT COUNT(usr_id) FROM user_groups WHERE grp_id=$1"

	var n int
	err := db.Get(&n, q, grpID)

	return n, err
}

// UserExists : check if a username exists.
func UserExists(username string, nsID int64) (bool, error) {
	q := "SELECT COUNT(id) FROM users WHERE (username=$1 AND ns_id=$2)"

	var n int
	err := db.Get(&n, q, username, nsID)
	if err != nil {
		log.S().Warning(err)
		return false, err
	}

	exists := (n > 0)
	return exists, nil
}

// DeleteUser : delete a user.
func DeleteUser(id int64) error {
	q := "DELETE FROM users WHERE id=$1"

	tx := db.MustBegin()
	tx.MustExec(q, id)

	err := tx.Commit()
	if err != nil {
		log.S().Warning(err)
	}
	return err
}

// IsUserInAdminGroup : check if a user is in quid admin group
func IsUserInAdminGroup(uID, nsID int64) (bool, error) {
	g, err := SelectGroup("quid_admin", nsID)
	if err != nil {
		return false, err
	}
	return IsUserInGroup(uID, g.ID)
}
