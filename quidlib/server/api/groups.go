package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllGroupsForNamespace : get all groups for a namespace http handler.
func AllGroupsForNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		logg.Warning("AllGroupsForNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := db.SelectGroupsForNamespace(nsID)
	if err != nil {
		logg.Warning("AllGroupsForNamespace: error selecting groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error selecting groups")
		return
	}

	gw.WriteOK(w, data)
}

// AllGroups : get all groups for a namespace http handler.
func AllGroups(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllGroups()
	if err != nil {
		logg.Warning("AllGroups: error selecting groups:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error selecting groups")
	}

	gw.WriteOK(w, data)
}

// GroupsInfo : group creation http handler.
func GroupsInfo(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		logg.Warning("GroupsInfo:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	n, err := db.CountUsersInGroup(id)
	if err != nil {
		logg.Warning("GroupsInfo: error counting in group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error counting in group")
		return
	}

	gw.WriteOK(w, "num_users", n)
}

// DeleteGroup : group deletion http handler.
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		logg.Warning("DeleteGroup:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !IsNsAdmin(r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := db.DeleteGroup(id); err != nil {
		logg.Warning("DeleteGroup: error deleting group:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting group")
		return
	}

	gw.WriteOK(w, "message", "ok")
}

// CreateGroup : group creation http handler.
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var m groupCreation
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		logg.Warning("CreateGroup:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name
	nsID := m.NamespaceID

	if p := garcon.Printable(name); p >= 0 {
		logg.Warning("CreateGroup: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsNsAdmin(r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
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
		logg.QueryError("createGroup GroupExists:", err)
		return ns, false, err
	}
	if exists {
		logg.QueryError("createGroup: group already exists")
		return ns, true, nil
	}

	uid, err := db.CreateGroup(name, namespaceID)
	if err != nil {
		logg.QueryError("createGroup:", err)
		return ns, false, err
	}

	ns.ID = uid
	ns.Name = name
	return ns, false, nil
}
