package server

type PasswordRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
}

type UserHandlerCreation struct {
	Username string `json:"username"`
	Password string `json:"password"`
	NsID     int64  `json:"ns_id"`
}

type UserSetEnabled struct {
	ID      int64 `json:"id"`
	Enabled bool  `json:"enabled"`
}

type GroupCreation struct {
	Name string `json:"name"`
	NsID int64  `json:"ns_id"`
}

type NamespaceIDRequest struct {
	NsID int64 `json:"ns_id"`
}

type RefreshMaxTTLRequest struct {
	ID            int64  `json:"id"`
	RefreshMaxTTL string `json:"refresh_max_ttl"`
}

type MaxTTLRequest struct {
	ID     int64  `json:"id"`
	MaxTTL string `json:"max_ttl"`
}

type InfoRequest struct {
	ID           int64  `json:"id"`
	EncodingForm string `json:"encoding_form"`
}

type NameRequest struct {
	Name string `json:"name"`
}

type Availability struct {
	ID     int64 `json:"id"`
	Enable bool
}

type NamespaceCreation struct {
	Name           string `json:"name"`
	Alg            string `json:"alg"`
	MaxTTL         string `json:"max_ttl"`
	RefreshMaxTTL  string `json:"refresh_max_ttl"`
	EnableEndpoint bool   `json:"enable_endpoint"`
}

type NonAdminUsersRequest struct {
	Username string `json:"username"`
	NsID     int64  `json:"ns_id"`
}

type AdministratorsCreation struct {
	UsrIDs []int64 `json:"usr_ids"`
	NsID   int64   `json:"ns_id"`
}

type AdministratorDeletion struct {
	UsrID int64 `json:"usr_id"`
	NsID  int64 `json:"ns_id"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Namespace    string `json:"namespace"`
}

type AccessTokenValidationRequest struct {
	AccessToken string `json:"access_token"`
	Namespace   string `json:"namespace"`
}

type NamespaceRequest struct {
	Namespace    string `json:"namespace"`
	EncodingForm string `json:"encoding_form"`
}

type UserOrgRequest struct {
	UsrID int64 `json:"usr_id"`
	OrgID int64 `json:"org_id"`
}

type UserGroupRequest struct {
	UsrID int64 `json:"usr_id"`
	GrpID int64 `json:"grp_id"`
	NsID  int64 `json:"ns_id"`
}

type UserRequest struct {
	ID   int64 `json:"id"`
	NsID int64 `json:"ns_id"`
}

type PublicKeyResponse struct {
	Alg string `json:"alg"`
	Key []byte `json:"key"`
}
