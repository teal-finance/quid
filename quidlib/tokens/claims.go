package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AccessClaims : standard claims for a user access token
type AccessClaims struct {
	UserName string   `json:"username"`
	Groups   []string `json:"groups"`
	Orgs     []string `json:"orgs"`
	jwt.StandardClaims
}

// RefreshClaims : standard claims for a user refresh token
type RefreshClaims struct {
	Namespace string `json:"namespace"`
	UserName  string `json:"username"`
	jwt.StandardClaims
}

// newAccessClaims : get a standard claim for a user access token
func newAccessClaims(username string, groups, orgs []string, timeout time.Time) *AccessClaims {
	return &AccessClaims{
		username, groups, orgs,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
}

// newRefreshClaims : get a standard claim for a user refresh token
func newRefreshClaims(namespaceName, username string, timeout time.Time) *RefreshClaims {
	return &RefreshClaims{
		namespaceName, username,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
}
