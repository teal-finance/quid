package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/tokens"
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

	// get the namespace id
	nsid, err := db.SelectNamespaceID(namespace)
	if err != nil {
		return err
	}
	// check the user password
	isAuthorized, u, err := checkUserPassword(username, password, nsid)
	if err != nil {
		return err
	}
	if isAuthorized == false {
		fmt.Println(username, "unauthorized")
		return c.NoContent(http.StatusUnauthorized)
	}
	// check the user admin group
	isAdmin, err := isUserInAdminGroup(u.ID, nsid)
	if err != nil {
		return err
	}
	if isAdmin == false {
		fmt.Println(username, "unauthorized: not in admin group")
		return c.NoContent(http.StatusUnauthorized)
	}

	// set the token
	token, err := tokens.GenAdminToken(u.Name)
	if err != nil {
		return err
	}

	fmt.Println("User", u.Name, "is connected")

	return c.JSON(http.StatusOK, echo.Map{
		"key": token,
	})
}
