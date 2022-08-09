package api

import (
	"net/http"

	"github.com/teal-finance/incorruptible/tvalues"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// AdminMiddleware : check the token claim to see if the user is admin.
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, ok := tvalues.FromCtx(r)
		if !ok {
			emo.Error("Missing Incorruptible token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		values, err := tv.Get(
			tv.KString(user),
			tv.KInt64(user_id),
			tv.KString(ns_name),
			tv.KInt64(ns_id),
			tv.KBool(is_admin))
		if err != nil {
			emo.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := values[0].String()
		userID := values[1].Int64()
		namespace := values[2].String()
		nsID := values[3].Int64()
		isAdmin := values[4].Bool()

		if !isAdmin {
			emo.ParamError("User " + userName + "is not Admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err = db.IsUserAdmin(namespace, nsID, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.ParamError("User " + userName + " is not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
		return
	})
}
