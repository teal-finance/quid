package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible/tvalues"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

var gw garcon.Writer

// Key IDs for the Incorruptible TValues
const (
	user = iota
	user_id
	ns_name
	ns_id
	is_admin
	is_ns_admin
)

// AdminLogin : http login handler for the admin interface.
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var m passwordRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := m.Username
	password := m.Password
	namespace := m.Namespace

	if p := garcon.Printable(username, password, namespace); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist")
		return
	}

	// check the user password
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if !isAuthorized {
		emo.Warning(username, "unauthorized: password check failed", password, ns.ID)
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUserAdmin {
		emo.Warning(username, "unauthorized: user is not admin")
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	_isAdmin := isUserAdmin && namespace == "quid"
	_isNsAdmin := isUserAdmin && namespace != "quid"
	emo.Info("Admin login successful for user", u.Name, "on namespace", ns.Name)

	// get or create an Incorruptible token
	tv, ok := tvalues.FromCtx(r)
	if !ok {
		emo.Error("No Incorruptible token => Create a new one")
	}

	// update the token fields
	err = tv.Set(
		tv.KString(user, u.Name),
		tv.KInt64(user_id, u.ID),
		tv.KString(ns_name, ns.Name),
		tv.KInt64(ns_id, ns.ID),
		tv.KBool(is_admin, _isAdmin),
		tv.KBool(is_ns_admin, _isNsAdmin),
	)
	if err != nil {
		emo.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// store the Incorruptible token in the request context
	r = tv.ToCtx(r)

	// set the session
	cookie, err := Incorruptible.NewCookieFromValues(tv)
	if err != nil {
		emo.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

// AdminLogout : http logout handler for the admin interface.
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	tv, ok := tvalues.FromCtx(r)
	if !ok {
		emo.Error("No cookie or cookie is not Incorruptible")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := tv.SetBool(is_admin, false); err != nil {
		emo.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie, err := Incorruptible.NewCookieFromValues(tv)
	if err != nil {
		emo.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

// RequestAdminAccessToken : request an access token from a refresh token
// for a namespace.
func RequestAdminAccessToken(w http.ResponseWriter, r *http.Request) {
	var m adminAccessTokenRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := m.RefreshToken
	nsName := m.Namespace

	if p := garcon.Printable(refreshToken, nsName); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emo.RefreshToken(nsName, refreshToken)

	// get the namespace
	_, ns, err := db.SelectNamespaceFromName(nsName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emo.Error(err)
		log.Fatal(err)
		return
	}

	emo.State("Verifying refresh token")
	// verify the refresh token
	var username string
	token, err := jwt.ParseWithClaims(refreshToken, &tokens.RefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ns.RefreshKey), nil
	})

	if claims, ok := token.Claims.(*tokens.RefreshClaims); ok && token.Valid {
		username = claims.UserName
		fmt.Printf("%v %v", claims.UserName, claims.ExpiresAt)
	} else {
		emo.Warning(err.Error())
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	// get the user
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if !found {
		emo.Warning("User not found: " + username)
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if the user is in the admin group
	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUserAdmin {
		emo.Warning("Admin access token request from user", u.Name, "that is not admin for namespace", ns.Name)
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
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
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// emo.AccessToken("Issued an admin access token for user", u.Name, "and namespace", ns.Name)
	gw.WriteOK(w, "token", t, "namespace", ns)
}
