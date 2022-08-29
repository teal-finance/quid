package server

//go:generate go run github.com/mailru/easyjson/... -all -byte -disable_members_unescape -disallow_unknown_fields -snake_case ${GOFILE}

// Namespace : base model.
type Namespace struct {
	Name                  string `json:"name"`
	SigningAlgo           string `db:"alg" json:"-"`
	AccessKey             []byte `db:"access_key" json:"-"`
	RefreshKey            []byte `db:"refresh_key" json:"-"`
	MaxTokenTTL           string `db:"max_access_ttl" json:"max_access_ttl"`
	MaxRefreshTokenTTL    string `db:"max_refresh_ttl" json:"max_refresh_ttl"`
	ID                    int64  `json:"id"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}

// NamespaceInfo : base model.
type NamespaceInfo struct {
	Groups   []Group `json:"groups"`
	NumUsers int     `json:"num_users"`
}

// NsAdmin : base model.
type NsAdmin struct {
	UserName    string `db:"username" json:"username"`
	ID          int64  `db:"id" json:"id"`
	UserID      int64  `db:"user_id" json:"user_id"`
	NamespaceID int64  `db:"namespace_id" json:"namespace_id"`
}

// NonNsAdmin : base model.
type NonNsAdmin struct {
	UserName    string `db:"username" json:"username"`
	UserID      int64  `db:"user_id" json:"user_id"`
	NamespaceID int64  `db:"namespace_id" json:"namespace_id"`
}

// User : base model.
type User struct {
	Name         string  `json:"username"`
	PasswordHash string  `json:"-"`
	Namespace    string  `json:"namespace,omitempty"`
	Org          string  `json:"org,omitempty"`
	Groups       []Group `json:"groups,omitempty"`
	ID           int64   `json:"id"`
}

// GroupNames : get the user group names.
func (user User) GroupNames() []string {
	names := make([]string, 0, len(user.Groups))
	for _, g := range user.Groups {
		names = append(names, g.Name)
	}
	return names
}

// Group : base model.
type Group struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	ID        int64  `json:"id"`
}

// UserGroup : base model.
type UserGroup struct {
	ID      int32 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}

// Org : base model.
type Org struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}
