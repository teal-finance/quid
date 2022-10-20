package server

type PasswordRequest struct {
	Username  string
	Password  string
	Namespace string
}

type UserHandlerCreation struct {
	Name     string
	Password string
	NsID     int64
}

type UserSetEnabled struct {
	ID      int64
	Enabled bool
}

type GroupCreation struct {
	Name string
	NsID int64
}

type NamespaceIDRequest struct {
	NsID int64
}

type RefreshMaxTTLRequest struct {
	ID            int64
	RefreshMaxTTL string
}

type MaxTTLRequest struct {
	ID     int64
	MaxTTL string
}

type InfoRequest struct {
	ID           int64
	EncodingForm string
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
	Alg            string
	MaxTTL         string
	RefreshMaxTTL  string
	EnableEndpoint bool
}

type NonAdminUsersRequest struct {
	Username string
	NsID     int64
}

type AdministratorsCreation struct {
	UsrIDs []int64 `json:"usr_ids"`
	NsID   int64
}

type AdministratorDeletion struct {
	UsrID int64
	NsID  int64
}

type AccessTokenRequest struct {
	RefreshToken string
	Namespace    string
}

type AccessTokenValidationRequest struct {
	AccessToken string
	Namespace   string
}

type NamespaceRequest struct {
	Namespace    string
	EncodingForm string
}

type UserOrgRequest struct {
	UsrID int64
	OrgID int64
}

type UserGroupRequest struct {
	UsrID int64
	GrpID int64
	NsID  int64
}

type UserRequest struct {
	ID   int64
	NsID int64
}

type PublicKeyResponse struct {
	Alg string
	Key []byte
}
