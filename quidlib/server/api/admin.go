package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

var gw garcon.Writer

// AdminLogin : http login handler for the admin interface.
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var m PasswordRequest
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
	// set the session
	sess, _ := session.Get("session", c)
	sess.Values["user"] = u.Name
	sess.Values["ns_name"] = ns.Name
	sess.Values["ns_id"] = ns.ID
	sess.Values["is_admin"] = _isAdmin
	sess.Values["is_ns_admin"] = _isNsAdmin
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   !conf.IsDevMode,
		SameSite: http.SameSiteStrictMode,
	}
	if conf.IsDevMode {
		sess.Options.SameSite = http.SameSiteLaxMode
	}

	if err = sess.Save(r, w); err != nil {
		emo.Error("Error saving session", err)
	}

	// set the refresh token
	token, err := tokens.GenRefreshToken("24h", ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		emo.Error("Error generating refresh token", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "cannot generate the refresh token")
		return
	}
	if token == "" {
		emo.Warning("Unauthorized: timeout max (", ns.MaxRefreshTokenTTL, ") for refresh token for namespace", ns.Name)
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	emo.Info("Admin user ", u.Name, " is connected")

	gw.WriteErr(w, r, http.StatusOK, "token", token, "namespace", ns)
}

// AdminLogout : http logout handler for the admin interface.
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Get("session", c)
	sess.Values["is_admin"] = "false"

	if err := sess.Save(r, w); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RequestAdminAccessToken : request an access token from a refresh token
// for a namespace.
func RequestAdminAccessToken(w http.ResponseWriter, r *http.Request) {
	var m AdminAccessTokenRequest
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

	//emo.AccessToken("Issued an admin access token for user", u.Name, "and namespace", ns.Name)
	gw.WriteOK(w, "token", t, "namespace", ns)
}
