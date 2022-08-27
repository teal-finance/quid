package tokens

import (
	"time"

	"github.com/golang-jwt/jwt"
)

//go:generate go run github.com/mailru/easyjson/... -all -byte -disable_members_unescape -disallow_unknown_fields ${GOFILE}

// AccessClaims is the standard claims for a user access token.
type AccessClaims struct {
	jwt.StandardClaims
	UserName string   `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Orgs     []string `json:"orgs,omitempty"`
}

type AdminAccessClaim struct {
	jwt.StandardClaims
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
	NsID      int64  `json:"ns_id,omitempty"`
	IsAdmin   bool   `json:"is_admin"`
	IsNsAdmin bool   `json:"is_ns_admin"`
}

// RefreshClaims is the standard claims for a user refresh token.
type RefreshClaims struct {
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	jwt.StandardClaims
}

// newAdminAccessClaims creates a standard claim for an admin user access token.
func newAdminAccessClaims(namespaceName, username string, userID, nsID int64, timeout time.Time, isAdmin, isNsAdmin bool) AdminAccessClaim {
	return AdminAccessClaim{
		jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: timeout.Unix(),
			Id:        "",
			IssuedAt:  0,
			Issuer:    "",
			NotBefore: 0,
			Subject:   "",
		},
		namespaceName,
		username,
		userID,
		nsID,
		isAdmin,
		isNsAdmin,
	}
}

// newAccessClaims creates a standard claim for a user access token.
func newAccessClaims(username string, groups, orgs []string, timeout time.Time) AccessClaims {
	return AccessClaims{
		jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: timeout.Unix(),
			Id:        "",
			IssuedAt:  0,
			Issuer:    "",
			NotBefore: 0,
			Subject:   "",
		},
		username,
		groups,
		orgs,
	}
}

// newRefreshClaims creates a standard claim for a user refresh token.
func newRefreshClaims(namespace, user string, timeout time.Time) RefreshClaims {
	return RefreshClaims{
		namespace,
		user,
		jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: timeout.Unix(),
			Id:        "",
			IssuedAt:  0,
			Issuer:    "",
			NotBefore: 0,
			Subject:   "",
		},
	}
}
