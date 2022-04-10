package db

import "time"

// base models to unpack data for Sqlx

type user struct {
	ID          int64     `db:"id" json:"id"`
	UserName    string    `db:"username" json:"username"`
	Password    string    `db:"password" json:"password"`
	Namespace   string    `db:"namespace" json:"namespace"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	IsDisabled  bool      `db:"is_disabled" json:"is_disabled"`
}

type userGroupName struct {
	Name string `db:"name" json:"name"`
}

type userOrgName struct {
	Name string `db:"name" json:"name"`
}

type namespace struct {
	ID                    int64  `db:"id" json:"id"`
	Name                  string `db:"name" json:"name"`
	Key                   string `db:"key" json:"key"`
	RefreshKey            string `db:"refresh_key" json:"refresh_key"`
	MaxTokenTTL           string `db:"max_token_ttl" json:"max_token_ttl"`
	MaxRefreshTokenTTL    string `db:"max_refresh_token_ttl" json:"max_refresh_token_ttl"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}

type org struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
