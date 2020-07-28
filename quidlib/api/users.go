package api

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/models"
)

// GroupsForNamespace : get the groups of a user
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

// AddUserInGroup : add a user in a group
func AddUserInGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userID := int64(m["user_id"].(float64))
	groupID := int64(m["group_id"].(float64))

	err := db.AddUserInGroup(userID, groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error adding user in group",
		})
	}
	return c.NoContent(http.StatusOK)

}

// RemoveUserFromGroup : add a user in a group
func RemoveUserFromGroup(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userID := int64(m["user_id"].(float64))
	groupID := int64(m["group_id"].(float64))

	err := db.RemoveUserFromGroup(userID, groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error removing user from group",
		})
	}
	return c.NoContent(http.StatusOK)

}

// UserInfo : get info for a user
func UserInfo(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	g, err := db.SelectGroupsForUser(ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting groups",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"groups": g,
	})

}

// DeleteUser : delete a user handler
func DeleteUser(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	err := db.DeleteUser(ID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error deleting user",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})

}

// CreateUserHandler : create a user handler
func CreateUserHandler(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m["name"].(string)
	password := m["password"].(string)
	namespaceID := int64(m["namespace_id"].(float64))

	// check if user exists
	exists, err := db.UserNameExists(name)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error checking user",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating user",
		})
	}

	// create user
	u, err := CreateUser(name, password, namespaceID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating user",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user_id": u.ID,
	})
}

// CreateUser : create a user
func CreateUser(username string, password string, namespaceID int64) (models.User, error) {
	user := models.User{}
	pwd := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	uid, err := db.CreateUserFromNameAndPassword(username, string(hashedPassword), namespaceID)
	if err != nil {
		return user, err
	}
	user.ID = uid
	user.Name = username
	return user, nil
}

func checkUserPassword(username string, password string, namespaceID int64) (bool, models.User, error) {
	found, u, err := db.SelectUser(username, namespaceID)
	if err != nil {
		return false, u, err
	}
	if !found {
		return false, u, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return false, u, nil
	}
	return true, u, nil
}

func isUserInAdminGroup(uid int64, nsid int64) (bool, error) {
	gid, err := db.SelectGroup("quid_admin", nsid)
	if err != nil {
		return false, err
	}
	isAdmin, err := db.IsUserInGroup(uid, gid.ID, nsid)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}
