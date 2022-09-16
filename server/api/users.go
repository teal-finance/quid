package api

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	db "github.com/teal-finance/quid/server/db"
)

// allNsUsers : select all users for a namespace.
func allNsUsers(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceIDRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("AllUsersInNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	nsID := m.NamespaceID

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	data, err := db.SelectNsUsers(nsID)
	if err != nil {
		log.QueryError("AllUsersInNamespace: error SELECT users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT users")
		return
	}

	gw.WriteOK(w, data)
}

// nsGroups : get the groups of a user.
func nsGroups(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("GroupsForNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	namespace := m.Namespace

	if p := gg.Printable(namespace); p >= 0 {
		log.ParamError("GroupsForNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	hasResult, ns, err := db.SelectNsFromName(namespace)
	if err != nil || !hasResult {
		log.QueryError("GroupsForNamespace: error SELECT namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT namespace")
		return
	}

	g, err := db.SelectNsGroups(ns.ID)
	if err != nil {
		log.QueryError("GroupsForNamespace: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT groups")
	}

	log.Result("GroupsForNamespace:", g)
	gw.WriteOK(w, "groups", g)
}

// addUserInOrg : add a user in an org.
func addUserInOrg(w http.ResponseWriter, r *http.Request) {
	var m server.UserOrgRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("AddUserInOrg:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	uID := m.UserID
	oID := m.OrgID

	err := db.AddUserInOrg(uID, oID)
	if err != nil {
		log.QueryError("AddUserInOrg: error adding user in org:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error adding user in org")
		return
	}

	log.Result("AddUserInOrg OK")
	w.WriteHeader(http.StatusOK)
}

// removeUserFromOrg : add a user in an org.
func removeUserFromOrg(w http.ResponseWriter, r *http.Request) {
	var m server.UserOrgRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RemoveUserFromOrg:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	uID := m.UserID
	oID := m.OrgID

	err := db.RemoveUserFromOrg(uID, oID)
	if err != nil {
		log.QueryError("RemoveUserFromOrg: error removing user from org:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error removing user from org")
		return
	}

	log.Result("RemoveUserFromOrg OK")
	w.WriteHeader(http.StatusOK)
}

// addUserInGroup : add a user in a group.
func addUserInGroup(w http.ResponseWriter, r *http.Request) {
	var m server.UserGroupRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("AddUserInGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	uID := m.UserID
	gID := m.GroupID
	nsID := m.NamespaceID

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	err := db.AddUserInGroup(uID, gID)
	if err != nil {
		log.QueryError("AddUserInGroup:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error adding user in group")
		return
	}

	log.Result("AddUserInGroup OK")
	w.WriteHeader(http.StatusOK)
}

// removeUserFromGroup : add a user in a group.
func removeUserFromGroup(w http.ResponseWriter, r *http.Request) {
	var m server.UserGroupRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RemoveUserFromGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	uID := m.UserID
	gID := m.GroupID
	nsID := m.NamespaceID

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	err := db.RemoveUserFromGroup(uID, gID)
	if err != nil {
		log.QueryError("RemoveUserFromGroup: error removing user from group:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error removing user from group")
		return
	}

	log.Result("RemoveUserFromGroup OK")
	w.WriteHeader(http.StatusOK)
}

// userGroupsInfo : get info for a user.
func userGroupsInfo(w http.ResponseWriter, r *http.Request) {
	var m server.UserRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("UserGroupsInfo:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	g, err := db.SelectGroupsForUser(id)
	if err != nil {
		log.QueryError("UserGroupsInfo: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT groups")
		return
	}

	log.Result("UserGroupsInfo:", g)
	gw.WriteOK(w, "groups", g)
}

// deleteUser : delete a user handler.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var m server.UserRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("DeleteUser:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	if err := db.DeleteUser(id); err != nil {
		log.QueryError("DeleteUser: error deleting user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting user")
		return
	}

	log.Result("DeleteUser OK")
	gw.WriteOK(w, "message", "ok")
}

// createUser : create a user handler.
func createUser(w http.ResponseWriter, r *http.Request) {
	var m server.UserHandlerCreation
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("CreateUser:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name
	password := m.Password
	nsID := m.NamespaceID

	if p := gg.Printable(name, password); p >= 0 {
		log.ParamError("CreateUser: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character")
		return
	}

	if !isNsAdmin(r, nsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "namespace_id", nsID)
		return
	}

	// check if user exists
	exists, err := db.UserExists(name, nsID)
	if err != nil {
		log.QueryError("CreateUser: error checking user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error checking user")
		return
	}
	if exists {
		log.Data("CreateUser: error user already exist")
		gw.WriteErr(w, r, http.StatusConflict, "error user already exist")
		return
	}

	// create user
	u, err := db.CreateUser(name, password, nsID)
	if err != nil {
		log.QueryError("CreateUser: error creating user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error creating user")
		return
	}

	log.Result("CreateUser:", u)
	gw.WriteOK(w, "user_id", u.ID)
}

func checkUserPassword(username, password string, namespaceID int64) (bool, server.User, error) {
	found, u, err := db.SelectEnabledUser(username, namespaceID)
	if !found || err != nil {
		return false, u, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Error(err)
		// "crypto/bcrypt: hashedPassword is not the hash of the given password"
		return false, u, err
	}

	return true, u, err
}
