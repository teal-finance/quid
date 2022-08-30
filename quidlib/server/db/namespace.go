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
	q := "SELECT id,name,alg,access_key,refresh_key,max_access_ttl,max_refresh_ttl,public_endpoint_enabled" +
		" FROM namespace WHERE name=$1"

	var ns server.Namespace

	row := db.QueryRowx(q, name)

	var data namespace
	err := row.StructScan(&data)
	if err != nil {
		logg.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return false, ns, nil
		}
		return true, ns, err
	}

	accessKey, err := crypt.AesGcmDecryptBin(data.AccessKey)
	if err != nil {
		logg.DecryptError(err)
		return true, ns, err
	}

	refreshKey, err := crypt.AesGcmDecryptBin(data.RefreshKey)
	if err != nil {
		logg.DecryptError(err)
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

// SelectNamespaceAccessPublicKey : get the AccessToken key for a namespace.
func SelectNamespaceAccessPublicKey(id int64) (found bool, algo string, der []byte, _ error) {
	row := db.QueryRowx("SELECT key FROM namespace WHERE id=$1", id)

	var data namespace
	if err := row.StructScan(&data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil, nil
		}

		logg.QueryError(err)
		return false, "", nil, err
	}

	private, err := crypt.AesGcmDecryptBin(data.AccessKey)
	if err != nil {
		logg.DecryptError(err)
		return true, "", nil, err
	}

	public, err := tokens.PrivateToPublicDER(data.SigningAlgo, private)
	if err != nil {
		logg.DecryptError(err)
		return true, "", nil, err
	}

	return true, data.SigningAlgo, public, nil
}

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
		logg.QueryError(err)
		return 0, err
	}

	for rows.Next() {
		var idi any
		err := rows.Scan(&idi)
		if err != nil {
			logg.QueryError(err)
			return 0, err
		}
		return idi.(int64), nil
	}

	logg.QueryError("no namespace", name)
	return 0, fmt.Errorf("no namespace %q", name)
}

// UpdateNamespaceTokenMaxTTL : update a max access token ttl for a namespace.
func UpdateNamespaceTokenMaxTTL(id int64, maxTTL string) error {
	q := "UPDATE namespace set max_access_ttl=$2 WHERE id=$1"
	_, err := db.Query(q, id, maxTTL)
	return err
}

// UpdateNamespaceRefreshTokenMaxTTL : update a max refresh token ttl for a namespace.
func UpdateNamespaceRefreshTokenMaxTTL(id int64, refreshMaxTTL string) error {
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
