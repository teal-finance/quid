package db

import "time"

// base models to unpack data for Sqlx

type user struct {
	DateCreated time.Time `db:"date_created" json:"date_created"`
	Name        string    `db:"name" json:"username"`
	Password    string    `db:"password" json:"password"`
	Namespace   string    `db:"namespace" json:"namespace"`
	ID          int64     `db:"id" json:"id"`
	Enabled     bool      `db:"enabled" json:"enabled"`
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
	AccessKey             []byte `db:"access_key" json:"access_key"`
	RefreshKey            []byte `db:"refresh_key" json:"refresh_key"`
	MaxAccessTTL          string `db:"max_access_ttl" json:"max_access_ttl"`
	MaxRefreshTTL         string `db:"max_refresh_ttl" json:"max_refresh_ttl"`
	ID                    int64  `db:"id" json:"id"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}

type org struct {
	Name string `db:"name" json:"name"`
	ID   int64  `db:"id" json:"id"`
}
