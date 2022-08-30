package api

import (
	"net/http"

	"github.com/teal-finance/quid/quidlib/server/db"
)

// QuidAdminMiddleware : check the token claim to see if the user is admin.
func QuidAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, err := Incorruptible.DecodeCookieToken(r)
		if err != nil {
			logg.Warning("QuidAdminMiddleware: no valid token:", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		values, err := tv.Get(
			tv.KString(keyUsername),
			tv.KInt64(KeyUserID),
			tv.KString(keyNsName),
			tv.KInt64(keyNsID),
			tv.KString(keyAdminType))
		if err != nil {
			logg.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := values[0].String()
		userID := values[1].Int64()
		namespace := values[2].String()
		nsID := values[3].Int64()
		adminType := values[4].String()

		if AdminType(adminType) != QuidAdmin {
			logg.ParamError("User '" + userName + "' is not QuidAdmin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userType, err := db.GetUserType(namespace, nsID, userID)
		if err != nil {
			logg.QueryError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userType != db.QuidAdmin {
			logg.Data("QuidAdminMiddleware: u=" + userName + " is not Admin in database")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logg.Param("QuidAdminMiddleware OK u="+userName+" (id=", userID, ") ns="+namespace+" (id=", nsID, ")")
		r = tv.ToCtx(r) // save the token in the request context
		next.ServeHTTP(w, r)
	})
}
