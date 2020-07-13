package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/synw/quid/quidlib/tokens"
)

// AdminMiddleware : check the token claim to see if the user is admin
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*tokens.StandardUserClaims)
		isAdmin := false
		for _, g := range claims.Groups {
			if g == "quid_admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			return c.NoContent(http.StatusUnauthorized)
		}
		return nil
	}
}
