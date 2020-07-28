package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// StandardUserClaims : standard claims for a user
type StandardUserClaims struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	jwt.StandardClaims
}

// SetStandardClaim : get a standard claim for a user
func standardUserClaims(namespaceName, name string, groups []string, timeout time.Time) *StandardUserClaims {
	claims := &StandardUserClaims{
		namespaceName, name, groups,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}
