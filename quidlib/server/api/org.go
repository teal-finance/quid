package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllOrgs : get all orgs http handler.
func AllOrgs(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllOrgs()
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error selecting orgs")
		return
	}

	b, err := json.Marshal(&data)
	if err != nil {
		emo.Error("%v while serializing %v", err, data)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// FindOrg : find an org from name.
func FindOrg(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return
	}

	name := m["name"].(string)

	data, err := db.SelectOrgStartsWith(name)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding org")
		return
	}

	b, err := json.Marshal(&data)
	if err != nil {
		emo.Error("%v while serializing %v", err, data)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// UserOrgsInfo : get orgs info for a user.
func UserOrgsInfo(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return
	}

	id := int64(m["id"].(float64))

	o, err := db.SelectOrgsForUser(id)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting orgs")
		return
	}

	gw.WriteOK(w, "orgs", o)
}

// DeleteOrg : org deletion http handler.
func DeleteOrg(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return
	}

	id := int64(m["id"].(float64))

	if err := db.DeleteOrg(id); err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error deleting org")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "message", "ok")
}

// CreateOrg : org creation http handler.
func CreateOrg(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return
	}

	name := m["name"].(string)

	org, exists, err := createOrg(name)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error creating org")
		return
	}
	if exists {
		gw.WriteErr(w, r, http.StatusConflict, "org already exists")
		return
	}

	gw.WriteOK(w, "org_id", org.ID)
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

	id, err := db.CreateOrg(name)
	if err != nil {
		return org, false, err
	}

	org.ID = id
	org.Name = name
	return org, false, nil
}
