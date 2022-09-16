package server

//go:generate go run github.com/mailru/easyjson/... -all -byte -disable_members_unescape -disallow_unknown_fields -snake_case ${GOFILE}

/* Documentation of the above generator command line:

$ go run github.com/mailru/easyjson/easyjson -help

Usage of /tmp/go-build1862268821/b001/exe/easyjson:
  -all
        generate marshaler/unmarshalers for all structs in a file
  -build_tags string
        build tags to add to generated file
  -byte
        use simple bytes instead of Base64Bytes for slice of bytes
  -disable_members_unescape
        don't perform unescaping of member names to improve performance
  -disallow_unknown_fields
        return error if any unknown field in json appeared
  -gen_build_flags string
        build flags when running the generator while bootstrapping
  -leave_temps
        do not delete temporary files
  -lower_camel_case
        use lowerCamelCase names instead of CamelCase by default
  -no_std_marshalers
        don't generate MarshalJSON/UnmarshalJSON funcs
  -noformat
        do not run 'gofmt -w' on output file
  -omit_empty
        omit empty fields by default
  -output_filename string
        specify the filename of the output
  -pkg
        process the whole package instead of just the given file
  -snake_case
        use snake_case names instead of CamelCase by default
  -stubs
        only generate stubs for marshaler/unmarshaler funcs
*/

type PasswordRequest struct {
	Username  string
	Password  string
	Namespace string
}

type UserHandlerCreation struct {
	Name        string
	Password    string
	NamespaceID int64
}

type GroupCreation struct {
	Name        string
	NamespaceID int64
}

type NamespaceIDRequest struct {
	NamespaceID int64
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
	Algo           string
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

type AccessTokenValidationRequest struct {
	AccessToken string
	Namespace   string
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

type StatusResponse struct {
	AdminType string
	Username  string
	Ns        NSInfo
}

type NSInfo struct {
	ID   int64
	Name string
}

type PublicKeyResponse struct {
	Alg string
	Key string
}
