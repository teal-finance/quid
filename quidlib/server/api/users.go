package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
)

// AllUsersInNamespace : select all users for a namespace.
func AllUsersInNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("AllUsersInNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := db.SelectUsersInNamespace(nsID)
	if err != nil {
		emo.QueryError("AllUsersInNamespace: error selecting users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting users")
		return
	}

	b, err := json.Marshal(&data)
	if err != nil {
		emo.Error("AllUsersInNamespace: %v while serializing %+v", err, data)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GroupsForNamespace : get the groups of a user.
func GroupsForNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("GroupsForNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namespace := m.Namespace

	if p := garcon.Printable(namespace); p >= 0 {
		emo.ParamError("GroupsForNamespace: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hasResult, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil || !hasResult {
		emo.QueryError("GroupsForNamespace: error selecting namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting namespace")
		return
	}

	g, err := db.SelectGroupsForNamespace(ns.ID)
	if err != nil {
		emo.QueryError("GroupsForNamespace: error selecting groups:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting groups")
	}

	emo.Result("GroupsForNamespace:", g)
	gw.WriteOK(w, "groups", g)
}

// AddUserInOrg : add a user in an org.
func AddUserInOrg(w http.ResponseWriter, r *http.Request) {
	var m userOrgRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("AddUserInOrg:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	oID := m.OrgID

	err := db.AddUserInOrg(uID, oID)
	if err != nil {
		emo.QueryError("AddUserInOrg: error adding user in org:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error adding user in org")
		return
	}

	emo.Result("AddUserInOrg OK")
	w.WriteHeader(http.StatusOK)
}

// RemoveUserFromOrg : add a user in an org.
func RemoveUserFromOrg(w http.ResponseWriter, r *http.Request) {
	var m userOrgRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("RemoveUserFromOrg:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	oID := m.OrgID

	err := db.RemoveUserFromOrg(uID, oID)
	if err != nil {
		emo.QueryError("RemoveUserFromOrg: error removing user from org:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error removing user from org")
		return
	}

	emo.Result("RemoveUserFromOrg OK")
	w.WriteHeader(http.StatusOK)
}

// AddUserInGroup : add a user in a group.
func AddUserInGroup(w http.ResponseWriter, r *http.Request) {
	var m userGroupRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("AddUserInGroup:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	gID := m.GroupID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := db.AddUserInGroup(uID, gID)
	if err != nil {
		emo.QueryError("AddUserInGroup:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error adding user in group")
		return
	}

	emo.Result("AddUserInGroup OK")
	w.WriteHeader(http.StatusOK)
}

// RemoveUserFromGroup : add a user in a group.
func RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	var m userGroupRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("RemoveUserFromGroup:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	gID := m.GroupID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := db.RemoveUserFromGroup(uID, gID)
	if err != nil {
		emo.QueryError("RemoveUserFromGroup: error removing user from group:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error removing user from group")
		return
	}

	emo.Result("RemoveUserFromGroup OK")
	w.WriteHeader(http.StatusOK)
}

// UserGroupsInfo : get info for a user.
func UserGroupsInfo(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("UserGroupsInfo:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	g, err := db.SelectGroupsForUser(id)
	if err != nil {
		emo.QueryError("UserGroupsInfo: error selecting groups:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting groups")
		return
	}

	emo.Result("UserGroupsInfo:", g)
	gw.WriteOK(w, "groups", g)
}

// DeleteUser : delete a user handler.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var m userRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("DeleteUser:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	nsID := m.NamespaceID

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := db.DeleteUser(id); err != nil {
		emo.QueryError("DeleteUser: error deleting user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error deleting user")
		return
	}

	emo.Result("DeleteUser OK")
	gw.WriteOK(w, "message", "ok")
}

// CreateUserHandler : create a user handler.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var m userHandlerCreation
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.ParamError("CreateUserHandler:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name
	password := m.Password
	nsID := m.NamespaceID

	if p := garcon.Printable(name, password); p >= 0 {
		emo.ParamError("CreateUserHandler: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !VerifyAdminNs(w, r, nsID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if user exists
	exists, err := db.UserNameExists(name, nsID)
	if err != nil {
		emo.QueryError("CreateUserHandler: error checking user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error checking user")
		return
	}
	if exists {
		emo.Data("CreateUserHandler: error user already exist")
		gw.WriteErr(w, r, http.StatusConflict, "error user already exist")
		return
	}

	// create user
	u, err := db.CreateUser(name, password, nsID)
	if err != nil {
		emo.QueryError("CreateUserHandler: error creating user:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error creating user")
		return
	}

	emo.Result("CreateUserHandler:", u)
	gw.WriteOK(w, "user_id", u.ID)
}

func checkUserPassword(username, password string, namespaceID int64) (bool, server.User, error) {
	found, u, err := db.SelectNonDisabledUser(username, namespaceID)
	if !found || err != nil {
		return false, u, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	if err != nil {
		fmt.Println("ERROR", "|"+err.Error()+"|")
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return false, u, err
		}
	}

	return true, u, err
}
