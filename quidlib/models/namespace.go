package models

// Namespace : base model
type Namespace struct {
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	Key                   string `db:"key" json:"-"`
	MaxTokenTTL           string `db:"max_token_ttl" json:"max_token_ttl"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}
