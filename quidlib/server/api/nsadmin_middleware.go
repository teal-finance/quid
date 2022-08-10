package api

import (
	"net/http"

	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// VerifyAdminNs checks that the requested namespace operation
// matches the request ns admin permissions
func VerifyAdminNs(w http.ResponseWriter, r *http.Request, nsID int64) bool {
	tv, ok := incorruptible.FromCtx(r)
	if !ok {
		emo.ParamError("VerifyAdminNs: missing Incorruptible token: cannot check Admin NsID=", nsID)
		return false
	}

	nsAdmin, err := tv.Bool(keyIsNsAdmin)
	if nsAdmin && (err == nil) {
		return true
	}

	gotID, err := tv.Int64(keyNsID)
	if err != nil {
		emo.ParamError("VerifyAdminNs: missing field 'ns_id' in Incorruptible token, want nsID=", nsID)
		return false
	}
	if gotID != nsID {
		emo.ParamError("VerifyAdminNs: user is nsAdmin for", gotID, " but not", nsID)
		return false
	}

	return true
}

// NsAdminMiddleware : check the token claim to see if the user is namespace admin.
func NsAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, ok := incorruptible.FromCtx(r)
		if !ok {
			emo.Error("NsAdminMiddleware: missing Incorruptible token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		values, err := tv.Get(
			tv.KString(keyUserName),
			tv.KInt64(KeyUserID),
			tv.KString(keyNsName),
			tv.KInt64(keyNsID),
			tv.KBool(keyIsNsAdmin))
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
			emo.ParamError("NsAdminMiddleware: user " + userName + "is not ns admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		isAdmin, err = db.IsUserAdmin(namespace, nsID, userID)
		if err != nil {
			emo.Error("NsAdminMiddleware:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			emo.ParamError("NsAdminMiddleware: user "+userName+" is not admin for namespace", namespace)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		emo.RequestPost("Request ok from ns admin middleware")
		next.ServeHTTP(w, r)
	})
}
