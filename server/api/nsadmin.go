package api

import (
	"net/http"

	_ "github.com/lib/pq"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server/db"
)

// allNsAdministrators : select all admin users for a namespace.
func allNsAdministrators(w http.ResponseWriter, r *http.Request) {
	var m namespaceIDRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("AllAdministratorsInNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	nsID := m.NamespaceID

	data, err := db.SelectNsAdministrators(nsID)
	if err != nil {
		log.QueryError("AllAdministratorsInNamespace: error SELECT admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT admin users")
		return
	}

	gw.WriteOK(w, data)
}

// listNonAdminUsersInNs : search from a username in namespace
func listNonAdminUsersInNs(w http.ResponseWriter, r *http.Request) {
	var m nonAdminUsersRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SearchForNonAdminUsersInNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	username := m.Username
	nsID := m.NamespaceID

	if p := gg.Printable(username); p >= 0 {
		log.Warn("SearchForNonAdminUsersInNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	u, err := db.SelectNonAdminUsersInNs(nsID, username)
	if err != nil {
		log.QueryError("SearchForNonAdminUsersInNamespace: error searching for non admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error searching for non admin users")
		return
	}

	gw.WriteOK(w, "users", u)
}

// CreateUserAdministrators : create admin users handler.
func createAdministrators(w http.ResponseWriter, r *http.Request) {
	var m administratorsCreation
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateAdministrators:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	uIDs := m.UserIDs
	nsID := m.NamespaceID

	for _, uID := range uIDs {
		// check if user is admin
		yes, err := db.IsUserAnAdmin(nsID, uID)
		if err != nil {
			log.QueryError("CreateAdministrators: error checking admin user:", err)
			gw.WriteErr(w, r, http.StatusConflict, "error checking admin user")
			return
		}
		if yes {
			log.Query("CreateAdministrators: user is already an administrator:", err)
			continue
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

// deleteAdministrator : delete an admin user handler.
func deleteAdministrator(w http.ResponseWriter, r *http.Request) {
	var m administratorDeletion
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("DeleteAdministrator:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
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
