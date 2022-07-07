package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllGroupsForNamespace : get all groups for a namespace http handler.
func AllGroupsForNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	nsID := int64(m["namespace_id"].(float64))

	data, err := db.SelectGroupsForNamespace(nsID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error selecting groups",
		})
	}

	return c.JSON(http.StatusOK, &data)
}

// AllGroups : get all groups for a namespace http handler.
func AllGroups(c echo.Context) error {
	data, err := db.SelectAllGroups()
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error selecting groups",
		})
	}

	return c.JSON(http.StatusOK, &data)
}

// GroupsInfo : group creation http handler.
func GroupsInfo(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	id := int64(m["id"].(float64))

	n, err := db.CountUsersInGroup(id)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error counting in group",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"num_users": n,
	})
}

// DeleteGroup : group deletion http handler.
func DeleteGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	id := int64(m["id"].(float64))

	if err := db.DeleteGroup(id); err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error deleting group",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

// CreateGroup : group creation http handler.
func CreateGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	name := m["name"].(string)
	nsID := int64(m["namespace_id"].(float64))

	ns, exists, err := createGroup(name, nsID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating group",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "group already exists",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"org_id": ns.ID,
	})
}

// createGroup : create a group.
func createGroup(name string, namespaceID int64) (server.Group, bool, error) {
	ns := server.Group{}

	exists, err := db.GroupExists(name, namespaceID)
	if err != nil {
		return ns, false, err
	}
	if exists {
		return ns, true, nil
	}

	uid, err := db.CreateGroup(name, namespaceID)
	if err != nil {
		return ns, false, err
	}

	ns.ID = uid
	ns.Name = name
	return ns, false, nil
}
