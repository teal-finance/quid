package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// AdminLogin : http login handler for the admin interface.
func AdminLogin(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	username := m["username"].(string)
	password := m["password"].(string)
	namespace := m["namespace"].(string)

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		return err
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "namespace does not exist",
		})
	}

	// check the user password
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if !isAuthorized {
		emo.Warning(username, "unauthorized: password check failed", password, ns.ID)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	if err != nil {
		return err
	}

	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		return err
	}
	if !isUserAdmin {
		emo.Warning(username, "unauthorized: user is not admin")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
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

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		emo.Error("Error saving session", err)
	}

	// set the refresh token
	token, err := tokens.GenRefreshToken("24h", ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		emo.Error("Error generating refresh token", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}
	if token == "" {
		emo.Warning("Unauthorized: timeout max (", ns.MaxRefreshTokenTTL, ") for refresh token for namespace", ns.Name)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	emo.Info("Admin user ", u.Name, " is connected")

	return c.JSON(http.StatusOK, echo.Map{
		"token":     token,
		"namespace": ns,
	})
}

// AdminLogout : http logout handler for the admin interface.
func AdminLogout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values["is_admin"] = "false"

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// RequestAdminAccessToken : request an access token from a refresh token
// for a namespace.
func RequestAdminAccessToken(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	refreshToken, ok := m["refresh_token"].(string)
	nsName := m["namespace"].(string)
	emo.RefreshToken(nsName, refreshToken)
	if !ok {
		emo.ParamError("provide a refresh_token parameter")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "provide a refresh_token parameter",
		})
	}

	// get the namespace
	_, ns, err := db.SelectNamespaceFromName(nsName)
	if err != nil {
		emo.Error(err)
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
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
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	// get the user
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if !found {
		emo.Warning("User not found: " + username)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	// check if the user is in the admin group
	isUserAdmin, err := db.IsUserAdmin(ns.Name, ns.ID, u.ID)
	if err != nil {
		return err
	}
	if !isUserAdmin {
		emo.Warning("Admin access token request from user", u.Name, "that is not admin for namespace", ns.Name)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
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
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	//emo.AccessToken("Issued an admin access token for user", u.Name, "and namespace", ns.Name)

	return c.JSON(http.StatusOK, echo.Map{
		"token":     t,
		"namespace": ns,
	})
}
