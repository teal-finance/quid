package api

import (
	"net/http"

	_ "github.com/lib/pq"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	"github.com/teal-finance/quid/server/db"
)

// allNsAdministrators : select all admin users for a namespace.
func allNsAdministrators(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceIDRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("AllAdministratorsInNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	data, err := db.SelectNsAdministrators(m.NamespaceID)
	if err != nil {
		log.QueryError("AllAdministratorsInNamespace: error SELECT admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT admin users")
		return
	}

	gw.WriteOK(w, data)
}

// listNonAdminUsersInNs : search from a username in namespace
func listNonAdminUsersInNs(w http.ResponseWriter, r *http.Request) {
	var m server.NonAdminUsersRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SearchForNonAdminUsersInNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if p := gg.Printable(m.Username); p >= 0 {
		log.Warn("SearchForNonAdminUsersInNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	u, err := db.SelectNonAdminUsersInNs(m.NamespaceID, m.Username)
	if err != nil {
		log.QueryError("SearchForNonAdminUsersInNamespace: error searching for non admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error searching for non admin users")
		return
	}

	gw.WriteOK(w, "users", u)
}

// CreateUserAdministrators : create admin users handler.
func createAdministrators(w http.ResponseWriter, r *http.Request) {
	var m server.AdministratorsCreation
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateAdministrators:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	for _, uID := range m.UserIDs {
		// check if user is admin
		yes, err := db.IsUserAnAdmin(m.NamespaceID, uID)
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
		if err = db.CreateAdministrator(m.NamespaceID, uID); err != nil {
			log.QueryError("CreateAdministrators: error creating admin user:", err)
			gw.WriteErr(w, r, http.StatusConflict, "error creating admin user")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// deleteAdministrator : delete an admin user handler.
func deleteAdministrator(w http.ResponseWriter, r *http.Request) {
	var m server.AdministratorDeletion
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("DeleteAdministrator:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	err := db.DeleteAdministrator(m.UserID, m.NamespaceID)
	if err != nil {
		log.QueryError("DeleteAdministrator: error deleting admin users:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error deleting admin users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
