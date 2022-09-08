package db

import (
	// pg import.
	"fmt"

	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/pkg/server"
)

// SelectAllOrgs : get all the orgs.
func SelectAllOrgs() ([]server.Org, error) {
	q := "SELECT id,name FROM orgtable"

	var data []server.Org
	err := db.Select(&data, q)

	return data, err
}

// SelectOrg : get a org.
func SelectOrg(name string) (server.Org, error) {
	q := "SELECT id,name FROM orgtable WHERE(name=$1)"

	var data []server.Org
	err := db.Select(&data, q, name)
	if err != nil {
		return server.Org{}, err
	}

	return data[0], nil
}

// SelectOrgsForUser : get the orgs for a user.
func SelectOrgsForUser(userID int64) ([]server.Org, error) {
	q := "SELECT orgtable.id as id, orgtable.name as name FROM userorg " +
		"JOIN orgtable ON userorg.org_id = orgtable.id " +
		"WHERE userorg.user_id=$1 ORDER BY orgtable.name"

	var data []server.Org
	err := db.Select(&data, q, userID)

	return data, err
}

// SelectOrgStartsWith : get a namespace.
func SelectOrgStartsWith(name string) ([]server.Org, error) {
	q := "SELECT id,name FROM orgtable WHERE name LIKE '" + name + "%'"

	var data []org
	err := db.Select(&data, q)
	if err != nil {
		return nil, err
	}

	res := make([]server.Org, 0, len(data))
	for _, u := range data {
		res = append(res, server.Org{
			ID:   u.ID,
			Name: u.Name,
		})
	}

	return res, nil
}

// SelectOrgsNamesForUser : get the orgs for a user.
func SelectOrgsNamesForUser(userID int64) ([]string, error) {
	q := "SELECT orgtable.name as name FROM userorg " +
		"JOIN orgtable ON userorg.org_id = orgtable.id " +
		"WHERE userorg.user_id=$1 ORDER BY orgtable.name"

	var data []userOrgName
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

// OrgExists : check if an org exists.
func OrgExists(name string) (bool, error) {
	q := "SELECT COUNT(id) FROM orgtable WHERE(name=$1)"

	var n int
	err := db.Get(&n, q, name)
	exists := (n == 1)

	return exists, err
}

// DeleteOrg : delete an org.
func DeleteOrg(id int64) error {
	q := "DELETE FROM orgtable WHERE id=$1"

	tx := db.MustBegin()
	tx.MustExec(q, id)

	return tx.Commit()
}

// CreateOrg : create an org.
func CreateOrg(name string) (int64, error) {
	q := "INSERT INTO orgtable(name) VALUES($1) RETURNING id"

	rows, err := db.Query(q, name)
	if err != nil {
		log.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi any
		err := rows.Scan(&idi)
		if err != nil {
			log.QueryError(err)
			return 0, err
		}
		return idi.(int64), nil
	}

	log.QueryError("no org", name)
	return 0, fmt.Errorf("no org %q", name)
}

// CreateOrgIfExist : create an org.
func CreateOrgIfExist(name string) (server.Org, bool, error) {
	var org server.Org

	exists, err := OrgExists(name)
	if err != nil {
		log.QueryError("createOrg OrgExists:", err)
		return org, false, err
	}
	if exists {
		log.QueryError("createOrg: already exist:", name)
		return org, true, nil
	}

	id, err := CreateOrg(name)
	if err != nil {
		log.QueryError("createOrg:", err)
		return org, false, err
	}

	org.ID = id
	org.Name = name
	return org, false, nil
}

// AddUserInOrg : add a user into an org.
func AddUserInOrg(userID, orgID int64) error {
	q := "INSERT INTO userorg(user_id,org_id) VALUES($1,$2)"

	tx := db.MustBegin()
	tx.MustExec(q, userID, orgID)

	return tx.Commit()
}

// RemoveUserFromOrg : remove a user from an org.
func RemoveUserFromOrg(userID, orgID int64) error {
	q := "DELETE FROM userorg WHERE user_id=$1 AND org_id=$2"

	tx := db.MustBegin()
	tx.MustExec(q, userID, orgID)

	return tx.Commit()
}
