package api

import (
	"net/http"

	"github.com/teal-finance/quid/server"
	"github.com/teal-finance/quid/server/db"
)

// quidAdminMiddleware : check the token claim to see if the user is admin.
func quidAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, err := Incorruptible.DecodeCookieToken(r)
		if err != nil {
			log.Warn("quidAdminMiddleware wants cookie", Incorruptible.Cookie(0).Name, "but", err)
			gw.WriteErr(w, r, http.StatusUnauthorized, "missing or invalid incorruptible cookie", "want_cookie_name", Incorruptible.Cookie(0).Name)
			return
		}

		values, err := tv.Get(
			tv.KString(keyUsername),
			tv.KInt64(KeyUserID),
			tv.KString(keyNsName),
			tv.KInt64(keyNsID),
			tv.KString(keyAdminType))
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := values[0].String()
		userID := values[1].Int64()
		namespace := values[2].String()
		nsID := values[3].Int64()
		adminType := values[4].Bool()

		if server.AdminType(adminType) != server.QuidAdmin {
			log.ParamError("User '" + userName + "' is not QuidAdmin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userType, err := db.GetUserType(namespace, nsID, userID)
		if err != nil {
			log.QueryError(err)
			gw.WriteErr(w, r, http.StatusUnauthorized, "DB error while getting user type", "namespace_id", nsID, "uid", userID)
			return
		}
		if userType != db.QuidAdmin {
			log.Data("quidAdminMiddleware: u=" + userName + " is not Admin in database")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Param("quidAdminMiddleware OK u="+userName+" (id=", userID, ") ns="+namespace+" (id=", nsID, ")")
		r = tv.ToCtx(r) // save the token in the request context
		next.ServeHTTP(w, r)
	})
}
