package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllNsGroups : get all groups for a namespace http handler.
func AllNsGroups(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("AllGroupsForNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	data, err := db.SelectNsGroups(nsID)
	if err != nil {
		log.QueryError("AllGroupsForNamespace: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error SELECT groups")
		return
	}

	gw.WriteOK(w, data)
}

// AllGroups : get all groups for a namespace http handler.
// Deprecated because this function is not used any longer.
func AllGroups(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllGroups()
	if err != nil {
		log.QueryError("AllGroups: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error SELECT groups")
	}

	gw.WriteOK(w, data)
}

// GroupsInfo : group creation http handler.
func GroupsInfo(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("GroupsInfo:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	n, err := db.CountUsersInGroup(id)
	if err != nil {
		log.QueryError("GroupsInfo: error counting in group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error counting in group")
		return
	}

	gw.WriteOK(w, "num_users", n)
}

// DeleteGroup : group deletion http handler.
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("DeleteGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	if err := db.DeleteGroup(id); err != nil {
		log.QueryError("DeleteGroup: error deleting group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting group")
		return
	}

	gw.WriteOK(w, "message", "ok")
}

// CreateGroup : group creation http handler.
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var m groupCreation
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name
	nsID := m.NamespaceID

	if p := garcon.Printable(name); p >= 0 {
		log.Warn("CreateGroup: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	if !IsNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	ns, exists, err := createGroup(name, nsID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error creating group")
		return
	}
	if exists {
		gw.WriteErr(w, r, http.StatusConflict, "group already exists")
		return
	}

	gw.WriteOK(w, "org_id", ns.ID)
}

// createGroup : create a group.
func createGroup(name string, namespaceID int64) (server.Group, bool, error) {
	var ns server.Group

	exists, err := db.GroupExists(name, namespaceID)
	if err != nil {
		log.QueryError("createGroup GroupExists:", err)
		return ns, false, err
	}
	if exists {
		log.QueryError("createGroup: group already exists")
		return ns, true, nil
	}

	uid, err := db.CreateGroup(name, namespaceID)
	if err != nil {
		log.QueryError("createGroup:", err)
		return ns, false, err
	}

	ns.ID = uid
	ns.Name = name
	return ns, false, nil
}
