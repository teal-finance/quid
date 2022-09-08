package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/server/db"
)

// Key IDs for the Incorruptible TValues
const (
	keyUsername = iota
	KeyUserID
	keyNsName
	keyNsID
	keyAdminType
)

// adminLogout : http logout handler for the admin interface.
func adminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, Incorruptible.DeadCookie())
}

// adminLogin : http login handler for the admin interface.
func adminLogin(w http.ResponseWriter, r *http.Request) {
	var m passwordRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("AdminLogin DecodeJSONBody:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	username := m.Username
	password := m.Password
	namespace := m.Namespace

	if p := garcon.Printable(username, password, namespace); p >= 0 {
		log.ParamError("AdminLogin: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	// get the namespace
	exists, ns, err := db.SelectNsFromName(namespace)
	if err != nil {
		log.QueryError("AdminLogin SelectNsFromName:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "DB error SELECT namespace", "namespace", namespace)
		return
	}
	if !exists {
		log.ParamError("AdminLogin: namespace " + namespace + " does not exist")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist", "namespace", namespace)
		return
	}

	// check the user password
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		log.Error("AdminLogin checkUserPassword:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "internal error when checking user password", "usr", username, "namespace", namespace)
		return
	}
	if !isAuthorized {
		log.Info("AdminLogin u=" + username + " ns=" + namespace + ": disabled user or bad password")
		gw.WriteErr(w, r, http.StatusUnauthorized, "disabled user or bad password", "usr", username, "namespace", namespace)
		return
	}

	userType, err := db.GetUserType(ns.Name, ns.ID, u.ID)
	if err != nil {
		log.QueryError("AdminLogin AdminType:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "DB error when getting user type", "usr", username, "namespace", namespace)
		return
	}
	if userType == db.UserNoAdmin {
		log.ParamError("AdminLogin: u=" + username + " is not admin")
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin", "usr", username, "namespace", namespace)
		return
	}

	adminType := QuidAdmin
	if userType == db.NsAdmin {
		adminType = NsAdmin
	}
	log.Result("AdminLogin OK u=" + u.Name + " ns=" + namespace + " AdminType=" + string(adminType))

	// create a new Incorruptible cookie
	cookie, tv, err := Incorruptible.NewCookie(r,
		incorruptible.String(keyUsername, u.Name),
		incorruptible.Int64(KeyUserID, u.ID),
		incorruptible.String(keyNsName, ns.Name),
		incorruptible.Int64(keyNsID, ns.ID),
		incorruptible.String(keyAdminType, string(adminType)),
	)
	if err != nil {
		log.Error("AdminLogin NewCookie:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "internal error when creating a new incorruptible cookie", "usr", username, "namespace", namespace)
		return
	}

	http.SetCookie(w, cookie)

	sendStatusResponse(w, tv)
}

// status returns 200 if user is admin.
func status(w http.ResponseWriter, r *http.Request) {
	tv, err := Incorruptible.DecodeCookieToken(r)
	if err != nil {
		log.Warn("/status: no valid token:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "missing or invalid incorruptible cookie")
		return
	}

	sendStatusResponse(w, tv)
}

// status returns 200 if user is admin.
func sendStatusResponse(w http.ResponseWriter, tv incorruptible.TValues) {
	adminType := tv.StringIfAny(keyAdminType)

	gw.WriteOK(w, statusResponse{
		AdminType: AdminType(adminType),
		Username:  tv.StringIfAny(keyUsername),
		Ns: nsInfo{
			ID:   tv.Int64IfAny(keyNsID),
			Name: tv.StringIfAny(keyNsName),
		},
	})
}