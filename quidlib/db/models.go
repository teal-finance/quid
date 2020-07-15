package db

import "time"

// base models to unpack data for Sqlx

type user struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Password    string    `db:"password" json:"password"`
	Namespace   string    `db:"namespace" json:"namespace"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

type group struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Namespace   string    `db:"namespace" json:"namespace"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

/*
type userGroup struct {
	ID      int32 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}*/

type userGroupName struct {
	Name string `db:"name" json:"name"`
}

type token struct {
	ID                  int32  `db:"id" json:"id"`
	Value               string `db:"value" json:"value"`
	UserID              int32  `db:"user_id" json:"user_id"`
	ExpirationTimestamp int32  `db:"expiration_timestamp" json:"expiration_timestamp"`
	Namespace           string `db:"namespace" json:"namespace"`
}

type namespace struct {
	ID                    int64  `db:"id" json:"id"`
	Name                  string `db:"name" json:"name"`
	Key                   string `db:"key" json:"key"`
	MaxTokenTTL           string `db:"max_token_ttl" json:"max_token_ttl"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}
