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
	keyIsAdmin
	keyIsNsAdmin
)

// AdminLogin : http login handler for the admin interface.
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var m passwordRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("AdminLogin DecodeJSONBody:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := m.Username
	password := m.Password
	namespace := m.Namespace

	if p := garcon.Printable(username, password, namespace); p >= 0 {
		emo.ParamError("AdminLogin: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		emo.QueryError("AdminLogin SelectNamespaceFromName:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		emo.ParamError("AdminLogin: namespace " + namespace + " does not exist")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist", "namespace", namespace)
		return
	}

	// check the user password
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		emo.Error("AdminLogin checkUserPassword:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAuthorized {
		emo.ParamError("AdminLogin: Bad password for " + username + " ns=" + namespace)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		emo.QueryError("AdminLogin IsUserAdmin:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUserAdmin {
		emo.ParamError("AdminLogin: u=" + username + " is not admin")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_isAdmin := isUserAdmin && namespace == "quid"
	_isNsAdmin := isUserAdmin && namespace != "quid"
	emo.Result("AdminLogin OK u=" + u.Name + " ns=" + namespace)

	// get or create an Incorruptible token
	tv, ok := incorruptible.FromCtx(r)
	if !ok {
		emo.Info("AdminLogin: no Incorruptible token => Create a new one u="+u.Name+" (id=", u.ID,
			") ns="+ns.Name+" (id=", ns.ID, ") admin=", _isAdmin, "NSAdmin=", _isNsAdmin)
	}

	// update the token fields
	err = tv.Set(
		tv.KString(keyUsername, u.Name),
		tv.KInt64(KeyUserID, u.ID),
		tv.KString(keyNsName, ns.Name),
		tv.KInt64(keyNsID, ns.ID),
		tv.KBool(keyIsAdmin, _isAdmin),
		tv.KBool(keyIsNsAdmin, _isNsAdmin),
	)
	if err != nil {
		emo.Error("AdminLogin tv.Set:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// store the Incorruptible token in the request context
	_ = tv.ToCtx(r)

	// set the session
	cookie, err := Incorruptible.NewCookieFromValues(tv)
	if err != nil {
		emo.Error("AdminLogin NewCookieFromValues:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	gw.WriteOK(w, statusResponse{
		User: unInfo{
			Admin: tv.BoolIfAny(keyIsAdmin),
			ID:    tv.Int64IfAny(KeyUserID),
			Name:  tv.StringIfAny(keyUsername),
		},
		Ns: unInfo{
			Admin: tv.BoolIfAny(keyIsNsAdmin),
			ID:    tv.Int64IfAny(keyNsID),
			Name:  tv.StringIfAny(keyNsName),
		},
	})
}

// AdminLogout : http logout handler for the admin interface.
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	tv, ok := incorruptible.FromCtx(r)
	if !ok {
		emo.Warning("AdminLogout: missing valid Incorruptible cookie")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emo.Result("AdminLogout OK")

	if err := tv.SetBool(keyIsAdmin, false); err != nil {
		emo.Error("AdminLogout tv.SetBool:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := Incorruptible.NewCookieFromValues(tv)
	if err != nil {
		emo.Error("AdminLogout NewCookieFromValues:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// RequestAdminAccessToken : request an access token from a refresh token
// for a namespace.
/*func RequestAdminAccessToken(w http.ResponseWriter, r *http.Request) {
	var m adminAccessTokenRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("RequestAdminAccessToken DecodeJSONBody:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := m.RefreshToken
	nsName := m.Namespace

	if p := garcon.Printable(refreshToken, nsName); p >= 0 {
		emo.ParamError("RequestAdminAccessToken: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emo.RefreshToken(nsName, refreshToken)

	// get the namespace
	_, ns, err := db.SelectNamespaceFromName(nsName)
	if err != nil {
		emo.QueryError("RequestAdminAccessToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emo.State("Verifying refresh token")
	// verify the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &tokens.RefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ns.RefreshKey), nil
	})
	if err != nil {
		emo.QueryError("RequestAdminAccessToken ParseWithClaims:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !token.Valid {
		emo.Warning("RequestAdminAccessToken: invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*tokens.RefreshClaims)
	if !ok {
		emo.Warning("RequestAdminAccessToken: cannot convert to RefreshClaims")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	emo.Result("RequestAdminAccessToken: u="+claims.UserName+" expires=", claims.ExpiresAt)

	// get the user
	username := claims.UserName
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if err != nil {
		emo.QueryError("RequestAdminAccessToken SelectNonDisabledUser:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		emo.Data("RequestAdminAccessToken: user not found: " + username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user is in the admin group
	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		emo.QueryError("RequestAdminAccessToken IsUserAdmin:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUserAdmin {
		emo.Data("RequestAdminAccessToken: Admin access token request from u="+u.Name+" is not admin for ns=", ns.Name)
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
		emo.Error("RequestAdminAccessToken GenAdminAccessToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emo.AccessToken("Issued an admin access token for user", u.Name, "and namespace", ns.Name)
	gw.WriteOK(w, "token", t, "namespace", ns)
}*/
