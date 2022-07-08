package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllUsersInNamespace : select all users for a namespace.
func AllUsersInNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	data, err := db.SelectUsersInNamespace(nsID)
	if err != nil {
		fmt.Println("ERROR", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting users",
		})
	}

	return c.JSON(http.StatusOK, &data)
}

// GroupsForNamespace : get the groups of a user.
func GroupsForNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	namespace := m["namespace"].(string)

	hasResult, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil || !hasResult {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting namespace",
		})
	}

	g, err := db.SelectGroupsForNamespace(ns.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting groups",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"groups": g,
	})
}

// AddUserInOrg : add a user in an org.
func AddUserInOrg(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	uID := int64(m["user_id"].(float64))
	oID := int64(m["org_id"].(float64))

	err := db.AddUserInOrg(uID, oID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error adding user in org",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// RemoveUserFromOrg : add a user in an org.
func RemoveUserFromOrg(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	uID := int64(m["user_id"].(float64))
	oID := int64(m["org_id"].(float64))

	err := db.RemoveUserFromOrg(uID, oID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error removing user from org",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// AddUserInGroup : add a user in a group.
func AddUserInGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	uID := int64(m["user_id"].(float64))
	gID := int64(m["group_id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	err := db.AddUserInGroup(uID, gID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error adding user in group",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// RemoveUserFromGroup : add a user in a group.
func RemoveUserFromGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	uID := int64(m["user_id"].(float64))
	gID := int64(m["group_id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	err := db.RemoveUserFromGroup(uID, gID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error removing user from group",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

// SearchForUsersInNamespace : search from a username in namespace.
/*func SearchForUsersInNamespace(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	fmt.Println("Search")
	username := m["username"].(string)
	nsID := int64(m["namespace_id"].(float64))
	fmt.Println("U", username, "NS", nsID)

	u, err := db.SearchUsersInNamespaceFromUsername(username, nsID)
	if err != nil {
		fmt.Println("ERR", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error searching for users",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"users": u,
	})
}*/

// UserGroupsInfo : get info for a user.
func UserGroupsInfo(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	id := int64(m["id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	g, err := db.SelectGroupsForUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting groups",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"groups": g,
	})
}

// DeleteUser : delete a user handler.
func DeleteUser(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	id := int64(m["id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	if err := db.DeleteUser(id); err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error deleting user",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

// CreateUserHandler : create a user handler.
func CreateUserHandler(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	name := m["name"].(string)
	password := m["password"].(string)
	nsID := int64(m["namespace_id"].(float64))

	if !VerifyAdminNs(c, nsID) {
		return c.NoContent(http.StatusUnauthorized)
	}

	// check if user exists
	exists, err := db.UserNameExists(name, nsID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error checking user",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error user already exist",
		})
	}

	// create user
	u, err := db.CreateUser(name, password, nsID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating user",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user_id": u.ID,
	})
}

func checkUserPassword(username, password string, namespaceID int64) (bool, server.User, error) {
	found, u, err := db.SelectNonDisabledUser(username, namespaceID)
	if !found || err != nil {
		return false, u, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	if err != nil {
		fmt.Println("ERROR", "|"+err.Error()+"|")
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return false, u, err
		}
	}

	return true, u, err
}
