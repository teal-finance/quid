package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllGroupsForNamespace : get all groups for a namespace http handler.
func AllGroupsForNamespace(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return
	}

	nsID := int64(m["namespace_id"].(float64))

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
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

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
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

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
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	name := m["name"].(string)
	nsID := int64(m["namespace_id"].(float64))

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
