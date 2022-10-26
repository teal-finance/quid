package server

import "time"

//go:generate go run github.com/mailru/easyjson/... -all -byte -disable_members_unescape -disallow_unknown_fields -snake_case ${GOFILE}

// Namespace : base model.
type Namespace struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Alg           string `json:"alg"`
	RefreshKey    []byte `json:"-"`
	AccessKey     []byte `json:"-"`
	MaxRefreshTTL string `json:"max_refresh_ttl"`
	MaxAccessTTL  string `json:"max_access_ttl"`
	Enabled       bool   `json:"public_endpoint_enabled"`
}

// NamespaceInfo : base model.
type NamespaceInfo struct {
	Groups   []Group `json:"groups"`
	NumUsers int     `json:"num_users"`
}

// Administrator : base model.
type Administrator struct {
	ID    int64  `json:"id"`
	Name  string `json:"username"`
	UsrID int64  `json:"usr_id"`
	NsID  int64  `json:"ns_id"`
}

// NonAdmin : base model.
type NonAdmin struct {
	Name  string `json:"username"`
	UsrID int64  `json:"usr_id"`
	NsID  int64  `json:"ns_id"`
}

// User : base model.
type User struct {
	DateCreated  time.Time `json:"date_created"`
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Namespace    string    `json:"namespace,omitempty"`
	Org          string    `json:"org,omitempty"`
	Groups       []Group   `json:"groups,omitempty"`
	Enabled      bool      `json:"enabled"`
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
	AdminType string
	Username  string
	Ns        NSInfo
}

type NSInfo struct {
	ID   int64
	Name string
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
