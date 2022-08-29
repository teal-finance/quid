package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllOrgs : get all orgs http handler.
func AllOrgs(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllOrgs()
	if err != nil {
		emo.QueryError("AllOrgs: error selecting orgs:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error selecting orgs")
		return
	}

	gw.WriteOK(w, data)
}

// FindOrg : find an org from name.
func FindOrg(w http.ResponseWriter, r *http.Request) {
	var m nameRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("FindOrg:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name

	if p := garcon.Printable(name); p >= 0 {
		emo.Warning("FindOrg: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := db.SelectOrgStartsWith(name)
	if err != nil {
		emo.QueryError("FindOrg:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding org")
		return
	}

	gw.WriteOK(w, data)
}

// UserOrgsInfo : get orgs info for a user.
func UserOrgsInfo(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("UserOrgsInfo:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

	o, err := db.SelectOrgsForUser(id)
	if err != nil {
		emo.QueryError("UserOrgsInfo: error selecting orgs:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting orgs")
		return
	}

	gw.WriteOK(w, "orgs", o)
}

// DeleteOrg : org deletion http handler.
func DeleteOrg(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("DeleteOrg:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

	if err := db.DeleteOrg(id); err != nil {
		emo.QueryError("DeleteOrg:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting org")
		return
	}

	gw.WriteOK(w, "message", "ok")
}

// CreateOrg : org creation http handler.
func CreateOrg(w http.ResponseWriter, r *http.Request) {
	var m nameRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		emo.ParamError("CreateOrg:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name

	if p := garcon.Printable(name); p >= 0 {
		emo.ParamError("CreateOrg: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	var org server.Org

	exists, err := db.OrgExists(name)
	if err != nil {
		emo.QueryError("createOrg OrgExists:", err)
		return org, false, err
	}
	if exists {
		emo.QueryError("createOrg: already exist:", name)
		return org, true, nil
	}

	id, err := db.CreateOrg(name)
	if err != nil {
		emo.QueryError("createOrg:", err)
		return org, false, err
	}

	org.ID = id
	org.Name = name
	return org, false, nil
}
