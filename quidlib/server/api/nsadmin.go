package api

import (
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
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting admin users",
		})
	}
	return c.JSON(http.StatusOK, &data)
}

/* SearchForNonAdminUsersInNamespace : search from a username in namespace
func SearchForNonAdminUsersInNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	username := m["username"].(string)
	nsID := int64(m["namespace_id"].(float64))

	u, err := db.SearchUsersInNamespaceFromUsername(username, nsID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error searching for users",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"users": u,
	})
}*/

// CreateUserAdministrators : create admin users handler
func CreateAdministrators(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userIDs := m["user_ids"].([]interface{})
	namespaceID := int64(m["namespace_id"].(float64))
	for _, fuserID := range userIDs {
		userID := int64(fuserID.(float64))
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
		_, err = db.CreateAdministrator(namespaceID, userID)
		if err != nil {
			return c.JSON(http.StatusConflict, echo.Map{
				"error": "error creating admin user",
			})
		}
	}
	return c.JSON(http.StatusOK, okResponse())
}

// DeleteAdministrator : delete an admin user handler
func DeleteAdministrator(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userID := int64(m["user_id"].(float64))
	namespaceID := int64(m["namespace_id"].(float64))
	err := db.DeleteAdministrator(userID, namespaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error deleting admin users",
		})
	}
	return c.JSON(http.StatusOK, okResponse())
}
