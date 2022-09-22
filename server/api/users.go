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

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	data, err := db.SelectNsUsers(m.NsID)
	if err != nil {
		log.QueryError("AllUsersInNamespace: error SELECT users:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error SELECT users")
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

	if p := gg.Printable(m.Namespace); p >= 0 {
		log.ParamError("GroupsForNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	// get the namespace
	ns, err := db.SelectNsFromName(m.Namespace)
	if err != nil {
		log.QueryError(err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot SELECT namespace", "namespace", m.Namespace, "error", err)
		return
	}

	g, err := db.SelectNsGroups(ns.ID)
	if err != nil {
		log.QueryError("GroupsForNamespace: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error SELECT groups")
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

	err := db.AddUserInOrg(m.UsrID, m.OrgID)
	if err != nil {
		log.QueryError("AddUserInOrg: error adding user in org:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error adding user in org")
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

	err := db.RemoveUserFromOrg(m.UsrID, m.OrgID)
	if err != nil {
		log.QueryError("RemoveUserFromOrg: error removing user from org:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error removing user from org")
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

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	err := db.AddUserInGroup(m.UsrID, m.GrpID)
	if err != nil {
		log.QueryError("AddUserInGroup:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error adding user in group")
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

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	err := db.RemoveUserFromGroup(m.UsrID, m.GrpID)
	if err != nil {
		log.QueryError("RemoveUserFromGroup: error removing user from group:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error removing user from group")
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

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	g, err := db.SelectGroupsForUser(m.ID)
	if err != nil {
		log.QueryError("UserGroupsInfo: error SELECT groups:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error SELECT groups")
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

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	if err := db.DeleteUser(m.ID); err != nil {
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

	if p := gg.Printable(m.Name, m.Password); p >= 0 {
		log.ParamError("CreateUser: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character")
		return
	}

	if !isNsAdmin(r, m.NsID) {
		gw.WriteErr(w, r, http.StatusUnauthorized, "user is not admin for requested namespace", "ns_id", m.NsID)
		return
	}

	// check if user exists
	exists, err := db.UserExists(m.Name, m.NsID)
	if err != nil {
		log.QueryError("CreateUser: error checking user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error checking user")
		return
	}
	if exists {
		log.Data("CreateUser: user already exists")
		gw.WriteErr(w, r, http.StatusConflict, "user already exists")
		return
	}

	// create user
	u, err := db.CreateUser(m.Name, m.Password, m.NsID)
	if err != nil {
		log.QueryError("CreateUser: error creating user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error creating user")
		return
	}

	log.Result("CreateUser:", u)
	gw.WriteOK(w, "usr_id", u.ID)
}

func checkUserPassword(username, password string, nsID int64) (bool, server.User, error) {
	found, u, err := db.SelectEnabledUser(username, nsID)
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
