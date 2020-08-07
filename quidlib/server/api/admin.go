package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/synw/quid/quidlib/conf"
	"github.com/synw/quid/quidlib/server/db"
	"github.com/synw/quid/quidlib/tokens"
)

// AdminLogout : http logout handler for the admin interface
func AdminLogout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values["is_admin"] = "false"
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

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
	if isAuthorized == false {
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
	if isAdmin == false {
		fmt.Println(username, "unauthorized: not in admin group")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	// set the session
	sameSite := http.SameSiteStrictMode
	if conf.IsDevMode {
		sameSite = http.SameSiteNoneMode
	}
	secure := !conf.IsDevMode
	if !conf.IsDevMode {
		// sessions are not used in dev mode
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600 * 24,
			HttpOnly: true,
			SameSite: sameSite,
			Secure:   secure,
		}
		sess.Values["is_admin"] = "true"
		sess.Values["user"] = u.Name
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			emo.Error("Error saving session", err)
		}
	}

	// set the refresh token
	exists, token, err := tokens.GenRefreshToken(ns.Name, ns.RefreshKey, ns.MaxRefreshTokenTTL, u.Name, "24h")
	if !exists {
		emo.Error("Unauthorized: timeout max (", ns.MaxRefreshTokenTTL, ") for refresh token for namespace", ns.Name)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	if err != nil {
		emo.Error("Error generating refresh token", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	fmt.Println("Admin user", u.Name, "is connected")

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
