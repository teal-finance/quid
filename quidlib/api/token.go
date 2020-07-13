package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/synw/quid/quidlib/conf"
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

	nsid, err := db.SelectNamespaceID(namespace)
	if err != nil {
		return err
	}

	isAuthorized, u, err := checkUserPassword(username, password, nsid)
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

	t, err := tokens.GenUserToken(u.Name, u.GroupNames(), time.Now().Add(conf.DefaultTokenTimeout))
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"key": t,
	})
}
