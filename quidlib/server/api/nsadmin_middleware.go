package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

func VerifyAdminNs(c echo.Context, nsID int64) bool {
	// check that the requested namespace operation
	// matches the request ns admin permissions
	isAdmin := c.Get("isAdmin").(bool)
	adminNs := c.Get("isAdminForNs").(int64)
	if isAdmin {
		return true
	} else {
		if adminNs == nsID {
			return true
		}
	}
	emo.ParamError("User is not nsadmin for namespace", nsID, "/", adminNs)
	return false
}

// NsAdminMiddleware : check the token claim to see if the user is namespaace admin.
func NsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		claims, ok := u.Claims.(*tokens.AdminAccessClaim)
		if !ok {
			emo.Error("Wrong AccessClaims type for claims: ", u.Claims)
			log.Panic("Wrong AccessClaims type for claims: ", u.Claims)
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "Wrong AccessClaims type for claims",
			})
		}

		emo.Data("Admin claims for", claims.Namespace, claims)

		isAdmin, err := db.IsUserAdmin(claims.Namespace, claims.NsID, claims.UserID)
		if err != nil {
			return err
		}
		if !isAdmin {
			emo.ParamError("The user "+claims.UserName+" is not nsadmin for namespace", claims.Namespace)
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("isAdmin", false)
		c.Set("isAdminForNs", claims.NsID)

		// check session data in production
		if conf.IsDevMode {
			emo.RequestPost("Request ok from nsadmin middleware")
			return next(c)
		}

		sess, _ := session.Get("session", c)
		if sess.Values["is_nsadmin"] == "true" {
			return next(c)
		}

		emo.Warning("Unauthorized session from nsadmin middleware")
		return c.NoContent(http.StatusUnauthorized)
	}
}
