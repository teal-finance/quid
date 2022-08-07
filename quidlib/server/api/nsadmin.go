package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/server/db"
)

// AllAdministratorsInNamespace : select all admin users for a namespace.
func AllAdministratorsInNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID := m.NamespaceID

	data, err := db.SelectAdministratorsInNamespace(nsID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting admin users")
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

// SearchForNonAdminUsersInNamespace : search from a username in namespace
func SearchForNonAdminUsersInNamespace(w http.ResponseWriter, r *http.Request) {
	var m nonAdminUsersRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := m.Username
	nsID := m.NamespaceID

	if p := garcon.Printable(username); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := db.SearchForNonAdminUsersInNamespace(nsID, username)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error searching for non admin users")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "users", u)
}

// CreateUserAdministrators : create admin users handler.
func CreateAdministrators(w http.ResponseWriter, r *http.Request) {
	var m administratorsCreation
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uIDs := m.UserIDs
	nsID := m.NamespaceID

	for _, uID := range uIDs {
		// check if user exists
		exists, err := db.AdministratorExists(nsID, uID)
		if err != nil {
			gw.WriteErr(w, r, http.StatusConflict, "error checking admin user")
			return
		}
		if exists {
			gw.WriteErr(w, r, http.StatusConflict, "error creating admin user")
			return
		}

		// create admin user
		if _, err = db.CreateAdministrator(nsID, uID); err != nil {
			gw.WriteErr(w, r, http.StatusConflict, "error creating admin user")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAdministrator : delete an admin user handler.
func DeleteAdministrator(w http.ResponseWriter, r *http.Request) {
	var m administratorDeletion
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uID := m.UserID
	nsID := m.NamespaceID

	err := db.DeleteAdministrator(uID, nsID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error deleting admin users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
