package api

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/tokens"

	"github.com/labstack/echo/v4"
)

// AdminMiddleware : check the token claim to see if the user is admin.
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check the access token to control that the user is admin
		u, ok := c.Get("user").(*jwt.Token)
		if !ok {
			emo.Error("Wrong JWT type in context user: ", c.Get("user"))
			log.Panicf("Wrong JWT type in context user: %T %v", c.Get("user"), c.Get("user"))
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "Wrong JWT type in context user",
			})
		}

		claims, ok := u.Claims.(*tokens.AccessClaims)
		if !ok {
			emo.Error("Wrong AccessClaims type for claims: ", u.Claims)
			log.Panic("Wrong AccessClaims type for claims: ", u.Claims)
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "Wrong AccessClaims type for claims",
			})
		}

		isAdmin := false
		for _, g := range claims.Groups {
			if g == "quid_admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			emo.ParamError("The user " + claims.UserName + " is not in the quid_admin group")
			return c.NoContent(http.StatusUnauthorized)
		}

		// check session data in production
		if conf.IsDevMode {
			return next(c)
		}

		sess, _ := session.Get("session", c)
		if sess.Values["is_admin"] == "true" {
			return next(c)
		}

		emo.Warning("Unauthorized session from admin middleware")
		return c.NoContent(http.StatusUnauthorized)
	}
}
