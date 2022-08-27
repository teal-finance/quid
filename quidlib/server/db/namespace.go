package db

import (
	"database/sql"
	"errors"
	"fmt"

	// pg import.
	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/crypt"
	"github.com/teal-finance/quid/quidlib/server"
)

// SelectAllNamespaces : get the namespaces.
func SelectAllNamespaces() ([]server.Namespace, error) {
	var data []namespace
	err := db.Select(&data, "SELECT id,name,max_token_ttl,max_refresh_token_ttl,public_endpoint_enabled FROM namespace ORDER BY name")
	if err != nil {
		return nil, err
	}

	res := make([]server.Namespace, 0, len(data))
	for _, u := range data {
		res = append(res, server.Namespace{
			ID:                    u.ID,
			Name:                  u.Name,
			MaxTokenTTL:           u.MaxTokenTTL,
			MaxRefreshTokenTTL:    u.MaxRefreshTokenTTL,
			PublicEndpointEnabled: u.PublicEndpointEnabled,
			AccessKey:             "",
			RefreshKey:            "",
		})
	}

	return res, nil
}

// SelectNamespaceStartsWith : get a namespace.
func SelectNamespaceStartsWith(name string) ([]server.Namespace, error) {
	var data []namespace
	err := db.Select(&data, "SELECT id,name FROM namespace WHERE name LIKE '"+name+"%'")
	if err != nil {
		return nil, err
	}

	res := make([]server.Namespace, 0, len(data))
	for _, u := range data {
		res = append(res, server.Namespace{
			ID:   u.ID,
			Name: u.Name,
		})
	}

	return res, nil
}

// SelectNamespaceFromName : get a namespace.
func SelectNamespaceFromName(name string) (bool, server.Namespace, error) {
	q := "SELECT id,name,key,refresh_key,max_token_ttl,max_refresh_token_ttl,public_endpoint_enabled" +
		" FROM namespace WHERE name=$1"
	// emo.Query(q, name)

	var ns server.Namespace

	row := db.QueryRowx(q, name)

	var data namespace
	err := row.StructScan(&data)
	if err != nil {
		emo.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return false, ns, nil
		}
		return true, ns, err
	}

	ns.AccessKey, err = crypt.AesGcmDecrypt(data.AccessKey, nil)
	if err != nil {
		emo.DecryptError(err)
		return true, ns, err
	}

	ns.RefreshKey, err = crypt.AesGcmDecrypt(data.RefreshKey, nil)
	if err != nil {
		emo.DecryptError(err)
		return true, ns, err
	}

	ns.ID = data.ID
	ns.Name = data.Name
	ns.MaxTokenTTL = data.MaxTokenTTL
	ns.MaxRefreshTokenTTL = data.MaxRefreshTokenTTL
	ns.PublicEndpointEnabled = data.PublicEndpointEnabled

	return true, ns, nil
}

// SelectNamespaceAccessKey : get the AccessToken key for a namespace.
func SelectNamespaceAccessKey(id int64) (bool, string, error) {
	row := db.QueryRowx("SELECT key FROM namespace WHERE id=$1", id)

	var data namespace
	if err := row.StructScan(&data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil
		}

		emo.QueryError(err)
		return false, "", err
	}

	key, err := crypt.AesGcmDecrypt(data.AccessKey, nil)
	if err != nil {
		emo.DecryptError(err)
		return true, "", err
	}

	return true, key, nil
}

/*
// SelectNamespaceKeys : get the refresh and access keys for a namespace
func SelectNamespaceKeys(name string) (bool, string, string, error) {
	var data namespace
	row := db.QueryRowx("SELECT key,refresh_key FROM namespace WHERE name=$1", name)
	err := row.StructScan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", "", nil
		}
		return false, "", "", err
	}
	//emo.Decrypt("Decrypting refresh key data for the namespace id " + strconv.FormatInt(ID, 10))
	rk, err := aesGcmDecrypt(data.RefreshKey, nil)
	if err != nil {
		emo.DecryptError(err)
		return true, "", "", err
	}
	ak, err := aesGcmDecrypt(data.AccessKey, nil)
	if err != nil {
		emo.DecryptError(err)
		return true, "", "", err
	}

	return true, rk, ak, nil
}*/

// SelectNamespaceID : get a namespace.
func SelectNamespaceID(name string) (int64, error) {
	var data []namespace
	err := db.Select(&data, "SELECT id,name FROM namespace WHERE name=$1", name)
	if err != nil {
		return 0, err
	}
	return data[0].ID, nil
}

// SetNamespaceEndpointAvailability : enable or disable public endpoint.
func SetNamespaceEndpointAvailability(id int64, enable bool) error {
	q := "UPDATE namespace SET public_endpoint_enabled=$2 WHERE id=$1"
	_, err := db.Query(q, id, enable)
	return err
}

// CreateNamespace : create a namespace.
func CreateNamespace(name, key, refreshKey, ttl, refreshTTL string, endpoint bool) (int64, error) {
	k, err := crypt.AesGcmEncrypt(key, nil)
	if err != nil {
		return 0, err
	}

	rk, err := crypt.AesGcmEncrypt(refreshKey, nil)
	if err != nil {
		return 0, err
	}

	q := "INSERT INTO namespace(name,key,refresh_key,max_token_ttl,max_refresh_token_ttl,public_endpoint_enabled)" +
		" VALUES($1,$2,$3,$4,$5,$6) RETURNING id"
	rows, err := db.Query(q, name, k, rk, ttl, refreshTTL, endpoint)
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

	emo.QueryError("no namespace", name)
	return 0, fmt.Errorf("no namespace %q", name)
}

// UpdateNamespaceTokenMaxTTL : update a max access token ttl for a namespace.
func UpdateNamespaceTokenMaxTTL(id int64, maxTTL string) error {
	q := "UPDATE namespace set max_token_ttl=$2 WHERE id=$1"
	_, err := db.Query(q, id, maxTTL)
	return err
}

// UpdateNamespaceRefreshTokenMaxTTL : update a max refresh token ttl for a namespace.
func UpdateNamespaceRefreshTokenMaxTTL(id int64, refreshMaxTTL string) error {
	q := "UPDATE namespace set max_refresh_token_ttl=$2 WHERE id=$1"
	_, err := db.Query(q, id, refreshMaxTTL)
	return err
}

// DeleteNamespace : delete a namespace.
func DeleteNamespace(id int64) QueryResult {
	q := "DELETE FROM namespace where id=$1"
	tx, _ := db.Begin()

	_, err := tx.Exec(q, id)
	if err != nil {
		return queryError(err)
	}

	err = tx.Commit()
	if err != nil {
		return queryError(err)
	}

	return queryNoError()
}

// NamespaceExists : check if an namespace exists.
func NamespaceExists(name string) (bool, error) {
	q := "SELECT COUNT(id) FROM namespace WHERE(name=$1)"

	var n int
	err := db.Get(&n, q, name)
	exists := (n == 1)

	return exists, err
}

// CountUsersForNamespace : count users in a namespace.
func CountUsersForNamespace(id int64) (int, error) {
	q := "SELECT COUNT(id) FROM usertable WHERE(namespace_id=$1)"

	var n int
	err := db.Get(&n, q, id)

	return n, err
}
