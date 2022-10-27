package db

import "time"

// base models to unpack data for Sqlx

type user struct {
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	Name         string    `db:"name" json:"name"`
	PasswordHash string    `db:"password" json:"password"`
	Namespace    string    `db:"namespace" json:"namespace"`
	ID           int64     `db:"id" json:"id"`
	Enabled      bool      `db:"enabled" json:"enabled"`
}

type userGroupName struct {
	Name string `db:"name" json:"name"`
}

type userOrgName struct {
	Name string `db:"name" json:"name"`
}

type namespace struct {
	Name                  string `db:"name" json:"name"`
	Alg                   string `db:"alg" json:"alg"`
	EncryptedAccessKey    []byte `db:"access_key" json:"access_key"`
	EncryptedRefreshKey   []byte `db:"refresh_key" json:"refresh_key"`
	MaxAccessTTL          string `db:"max_access_ttl" json:"max_access_ttl"`
	MaxRefreshTTL         string `db:"max_refresh_ttl" json:"max_refresh_ttl"`
	ID                    int64  `db:"id" json:"id"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}

type org struct {
	Name string `db:"name" json:"name"`
	ID   int64  `db:"id" json:"id"`
}

// Administrator : base model.
type Administrator struct {
	ID    int64  `json:"id"      db:"id"`
	Name  string `json:"name"    db:"name"`
	UsrID int64  `json:"usr_id"  db:"usr_id"`
	NsID  int64  `json:"ns_id"   db:"ns_id"`
}

// NonAdmin : base model.
type NonAdmin struct {
	Name  string `json:"name"    db:"name"`
	UsrID int64  `json:"usr_id"  db:"usr_id"`
	NsID  int64  `json:"ns_id"   db:"ns_id"`
}
