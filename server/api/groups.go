package api

import (
	"net/http"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	db "github.com/teal-finance/quid/server/db"
)

// allNsGroups : get all groups for a namespace http handler.
func allNsGroups(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceIDRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("AllGroupsForNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if !isNsAdmin(r, m.NamespaceID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", m.NamespaceID)
		return
	}

	data, err := db.SelectNsGroups(m.NamespaceID)
	if err != nil {
		log.QueryError("AllGroupsForNamespace: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error SELECT groups")
		return
	}

	gw.WriteOK(w, data)
}

// allGroups : get all groups for a namespace http handler.
// Deprecated because this function is not used any longer.
func AllGroups(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllGroups()
	if err != nil {
		log.QueryError("AllGroups: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error SELECT groups")
	}

	gw.WriteOK(w, data)
}

// groupsInfo : group creation http handler.
func groupsInfo(w http.ResponseWriter, r *http.Request) {
	var m server.UserRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("GroupsInfo:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if !isNsAdmin(r, m.NamespaceID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", m.NamespaceID)
		return
	}

	n, err := db.CountUsersInGroup(m.ID)
	if err != nil {
		log.QueryError("GroupsInfo: error counting in group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error counting in group")
		return
	}

	gw.WriteOK(w, "num_users", n)
}

// deleteGroup : group deletion http handler.
func deleteGroup(w http.ResponseWriter, r *http.Request) {
	var m server.UserRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("DeleteGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if !isNsAdmin(r, m.NamespaceID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", m.NamespaceID)
		return
	}

	if err := db.DeleteGroup(m.ID); err != nil {
		log.QueryError("DeleteGroup: error deleting group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting group")
		return
	}

	gw.WriteOK(w, "message", "ok")
}

// createGroup : group creation http handler.
func createGroup(w http.ResponseWriter, r *http.Request) {
	var m server.GroupCreation
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if p := gg.Printable(m.Name); p >= 0 {
		log.Warn("CreateGroup: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	if !isNsAdmin(r, m.NamespaceID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", m.NamespaceID)
		return
	}

	grp, exists, err := db.CreateGroupIfExist(m.Name, m.NamespaceID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error creating group")
		return
	}
	if exists {
		gw.WriteErr(w, r, http.StatusConflict, "group already exists")
		return
	}

	gw.WriteOK(w, "grp_id", grp.ID)
}
