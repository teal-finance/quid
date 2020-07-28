package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/tokens"
)

// RequestToken : http login handler
func RequestToken(c echo.Context) error {

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	username := m["username"].(string)
	password := m["password"].(string)
	namespace := m["namespace"].(string)
	timeout := c.Param("timeout")

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

	// check if the user password matches
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		return err
	}
	// Respond with unauthorized status
	if isAuthorized == false {
		fmt.Println(username, "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}

	// get the user group names
	groupNames, err := db.SelectGroupsNamesForUser(u.ID)
	if err != nil {
		return err
	}

	// generate the token
	isAuth, t, err := tokens.GenUserToken(ns, u.Name, groupNames, timeout, ns.MaxTokenTTL)
	if err != nil {
		log.Fatal(err)
	}
	if !isAuth {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "max timeout exceeded",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
