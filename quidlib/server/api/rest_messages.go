package api

//go:generate go run github.com/mailru/easyjson/... -all -snake_case ${GOFILE}

type LoginRequest struct {
	Username  string
	Password  string
	Namespace string
}

type AdminAccessTokenRequest struct {
	RefreshToken string
	Namespace    string
}

type GroupCreation struct {
	Name        string
	NamespaceID int64
}

type NamespaceRefreshTokenMaxTTLRequest struct {
	ID            int64
	RefreshMaxTTL string
}

type NamespaceTokenMaxTTLRequest struct {
	ID     int64
	MaxTTL string
}

type InfoRequest struct {
	ID int64
}

type NameRequest struct {
	Name string
}

type Availability struct {
	ID     int64
	Enable bool
}

type NamespaceCreation struct {
	Name           string
	MaxTTL         string
	RefreshMaxTTL  string
	EnableEndpoint bool
}

type NonAdminUsersRequest struct {
	Username    string
	NamespaceID int64
}

type AdministratorsCreation struct {
	UserIDs     []int64
	NamespaceID int64
}

type AdministratorDeletion struct {
	UserID      int64
	NamespaceID int64
}

type AccessTokenRequest struct {
	RefreshToken string
	Namespace    string
}

type RefreshTokenRequest struct {
	Username  string
	Password  string
	Namespace string
}

type NamespaceRequest struct {
	Namespace string
}

type UserOrgRequest struct {
	UserID int64
	OrgID  int64
}

type UserGroupRequest struct {
	UserID      int64
	GroupID     int64
	NamespaceID int64
}

type UserRequest struct {
	ID          int64
	NamespaceID int64
}

type UserHandlerCreation struct {
	Name        string
	Password    string
	NamespaceID int64
}
