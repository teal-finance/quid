package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAdministratorsInNamespace : get the admin users in a namespace
func SelectAdministratorsInNamespace(namespaceID int64) ([]server.User, error) {
	data := []user{}
	usrs := []server.User{}
	err := db.Select(&data,
		"SELECT namespaceadmin.id,namespaceadmin.user_id,namespace.name as namespace FROM namespaceadmin "+
			"JOIN namespace ON namespaceadmin.namespace_id = namespace.id "+
			//"JOIN user ON namespaceadmin.user_id = usertable.id "+
			"WHERE namespaceadmin.namespace_id=$1 ORDER BY user", namespaceID)
	if err != nil {
		return usrs, err
	}
	fmt.Println("RESULT", data)
	for _, u := range data {
		usrs = append(usrs, server.User{
			ID:        u.ID,
			UserName:  u.UserName,
			Namespace: u.Namespace,
		})
	}
	return usrs, nil
}

// CreateAdministrator : create an admin user
func CreateAdministrator(namespaceID int64, userID int64) (int64, error) {
	q := "INSERT INTO namespaceadmin(namespace_id, user_id) VALUES($1,$2) RETURNING id"
	rows, err := db.Query(q, namespaceID, userID)
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

// AdministratorExists : check if an admin user exists
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
