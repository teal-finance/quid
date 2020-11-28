package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// StandardAccessClaims : standard claims for a user access token
type StandardAccessClaims struct {
	UserName string   `json:"username"`
	Groups   []string `json:"groups"`
	Orgs     []string `json:"orgs"`
	jwt.StandardClaims
}

// StandardRefreshClaims : standard claims for a user refresh token
type StandardRefreshClaims struct {
	Namespace string `json:"namespace"`
	UserName  string `json:"username"`
	jwt.StandardClaims
}

// standardAccessClaims : get a standard claim for a user access token
func standardAccessClaims(username string, groups, orgs []string, timeout time.Time) *StandardAccessClaims {
	claims := &StandardAccessClaims{
		username, groups, orgs,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}

// standardRefreshClaims : get a standard claim for a user refresh token
func standardRefreshClaims(namespaceName, username string, timeout time.Time) *StandardRefreshClaims {
	claims := &StandardRefreshClaims{
		namespaceName, username,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}
