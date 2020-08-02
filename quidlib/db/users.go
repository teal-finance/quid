package db

import (
	"database/sql"
	"log"

	// pg import
	_ "github.com/lib/pq"

	"github.com/synw/quid/quidlib/models"
)

// SelectNonDisabledUser : get a user from it's name
func SelectNonDisabledUser(name string, namespaceID int64) (bool, models.User, error) {
	u := user{}
	ux := models.User{}
	row := db.QueryRowx("SELECT id,name,password,is_disabled FROM usertable WHERE(name=$1 AND namespace_id=$2)", name, namespaceID)
	err := row.StructScan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ux, nil
		}
		return false, ux, err
	}
	if u.IsDisabled {
		return false, ux, nil
	}
	ux.Name = u.Name
	ux.PasswordHash = u.Password
	ux.ID = u.ID
	return true, ux, nil
}

// SelectAllUsers : get the users
func SelectAllUsers() ([]models.User, error) {
	data := []user{}
	usrs := []models.User{}
	err := db.Select(&data,
		"SELECT usertable.id,usertable.name,namespace.name as namespace FROM usertable "+
			"JOIN namespace ON usertable.namespace_id = namespace.id ORDER BY name")
	//err := db.Select(&data, "SELECT id,name,password FROM usertable ORDER BY name")
	if err != nil {
		return usrs, err
	}
	for _, u := range data {
		usrs = append(usrs, models.User{
			ID:        u.ID,
			Name:      u.Name,
			Namespace: u.Namespace,
		})
	}
	return usrs, nil
}

// SelectUsersInNamespace : get the users in a namespace
func SelectUsersInNamespace(namespaceID int64) ([]models.User, error) {
	data := []user{}
	usrs := []models.User{}
	err := db.Select(&data,
		"SELECT usertable.id,usertable.name,namespace.name as namespace FROM usertable WHERE usertable.namespace_id=$1"+
			"JOIN namespace ON usertable.namespace_id = namespace.id ORDER BY name", namespaceID)
	if err != nil {
		return usrs, err
	}
	for _, u := range data {
		usrs = append(usrs, models.User{
			ID:        u.ID,
			Name:      u.Name,
			Namespace: u.Namespace,
		})
	}
	return usrs, nil
}

// SelectUsersInGroup : get the users in a group
func SelectUsersInGroup(name string, namespaceID int64) (models.Group, error) {
	data := []models.Group{}
	err := db.Select(&data, "SELECT id,name FROM grouptable WHERE(name=$1 AND namespace_id=$2)", name, namespaceID)
	if err != nil {
		return data[0], err
	}
	return data[0], nil
}

// CreateUserFromNameAndPassword : create a user
func CreateUserFromNameAndPassword(name string, passwordHash string, namespaceID int64) (int64, error) {
	q := "INSERT INTO usertable(name,password,namespace_id) VALUES($1,$2,$3) RETURNING id"
	rows, err := db.Query(q, name, passwordHash, namespaceID)
	if err != nil {
		log.Fatal(err)
	}
	var id int64
	for rows.Next() {
		var idi interface{}
		err := rows.Scan(&idi)
		if err != nil {
			log.Fatalln(err)
		}
		id = idi.(int64)
	}
	return id, nil
}

/*
// SelectGroupsForUser : get the groups for a user in a namespace
func SelectGroupsForUser(userID int64) ([]models.Group, error) {
	data := []group{}
	err := db.Select(&data, "SELECT grouptable.id,grouptable.name FROM usergroup "+
		"JOIN grouptable ON usergroup.group_id = grouptable.id WHERE usergroup.user_id=$1 ORDER BY grouptable.name",
		userID)
	gr := []models.Group{}
	if err != nil {
		return gr, err
	}
	for _, g := range gr {
		res := models.Group{
			ID:   g.ID,
			Name: g.Name,
		}
		gr = append(gr, res)
	}
	return gr, nil
}*/

// CountUsersInGroup : count the users in a group
func CountUsersInGroup(groupID int64) (int, error) {
	var n int
	q := "SELECT COUNT(user_id) FROM usergroup WHERE group_id=$1"
	err := db.Get(&n, q, groupID)
	if err != nil {
		return n, err
	}
	return n, nil
}

// UserNameExists : check if a username exists
func UserNameExists(name string) (bool, error) {
	var n int
	q := "SELECT COUNT(id) FROM usertable WHERE name=$1"
	err := db.Get(&n, q, name)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}

// DeleteUser : delete a user
func DeleteUser(ID int64) error {
	q := "DELETE FROM usertable WHERE id=$1"
	tx := db.MustBegin()
	tx.MustExec(q, ID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
