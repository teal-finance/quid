package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AccessClaims is the standard claims for a user access token.
type AccessClaims struct {
	jwt.RegisteredClaims
	UserName string   `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Orgs     []string `json:"orgs,omitempty"`
}

// RefreshClaims is the standard claims for a user refresh token.
type RefreshClaims struct {
	Namespace string `json:"namespace,omitempty"`
	UserName  string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

// newAccessClaims creates a standard claim for a user access token.
func newAccessClaims(username string, groups, orgs []string, timeout time.Time) AccessClaims {
	return AccessClaims{
		jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(timeout),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
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
		jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(timeout),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}
}
