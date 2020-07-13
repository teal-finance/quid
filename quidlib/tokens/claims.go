package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// StandardUserClaims : standard claims for a user
type StandardUserClaims struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
	jwt.StandardClaims
}

// SetStandardClaim : get a standard claim for a user
func standardUserClaims(name string, groups []string, timeout time.Time) *StandardUserClaims {
	claims := &StandardUserClaims{
		name, groups,
		jwt.StandardClaims{
			ExpiresAt: timeout.Unix(),
		},
	}
	return claims
}
