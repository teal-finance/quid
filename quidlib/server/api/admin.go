package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// Key IDs for the Incorruptible TValues
const (
	keyUsername = iota
	KeyUserID
	keyNsName
	keyNsID
	keyAdminType
)

// AdminLogin : http login handler for the admin interface.
func AdminLogin(w http.ResponseWriter, r *http.Request) {
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

// AdminLogout : http logout handler for the admin interface.
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, Incorruptible.DeadCookie())
}

// RequestAdminAccessToken : request an access token from a refresh token
// for a namespace.
/*func RequestAdminAccessToken(w http.ResponseWriter, r *http.Request) {
	var m adminAccessTokenRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RequestAdminAccessToken DecodeJSONBody:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	refreshToken := m.RefreshToken
	nsName := m.Namespace

	if p := garcon.Printable(refreshToken, nsName); p >= 0 {
		log.ParamError("RequestAdminAccessToken: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.RefreshToken(nsName, refreshToken)

	// get the namespace
	_, ns, err := db.SelectNsFromName(nsName)
	if err != nil {
		log.QueryError("RequestAdminAccessToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.State("Verifying refresh token")
	// verify the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &tokens.RefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ns.RefreshKey), nil
	})
	if err != nil {
		log.QueryError("RequestAdminAccessToken ParseWithClaims:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !token.Valid {
		log.Warn("RequestAdminAccessToken: invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*tokens.RefreshClaims)
	if !ok {
		log.Warn("RequestAdminAccessToken: cannot convert to RefreshClaims")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Result("RequestAdminAccessToken: u="+claims.UserName+" expires=", claims.ExpiresAt)

	// get the user
	username := claims.UserName
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if err != nil {
		log.QueryError("RequestAdminAccessToken SelectNonDisabledUser:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		log.Data("RequestAdminAccessToken: user not found: " + username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user is in the admin group
	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		log.QueryError("RequestAdminAccessToken IsUserAdmin:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUserAdmin {
		log.Data("RequestAdminAccessToken: Admin access token request from u="+u.Name+" is not admin for ns=", ns.Name)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_isAdmin := false
	_isNsAdmin := false
	if ns.Name != "quid" {
		_isAdmin = true
	} else {
		_isNsAdmin = true
	}

	// generate the access token
	t, err := tokens.GenAdminAccessToken(ns.Name, "5m", ns.MaxTokenTTL, u.Name, u.ID, ns.ID, []byte(AdminNsKey), _isAdmin, _isNsAdmin)
	if err != nil {
		log.Error("RequestAdminAccessToken GenAdminAccessToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.AccessToken("Issued an admin access token for user", u.Name, "and namespace", ns.Name)
	gw.WriteOK(w, "token", t, "namespace", ns)
}*/
