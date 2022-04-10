package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// AdminLogin : http login handler for the admin interface
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
	if err != nil {
		return err
	}
	if !isAuthorized {
		fmt.Println(username, "unauthorized")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	// check the user admin group
	isAdmin, err := isUserInAdminGroup(u.ID, ns.ID)
	if err != nil {
		return err
	}
	if !isAdmin {
		fmt.Println(username, "unauthorized: not in admin group")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	// set the session
	sameSite := http.SameSiteStrictMode
	if conf.IsDevMode {
		sameSite = http.SameSiteLaxMode
	}
	secure := !conf.IsDevMode
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   secure,
	}
	sess.Values["is_admin"] = "true"
	sess.Values["user"] = u.UserName
	//emo.Info("Setting session", u.Name, sess.Values["is_admin"])
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		emo.Error("Error saving session", err)
	}

	// set the refresh token
	token, err := tokens.GenRefreshToken("24h", ns.MaxRefreshTokenTTL, ns.Name, u.UserName, []byte(ns.RefreshKey))
	if err != nil {
		emo.Error("Error generating refresh token", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}
	if token == "" {
		emo.Error("Unauthorized: timeout max (", ns.MaxRefreshTokenTTL, ") for refresh token for namespace", ns.Name)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	fmt.Println("Admin user", u.UserName, "is connected")

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

// AdminLogout : http logout handler for the admin interface
func AdminLogout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values["is_admin"] = "false"
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
