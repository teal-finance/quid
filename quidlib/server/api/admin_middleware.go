package api

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

var c echo.Context // FIXME TODO

// AdminMiddleware : check the token claim to see if the user is admin.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the access token to control that the user is admin
		u, ok := c.Get("user").(*jwt.Token)
		if !ok {
			emo.Error("Wrong JWT type in context user: ", c.Get("user"))
			log.Panicf("Wrong JWT type in context user: %T %v", c.Get("user"), c.Get("user"))
			gw.WriteErr(w, r, http.StatusConflict, "Wrong JWT type in context user")
		}

		claims, ok := u.Claims.(*tokens.AdminAccessClaim)
		if !ok {
			emo.Error("Wrong AccessClaims type for claims: ", u.Claims)
			log.Panic("Wrong AccessClaims type for claims: ", u.Claims)
			gw.WriteErr(w, r, http.StatusConflict, "Wrong AccessClaims type for claims")
		}

		emo.Data("Admin claims for", claims.Namespace, claims)

		if claims.Namespace != "quid" {
			// The user is only namespace admin
			emo.ParamError("The user " + claims.UserName + " is not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err := db.IsUserAdmin(claims.Namespace, claims.NsID, claims.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.ParamError("The user " + claims.UserName + " is not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = PutAdminInfoInCtx(r, true, 0)

		// check session data in production
		if conf.IsDevMode {
			next.ServeHTTP(w, r)
			return
		}

		sess, _ := session.Get("session", c)
		if sess.Values["is_admin"] == "true" {
			next.ServeHTTP(w, r)
			return
		}

		emo.Warning("Unauthorized session from admin middleware")
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// --------------------------------------
// Read/write admin info to/from context

//nolint:gochecknoglobals // adminKey is a Context key and need to be global
var adminKey struct{}

type AdminInfo struct {
	isAdmin     bool
	namespaceID int64
}

// GetAdminInfoFromCtx gets the admin information from the request context.
func GetAdminInfoFromCtx(r *http.Request) AdminInfo {
	info, ok := r.Context().Value(adminKey).(AdminInfo)
	if !ok {
		log.Print("WRN JWT: no permission in context ", r.URL.Path)
	}
	return info
}

// PutAdminInfoInCtx stores the admin information within the request context.
func PutAdminInfoInCtx(r *http.Request, isAdmin bool, namespaceID int64) *http.Request {
	info := AdminInfo{isAdmin, namespaceID}
	parent := r.Context()
	child := context.WithValue(parent, adminKey, info)
	return r.WithContext(child)
}
