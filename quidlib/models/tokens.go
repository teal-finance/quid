package models

// Token : base model
type Token struct {
	ID                  int64
	Value               string
	User                User
	UserGroups          []Group
	ExpirationTimestamp int32
	Namespace           string `json:"namespace"`
}
