package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllOrgs : get all orgs http handler.
func AllOrgs(c echo.Context) error {
	data, err := db.SelectAllOrgs()
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error selecting orgs",
		})
	}
	return c.JSON(http.StatusOK, &data)
}

// FindOrg : find an org from name.
func FindOrg(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m["name"].(string)

	data, err := db.SelectOrgStartsWith(name)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error finding org",
		})
	}
	return c.JSON(http.StatusOK, &data)
}

// UserOrgsInfo : get orgs info for a user.
func UserOrgsInfo(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	ID := int64(m["id"].(float64))

	o, err := db.SelectOrgsForUser(ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "error selecting orgs",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"orgs": o,
	})
}

// DeleteOrg : org deletion http handler.
func DeleteOrg(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}

	id := int64(m["id"].(float64))
	if err := db.DeleteOrg(id); err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error deleting org",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "ok",
	})
}

// CreateOrg : org creation http handler.
func CreateOrg(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m["name"].(string)

	org, exists, err := createOrg(name)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "error creating org",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "org already exists",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"org_id": org.ID,
	})
}

// createOrg : create an org.
func createOrg(name string) (server.Org, bool, error) {
	org := server.Org{}

	exists, err := db.OrgExists(name)
	if err != nil {
		return org, false, err
	}
	if exists {
		return org, true, nil
	}

	uid, err := db.CreateOrg(name)
	if err != nil {
		return org, false, err
	}
	org.ID = uid
	org.Name = name
	return org, false, nil
}
