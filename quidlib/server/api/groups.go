package api

import (
	"encoding/json"
	"net/http"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllGroupsForNamespace : get all groups for a namespace http handler.
func AllGroupsForNamespace(w http.ResponseWriter, r *http.Request) {
	var m NamespaceIDRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := db.SelectGroupsForNamespace(nsID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error selecting groups")
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

// AllGroups : get all groups for a namespace http handler.
func AllGroups(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllGroups()
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error selecting groups")
	}

	b, err := json.Marshal(&data)
	if err != nil {
		emo.Error("%v while serializing %v", err, data)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GroupsInfo : group creation http handler.
func GroupsInfo(w http.ResponseWriter, r *http.Request) {
	var m UserRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	n, err := db.CountUsersInGroup(id)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error counting in group")
		return
	}

	gw.WriteOK(w, "num_users", n)
}

// DeleteGroup : group deletion http handler.
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var m UserRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := db.DeleteGroup(id); err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error deleting group")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "message", "ok")
}

// CreateGroup : group creation http handler.
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var m GroupCreation
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name
	nsID := m.NamespaceID

	if p := garcon.Printable(name); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !VerifyAdminNs(w, r, nsID) {
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
