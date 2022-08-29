package tokens

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

//go:generate go run github.com/mailru/easyjson/... -all -byte -disable_members_unescape -disallow_unknown_fields ${GOFILE}

// AccessClaims is the standard claims for a user access token.
type AccessClaims struct {
	UserName string   `json:"usr,omitempty"`
	Groups   []string `json:"grp,omitempty"`
	Orgs     []string `json:"org,omitempty"`
	jwt.RegisteredClaims
}

type AdminAccessClaim struct {
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
	NsID      int64  `json:"ns_id,omitempty"`
	IsAdmin   bool   `json:"is_admin"`
	IsNsAdmin bool   `json:"is_ns_admin"`
	jwt.RegisteredClaims
}

// RefreshClaims is the standard claims for a user refresh token.
type RefreshClaims struct {
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

// newAdminAccessClaims creates a standard claim for an admin user access token.
func newAdminAccessClaims(namespaceName, username string, userID, nsID int64, expiry time.Time, isAdmin, isNsAdmin bool) AdminAccessClaim {
	return AdminAccessClaim{
		namespaceName,
		username,
		userID,
		nsID,
		isAdmin,
		isNsAdmin,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiry)},
	}
}

// newAccessClaims creates a standard claim for a user access token.
func newAccessClaims(username string, groups, orgs []string, expiry time.Time) AccessClaims {
	return AccessClaims{
		username,
		groups,
		orgs,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiry)},
	}
}

// newRefreshClaims creates a standard claim for a user refresh token.
func newRefreshClaims(namespace, user string, expiry time.Time) RefreshClaims {
	return RefreshClaims{
		namespace,
		user,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiry)},
	}
}
