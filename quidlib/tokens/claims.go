package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// StandardAccessClaims : standard claims for a user access token
type StandardAccessClaims struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	jwt.StandardClaims
}

// StandardRefreshClaims : standard claims for a user refresh token
type StandardRefreshClaims struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	jwt.StandardClaims
}

// standardAccessClaims : get a standard claim for a user access token
func standardAccessClaims(namespaceName, name string, groups []string, timeout time.Time) *StandardAccessClaims {
	claims := &StandardAccessClaims{
		namespaceName, name, groups,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}

// standardRefreshClaims : get a standard claim for a user refresh token
func standardRefreshClaims(namespaceName, name string, timeout time.Time) *StandardRefreshClaims {
	claims := &StandardRefreshClaims{
		namespaceName, name,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}
