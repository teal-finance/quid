package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// AllAdministratorsInNamespace : select all admin users for a namespace
func AllAdministratorsInNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	namespaceID := int64(m["namespace_id"].(float64))
	data, err := db.SelectAdministratorsInNamespace(namespaceID)
	if err != nil {
		fmt.Println("ERROR", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting admin users",
		})
	}
	return c.JSON(http.StatusOK, &data)
}

// CreateUserAdministrator : create an admin user handler
func CreateAdministrator(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userID := int64(m["user_id"].(float64))
	namespaceID := int64(m["namespace_id"].(float64))

	// check if user exists
	exists, err := db.AdministratorExists(namespaceID, userID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error checking admin user",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating admin user",
		})
	}

	// create admin user
	id, err := db.CreateAdministrator(namespaceID, userID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating admin user",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"administrator_id": id,
	})
}
