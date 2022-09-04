package api

import (
	"net/http"

	_ "github.com/lib/pq"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// AllAdministratorsInNamespace : select all admin users for a namespace.
func AllAdministratorsInNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("AllAdministratorsInNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID := m.NamespaceID

	data, err := db.SelectAdministratorsInNamespace(nsID)
	if err != nil {
		log.QueryError("AllAdministratorsInNamespace: error selecting admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting admin users")
		return
	}

	gw.WriteOK(w, data)
}

// SearchForNonAdminUsersInNamespace : search from a username in namespace
func SearchForNonAdminUsersInNamespace(w http.ResponseWriter, r *http.Request) {
	var m nonAdminUsersRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SearchForNonAdminUsersInNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := m.Username
	nsID := m.NamespaceID

	if p := garcon.Printable(username); p >= 0 {
		log.Warn("SearchForNonAdminUsersInNamespace: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := db.SearchForNonAdminUsersInNamespace(nsID, username)
	if err != nil {
		log.QueryError("SearchForNonAdminUsersInNamespace: error searching for non admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error searching for non admin users")
		return
	}

	gw.WriteOK(w, "users", u)
}

// CreateUserAdministrators : create admin users handler.
func CreateAdministrators(w http.ResponseWriter, r *http.Request) {
	var m administratorsCreation
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateAdministrators:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uIDs := m.UserIDs
	nsID := m.NamespaceID

	for _, uID := range uIDs {
		// check if user exists
		exists, err := db.AdministratorExists(nsID, uID)
		if err != nil {
			log.QueryError("CreateAdministrators: error checking admin user:", err)
			gw.WriteErr(w, r, http.StatusConflict, "error checking admin user")
			return
		}
		if exists {
			log.QueryError("CreateAdministrators: admin user already exist:", err)
			gw.WriteErr(w, r, http.StatusConflict, "admin user already exist")
			return
		}

		// create admin user
		if _, err = db.CreateAdministrator(nsID, uID); err != nil {
			log.QueryError("CreateAdministrators: error creating admin user:", err)
			gw.WriteErr(w, r, http.StatusConflict, "error creating admin user")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAdministrator : delete an admin user handler.
func DeleteAdministrator(w http.ResponseWriter, r *http.Request) {
	var m administratorDeletion
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("DeleteAdministrator:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	nsID := m.NamespaceID

	err := db.DeleteAdministrator(uID, nsID)
	if err != nil {
		log.QueryError("DeleteAdministrator: error deleting admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error deleting admin users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
