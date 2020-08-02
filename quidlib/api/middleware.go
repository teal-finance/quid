package api

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo-contrib/session"
	"github.com/synw/quid/quidlib/conf"
	"github.com/synw/quid/quidlib/tokens"

	"github.com/labstack/echo/v4"
)

// AdminMiddleware : check the token claim to see if the user is admin
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check the access token to control that the user is admin
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*tokens.StandardAccessClaims)
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
		// check session in production
		if conf.IsDevMode {
			return next(c)
		}
		sess, _ := session.Get("session", c)
		fmt.Println("IS ADMIN", sess.Values["is_admin"])
		fmt.Println("USER", sess.Values["user"])
		if sess.Values["is_admin"] == "true" {
			return next(c)
		}
		return c.NoContent(http.StatusUnauthorized)
	}
}
