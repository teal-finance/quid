package api

import (
	"net/http"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	db "github.com/teal-finance/quid/server/db"
)

// allOrgs : get all orgs http handler.
func allOrgs(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllOrgs()
	if err != nil {
		log.QueryError("AllOrgs: error SELECT orgs:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error SELECT orgs")
		return
	}

	gw.WriteOK(w, data)
}

// findOrg : find an org from name.
func findOrg(w http.ResponseWriter, r *http.Request) {
	var m server.NameRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("FindOrg:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name

	if p := gg.Printable(name); p >= 0 {
		log.Warn("FindOrg: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	data, err := db.SelectOrgStartsWith(name)
	if err != nil {
		log.QueryError("FindOrg:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding org")
		return
	}

	gw.WriteOK(w, data)
}

// userOrgsInfo : get orgs info for a user.
func userOrgsInfo(w http.ResponseWriter, r *http.Request) {
	var m server.InfoRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("UserOrgsInfo:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID

	o, err := db.SelectOrgsForUser(id)
	if err != nil {
		log.QueryError("UserOrgsInfo: error SELECT orgs:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT orgs")
		return
	}

	gw.WriteOK(w, "orgs", o)
}

// deleteOrg : org deletion http handler.
func deleteOrg(w http.ResponseWriter, r *http.Request) {
	var m server.InfoRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("DeleteOrg:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID

	if err := db.DeleteOrg(id); err != nil {
		log.QueryError("DeleteOrg:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting org")
		return
	}

	gw.WriteOK(w, "message", "ok")
}

// createOrg : org creation http handler.
func createOrg(w http.ResponseWriter, r *http.Request) {
	var m server.NameRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("CreateOrg:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name

	if p := gg.Printable(name); p >= 0 {
		log.ParamError("CreateOrg: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	org, exists, err := db.CreateOrgIfExist(name)
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
