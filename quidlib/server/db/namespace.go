package db

import (
	"database/sql"
	"errors"
	"fmt"

	// pg import.
	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/crypt"
	"github.com/teal-finance/quid/quidlib/server"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// SelectAllNamespaces : get the namespaces.
func SelectAllNamespaces() ([]server.Namespace, error) {
	var data []namespace
	err := db.Select(&data, "SELECT id,name,max_access_ttl,max_refresh_ttl,public_endpoint_enabled FROM namespace ORDER BY name")
	if err != nil {
		return nil, err
	}

	res := make([]server.Namespace, 0, len(data))
	for _, u := range data {
		res = append(res, server.Namespace{
			Name:                  u.Name,
			SigningAlgo:           "",
			AccessKey:             nil,
			RefreshKey:            nil,
			MaxTokenTTL:           u.MaxAccessTTL,
			MaxRefreshTokenTTL:    u.MaxRefreshTTL,
			ID:                    u.ID,
			PublicEndpointEnabled: u.PublicEndpointEnabled,
		})
	}

	return res, nil
}

// SelectNsStartsWith : get a namespace.
func SelectNsStartsWith(name string) ([]server.Namespace, error) {
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

// SelectNsFromName : get a namespace.
func SelectNsFromName(name string) (bool, server.Namespace, error) {
	q := "SELECT id,name,alg,access_key,refresh_key,max_access_ttl,max_refresh_ttl,public_endpoint_enabled" +
		" FROM namespace WHERE name=$1"

	var ns server.Namespace

	row := db.QueryRowx(q, name)

	var data namespace
	err := row.StructScan(&data)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return false, ns, nil
		}
		return true, ns, err
	}

	accessKey, err := crypt.AesGcmDecryptBin(data.AccessKey)
	if err != nil {
		log.DecryptError(err)
		return true, ns, err
	}

	refreshKey, err := crypt.AesGcmDecryptBin(data.RefreshKey)
	if err != nil {
		log.DecryptError(err)
		return true, ns, err
	}

	ns = server.Namespace{
		Name:                  data.Name,
		SigningAlgo:           "HS256",
		AccessKey:             accessKey,
		RefreshKey:            refreshKey,
		MaxTokenTTL:           data.MaxAccessTTL,
		MaxRefreshTokenTTL:    data.MaxRefreshTTL,
		ID:                    data.ID,
		PublicEndpointEnabled: data.PublicEndpointEnabled,
	}

	return true, ns, nil
}

// SelectVerificationKeyDER get the AccessToken key (in DER form) for a namespace.
func SelectVerificationKeyDER(id int64) (found bool, algo string, der []byte, _ error) {
	row := db.QueryRowx("SELECT key FROM namespace WHERE id=$1", id)

	var data namespace
	if err := row.StructScan(&data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil, nil
		}
		log.QueryError(err)
		return false, "", nil, err
	}

	der, err := tokens.DecryptVerificationKeyDER(data.SigningAlgo, data.AccessKey)
	if err != nil {
		log.Error(err)
		return true, "", nil, err
	}

	return true, data.SigningAlgo, der, nil
}

// SelectNsID : get a namespace.
func SelectNsID(name string) (int64, error) {
	var data []namespace
	err := db.Select(&data, "SELECT id,name FROM namespace WHERE name=$1", name)
	if err != nil {
		return 0, err
	}
	return data[0].ID, nil
}

// EnableNsEndpoint : enable or disable public endpoint.
func EnableNsEndpoint(id int64, enable bool) error {
	q := "UPDATE namespace SET public_endpoint_enabled=$2 WHERE id=$1"
	_, err := db.Query(q, id, enable)
	return err
}

// CreateNamespace : create a namespace.
func CreateNamespace(name, ttl, refreshTTL, algo string, accessKey, refreshKey []byte, endpoint bool) (int64, error) {
	q := "INSERT INTO namespace(name,alg,access_key,refresh_key,max_access_ttl,max_refresh_ttl,public_endpoint_enabled)" +
		" VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id"

	ak, err := crypt.AesGcmEncryptBin(accessKey)
	if err != nil {
		return 0, err
	}

	rk, err := crypt.AesGcmEncryptBin(refreshKey)
	if err != nil {
		return 0, err
	}

	rows, err := db.Query(q, name, algo, ak, rk, ttl, refreshTTL, endpoint)
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

	log.QueryError("no namespace", name)
	return 0, fmt.Errorf("no namespace %q", name)
}

// CreateNamespaceIfExist : create a namespace.
func CreateNamespaceIfExist(name, ttl, refreshMaxTTL, algo string, accessKey, refreshKey []byte, endpoint bool) (int64, bool, error) {
	exists, err := NamespaceExists(name)
	if err != nil {
		log.QueryError("createNamespace NamespaceExists:", err)
		return 0, false, err
	}
	if exists {
		log.QueryError("createNamespace: already exist")
		return 0, true, nil
	}

	nsID, err := CreateNamespace(name, ttl, refreshMaxTTL, algo, accessKey, refreshKey, endpoint)
	if err != nil {
		log.QueryError("createNamespace:", err)
		return 0, false, err
	}

	return nsID, false, nil
}

// UpdateNsTokenMaxTTL : update a max access token ttl for a namespace.
func UpdateNsTokenMaxTTL(id int64, maxTTL string) error {
	q := "UPDATE namespace set max_access_ttl=$2 WHERE id=$1"
	_, err := db.Query(q, id, maxTTL)
	return err
}

// UpdateNsRefreshMaxTTL : update a max refresh token ttl for a namespace.
func UpdateNsRefreshMaxTTL(id int64, refreshMaxTTL string) error {
	q := "UPDATE namespace set max_refresh_ttl=$2 WHERE id=$1"
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

// NamespaceExists : check if a namespace exists.
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
