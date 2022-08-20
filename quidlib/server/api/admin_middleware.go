package api

import (
	"net/http"

	"github.com/teal-finance/quid/quidlib/server/db"
)

// AdminMiddleware : check the token claim to see if the user is admin.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, err := Incorruptible.DecodeCookieToken(r)
		if err != nil {
			emo.Warning("AdminMiddleware: no valid token:", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		values, err := tv.Get(
			tv.KString(keyUsername),
			tv.KInt64(KeyUserID),
			tv.KString(keyNsName),
			tv.KInt64(keyNsID),
			tv.KBool(keyAdminType))
		if err != nil {
			emo.Error("AdminMiddleware tv.Get:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := values[0].String()
		userID := values[1].Int64()
		namespace := values[2].String()
		nsID := values[3].Int64()
		isAdmin := values[4].Bool()

		if !isAdmin {
			emo.ParamError("AdminMiddleware: User " + userName + " is not Admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err = db.IsUserAdmin(namespace, nsID, userID)
		if err != nil {
			emo.QueryError("AdminMiddleware IsUserAdmin:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.Data("AdminMiddleware: u=" + userName + " is not Admin in database")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		emo.Param("AdminMiddleware OK u="+userName+" (id=", userID, ") ns="+namespace+" (id=", nsID, ")")
		r = tv.ToCtx(r) // save the token in the request context
		next.ServeHTTP(w, r)
	})
}
