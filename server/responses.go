package server

import "time"

// Namespace : base model.
type Namespace struct {
	ID            int64  `json:"id"                      db:"id"`
	Name          string `json:"name"                    db:"name"`
	Alg           string `json:"alg"                     db:"alg"`
	RefreshKey    []byte `json:"-"                       db:"refresh_key"`
	AccessKey     []byte `json:"-"                       db:"access_key"`
	MaxRefreshTTL string `json:"max_refresh_ttl"         db:"max_refresh_ttl"`
	MaxAccessTTL  string `json:"max_access_ttl"          db:"max_access_ttl"`
	Enabled       bool   `json:"public_endpoint_enabled" db:"public_endpoint_enabled"`
}

// NamespaceInfo : base model.
type NamespaceInfo struct {
	Groups   []Group `json:"groups"`
	NumUsers int     `json:"num_users"`
}

// User : base model.
type User struct {
	DateCreated  time.Time `json:"date_created"        db:"date_created"`
	ID           int64     `json:"id"                  db:"id"`
	Name         string    `json:"name"                db:"name"`
	PasswordHash string    `json:"-"                   db:"password_hash"`
	Namespace    string    `json:"namespace,omitempty" db:"namespace"`
	Org          string    `json:"org,omitempty"       db:"org"`
	Groups       []Group   `json:"groups,omitempty"    db:"groups"`
	Enabled      bool      `json:"enabled"             db:"enabled"`
}

// groupNames : list the names of all groups.
// Deprecated because this function is not used.
func (user User) groupNames() []string {
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

// // UserGroup : base model.
// type UserGroup struct {
// 	ID    int32 `json:"id"`
// 	UsrID int64 `json:"usr_id"`
// 	GrpID int64 `json:"grp_id"`
// }

// Org : base model.
type Org struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

type StatusResponse struct {
	AdminType string `json:"admin_type"`
	Username  string `json:"username"`
	Ns        NSInfo `json:"ns"`
}

type NSInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AdminType bool

const (
	QuidAdmin AdminType = false
	NsAdmin   AdminType = true
)

func (t AdminType) String() string {
	if t == QuidAdmin {
		return "QuidAdmin"
	}
	return "NsAdmin"
}
