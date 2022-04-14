package tokens

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// AccessClaims is the standard claims for a user access token.
type AccessClaims struct {
	jwt.StandardClaims
	UserName string   `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Orgs     []string `json:"orgs,omitempty"`
}

// RefreshClaims is the standard claims for a user refresh token.
type RefreshClaims struct {
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	jwt.StandardClaims
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
