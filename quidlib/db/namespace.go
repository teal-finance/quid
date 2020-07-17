package db

import (
	"database/sql"

	// pg import
	_ "github.com/lib/pq"

	"github.com/synw/quid/quidlib/models"
)

// SelectAllNamespaces : get the namespaces
func SelectAllNamespaces() ([]models.Namespace, error) {
	data := []namespace{}
	res := []models.Namespace{}
	err := db.Select(&data, "SELECT id,name,max_token_ttl,public_endpoint_enabled FROM namespace ORDER BY name")
	if err != nil {
		return res, err
	}
	for _, u := range data {
		res = append(res, models.Namespace{
			ID:                    u.ID,
			Name:                  u.Name,
			MaxTokenTTL:           u.MaxTokenTTL,
			PublicEndpointEnabled: u.PublicEndpointEnabled,
		})
	}
	return res, nil
}

// SelectNamespaceStartsWith : get a namespace
func SelectNamespaceStartsWith(name string) ([]models.Namespace, error) {
	data := []namespace{}
	res := []models.Namespace{}
	err := db.Select(&data, "SELECT id,name FROM namespace WHERE name LIKE '"+name+"%'")
	if err != nil {
		return res, err
	}
	for _, u := range data {
		res = append(res, models.Namespace{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	return res, nil
}

// SelectNamespaceKey : get the key for a namespace
func SelectNamespaceKey(ID int64) (bool, string, error) {
	data := namespace{}
	var key string
	row := db.QueryRowx("SELECT key FROM namespace WHERE id=$1", ID)
	err := row.StructScan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, key, nil
		}
		return false, key, err
	}
	key = data.Key
	return true, key, nil
}

// SelectNamespaceID : get a namespace
func SelectNamespaceID(name string) (int64, error) {
	data := []namespace{}
	err := db.Select(&data, "SELECT id,name FROM namespace WHERE name=$1", name)
	if err != nil {
		return 0, err
	}
	return data[0].ID, nil
}

// SetNamespaceEndpointAvailability
func SetNamespaceEndpointAvailability(ID int64, enable bool) error {
	q := "UPDATE namespace SET public_endpoint_enabled=$2 WHERE id=$1"
	_, err := db.Query(q, ID, enable)
	if err != nil {
		return err
	}
	return nil
}

// CreateNamespace : create a namespace
func CreateNamespace(name, key, ttl string, endpoint bool) (int64, error) {
	/*k, err := encrypt(key)
	if err != nil {
		return 0, err
	}*/
	q := "INSERT INTO namespace(name,key,max_token_ttl,public_endpoint_enabled) VALUES($1,$2,$3,$4) RETURNING id"
	rows, err := db.Query(q, name, key, ttl, endpoint)
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

// DeleteNamespace : delete a namespace
func DeleteNamespace(ID int64) error {
	q := "DELETE FROM namespace where id=$1"
	tx := db.MustBegin()
	tx.MustExec(q, ID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// NamespaceExists : check if an namespace exists
func NamespaceExists(name string) (bool, error) {
	q := "SELECT COUNT(id) FROM namespace WHERE(name=$1)"
	var n int
	err := db.Get(&n, q, name)
	if err != nil {
		return false, err
	}
	exists := false
	if n == 1 {
		exists = true
	}
	return exists, nil
}

// CountUsersForNamespace : count users in a namespace
func CountUsersForNamespace(ID int64) (int, error) {
	q := "SELECT COUNT(id) FROM usertable WHERE(namespace_id=$1)"
	var n int
	err := db.Get(&n, q, ID)
	if err != nil {
		return n, err
	}
	return n, nil
}
