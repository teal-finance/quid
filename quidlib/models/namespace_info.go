package models

// NamespaceInfo : base model
type NamespaceInfo struct {
	NumUsers int     `json:"num_users"`
	Groups   []Group `json:"groups"`
}
