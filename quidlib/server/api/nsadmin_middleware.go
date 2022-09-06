package api

import (
	"net/http"

	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// isNsAdmin checks that the requested namespace operation
// matches the request ns admin permissions
func isNsAdmin(r *http.Request, nsID int64) bool {
	tv, ok := incorruptible.FromCtx(r)
	if !ok {
		log.ParamError("VerifyAdminNs: missing Incorruptible token: cannot check Admin nsID=", nsID)
		return false
	}

	adminType := AdminType(tv.StringIfAny(keyAdminType))
	if adminType == QuidAdmin {
		log.Param("VerifyAdminNs OK: Incorruptible token contains IsNsAdmin=true => Do not check the nsID")
		return true
	}

	gotID, err := tv.Int64(keyNsID)
	if err != nil {
		log.ParamError("VerifyAdminNs: missing field nsID in Incorruptible token, want nsID=", nsID, err)
		return false
	}
	if gotID != nsID {
		log.ParamError("VerifyAdminNs: user is nsAdmin for", gotID, ", but not", nsID)
		return false
	}

	return true
}

// nsAdminMiddleware : check the token claim to see if the user is namespace admin.
func nsAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tv, err := Incorruptible.DecodeCookieToken(r)
		if err != nil {
			log.Warn("QuidAdminMiddleware: no valid token:", err.Error())
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
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := values[0].String()
		userID := values[1].Int64()
		namespace := values[2].String()
		nsID := values[3].Int64()

		userType, err := db.GetUserType(namespace, nsID, userID)
		if err != nil {
			log.QueryError("NsAdminMiddleware:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userType == db.UserNoAdmin {
			log.ParamError("NsAdminMiddleware: u=" + userName + " is admin, but not for ns=" + namespace)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.RequestPost("NsAdminMiddleware OK u="+userName+" (id=", userID, ") ns="+namespace+" (id=", nsID, ")")
		r = tv.ToCtx(r) // save the token in the request context
		next.ServeHTTP(w, r)
	})
}
