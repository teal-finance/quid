package api

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"

	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

func VerifyAdminNs(w http.ResponseWriter, r *http.Request, nsID int64) bool {
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

// NsAdminMiddleware : check the token claim to see if the user is namespace admin.
func NsAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the access token to control that the user is admin
		u, ok := c.Get("user").(*jwt.Token)
		if !ok {
			emo.Error("Wrong JWT type in context user: ", c.Get("user"))
			log.Panicf("Wrong JWT type in context user: %T %v", c.Get("user"), c.Get("user"))
			gw.WriteErr(w, r, http.StatusConflict, "Wrong JWT type in context user")
			return
		}

		claims, ok := u.Claims.(*tokens.AdminAccessClaim)
		if !ok {
			emo.Error("Wrong AccessClaims type for claims: ", u.Claims)
			log.Panic("Wrong AccessClaims type for claims: ", u.Claims)
			gw.WriteErr(w, r, http.StatusConflict, "Wrong AccessClaims type for claims")
			return
		}

		emo.Data("Admin claims for", claims.Namespace, claims)

		isAdmin, err := db.IsUserAdmin(claims.Namespace, claims.NsID, claims.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.ParamError("The user "+claims.UserName+" is not nsadmin for namespace", claims.Namespace)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c.Set("isAdmin", false)
		c.Set("isAdminForNs", claims.NsID)

		// check session data in production
		if conf.IsDevMode {
			emo.RequestPost("Request ok from nsadmin middleware")
			next.ServeHTTP(w, r)
			return
		}

		sess, _ := session.Get("session", c)
		if sess.Values["is_ns_admin"] == "true" {
			next.ServeHTTP(w, r)
			return
		}

		emo.Warning("Unauthorized session from nsadmin middleware")
		w.WriteHeader(http.StatusUnauthorized)
	})
}
