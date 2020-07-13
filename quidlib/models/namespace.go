package models

// Namespace : base model
type Namespace struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Key  string `db:"key" json:"-"`
}
