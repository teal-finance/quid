package api

import (
	"net/http"

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

	isUserAdmin := false
	if ns.Name == "quid" {
		// check the user quid admin group
		isAdmin, err := db.IsUserInAdminGroup(u.ID, ns.ID)
		if err != nil {
			return err
		}
		if isAdmin {
			isUserAdmin = true
			emo.Info("Admin login successfull for user", u.Name, "on namespace", ns.Name)
		} else {
			emo.Warning(username, "unauthorized: user not in quid admin group")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "unauthorized",
			})
		}

	} else {
		// check if the user is namespace administrator
		exists, err := db.AdministratorExists(u.ID, ns.ID)
		if err != nil {
			return err
		}
		if !exists {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "unauthorized",
			})
		}
		emo.Info("Namespace administrator login successfull for user", u.Name, "on namespace", ns.Name)
	}

	// set the session
	sess, _ := session.Get("session", c)
	sess.Values["user"] = u.Name
	sess.Values["is_admin"] = isUserAdmin
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
		"token": token,
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
