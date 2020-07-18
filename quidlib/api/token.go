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

	exists, ns, err := db.SelectNamespace(namespace)
	if err != nil {
		return err
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "namespace does not exist",
		})
	}

	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		return err
	}

	// Respond with unauthorized status
	if isAuthorized == false {
		fmt.Println(username, "unauthorized")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthorized",
		})
	}

	fmt.Println("User", u.Name, "is connected")

	t, err := tokens.GenUserToken(u.Name, ns.Key, u.GroupNames(), ns.MaxTokenTTL)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"key": t,
	})
}
