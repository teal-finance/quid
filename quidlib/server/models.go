package server

// Namespace : base model
type Namespace struct {
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	Key                   string `db:"key" json:"-"`
	RefreshKey            string `db:"refresh_key" json:"-"`
	MaxTokenTTL           string `db:"max_token_ttl" json:"max_token_ttl"`
	MaxRefreshTokenTTL    string `db:"max_refresh_token_ttl" json:"max_refresh_token_ttl"`
	PublicEndpointEnabled bool   `db:"public_endpoint_enabled" json:"public_endpoint_enabled"`
}

// NamespaceInfo : base model
type NamespaceInfo struct {
	NumUsers int     `json:"num_users"`
	Groups   []Group `json:"groups"`
}

// User : base model
type User struct {
	ID           int64   `json:"id"`
	UserName     string  `json:"username"`
	PasswordHash string  `json:"-"`
	Groups       []Group `json:"groups,omitempty"`
	Namespace    string  `json:"namespace,omitempty"`
	Org          string  `json:"org,omitempty"`
}

// GroupNames : get the user group names
func (user User) GroupNames() []string {
	u := []string{}
	for _, g := range user.Groups {
		u = append(u, g.Name)
	}
	return u
}

// Group : base model
type Group struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// UserGroup : base model
type UserGroup struct {
	ID      int32 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"user_id"`
	GroupID int64 `db:"group_id" json:"group_id"`
}

// Org : base model
type Org struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
