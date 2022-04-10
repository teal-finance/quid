package db

import (
	// pg import
	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAllOrgs : get all the orgs
func SelectAllOrgs() ([]server.Org, error) {
	data := []server.Org{}
	err := db.Select(&data, "SELECT id,name FROM orgtable")
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectOrg : get a org
func SelectOrg(name string) (server.Org, error) {
	data := []server.Org{}
	err := db.Select(&data, "SELECT id,name FROM orgtable WHERE(name=$1)", name)
	if err != nil {
		return data[0], err
	}
	return data[0], nil
}

// SelectOrgsForUser : get the orgs for a user
func SelectOrgsForUser(userID int64) ([]server.Org, error) {
	data := []server.Org{}
	err := db.Select(&data, "SELECT orgtable.id as id, orgtable.name as name FROM userorg "+
		"JOIN orgtable ON userorg.org_id = orgtable.id WHERE userorg.user_id=$1 ORDER BY orgtable.name",
		userID)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SelectOrgStartsWith : get a namespace
func SelectOrgStartsWith(name string) ([]server.Org, error) {
	data := []org{}
	res := []server.Org{}
	err := db.Select(&data, "SELECT id,name FROM orgtable WHERE name LIKE '"+name+"%'")
	if err != nil {
		return res, err
	}
	for _, u := range data {
		res = append(res, server.Org{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	return res, nil
}

// SelectOrgsNamesForUser : get the orgs for a user
func SelectOrgsNamesForUser(userID int64) ([]string, error) {
	data := []userOrgName{}
	err := db.Select(&data, "SELECT orgtable.name as name FROM userorg "+
		"JOIN orgtable ON userorg.org_id = orgtable.id WHERE userorg.user_id=$1 ORDER BY orgtable.name",
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

// OrgExists : check if an org exists
func OrgExists(name string) (bool, error) {
	q := "SELECT COUNT(id) FROM orgtable WHERE(name=$1)"
	var n int
	err := db.Get(&n, q, name)
	exists := (n == 1)
	return exists, err
}

// DeleteOrg : delete an org
func DeleteOrg(ID int64) error {
	q := "DELETE FROM orgtable WHERE id=$1"
	tx := db.MustBegin()
	tx.MustExec(q, ID)
	return tx.Commit()
}

// CreateOrg : create an org
func CreateOrg(name string) (int64, error) {
	q := "INSERT INTO orgtable(name) VALUES($1) RETURNING id"
	rows, err := db.Query(q, name)
	if err != nil {
		return 0, err
	}

	var id int64
	for rows.Next() {
		var idi interface{}
		if err := rows.Scan(&idi); err != nil {
			return 0, err
		}

		id = idi.(int64)
	}

	return id, nil
}

// AddUserInOrg : add a user into an org
func AddUserInOrg(userID int64, orgID int64) error {
	q := "INSERT INTO userorg(user_id,org_id) VALUES($1,$2)"
	tx := db.MustBegin()
	tx.MustExec(q, userID, orgID)
	return tx.Commit()
}

// RemoveUserFromOrg : remove a user from an org
func RemoveUserFromOrg(userID int64, orgID int64) error {
	q := "DELETE FROM userorg WHERE user_id=$1 AND org_id=$2"
	tx := db.MustBegin()
	tx.MustExec(q, userID, orgID)
	return tx.Commit()
}
