package api

import (
	"net/http"

	"github.com/teal-finance/incorruptible/tvalues"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// VerifyAdminNs checks that the requested namespace operation
// matches the request ns admin permissions
func VerifyAdminNs(w http.ResponseWriter, r *http.Request, nsID int64) bool {
	tv, ok := tvalues.FromCtx(r)
	if !ok {
		emo.ParamError("Missing Incorruptible token: cannot check Admin NsID=", nsID)
		return false
	}

	nsAdmin, err := tv.Bool(is_ns_admin)
	if nsAdmin && (err == nil) {
		return true
	}

	gotID, err := tv.Int64(ns_id)
	if err != nil {
		emo.ParamError("Missing field 'ns_id' in Incorruptible token, want nsID=", nsID)
		return false
	}
	if gotID != nsID {
		emo.ParamError("User is nsAdmin for ", gotID, " but not ", nsID)
		return false
	}

	return true
}

// NsAdminMiddleware : check the token claim to see if the user is namespace admin.
func NsAdminMiddleware(next http.Handler) http.Handler {
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
			tv.KBool(is_ns_admin))
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
			emo.Error("User " + userName + "is not a ns admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err = db.IsUserAdmin(namespace, nsID, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.ParamError("The user "+userName+" is not admin for namespace", namespace)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		emo.RequestPost("Request ok from ns admin middleware")
		next.ServeHTTP(w, r)
	})
}
