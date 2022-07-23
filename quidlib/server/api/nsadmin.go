package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/teal-finance/quid/quidlib/server/db"
)

// AllAdministratorsInNamespace : select all admin users for a namespace.
func AllAdministratorsInNamespace(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	nsID := int64(m["namespace_id"].(float64))

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
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	username := m["username"].(string)
	nsID := int64(m["namespace_id"].(float64))

	u, err := db.SearchForNonAdminUsersInNamespace(nsID, username)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error searching for non admin users")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "users", u)
}

// CreateUserAdministrators : create admin users handler.
func CreateAdministrators(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	uIDs := m["user_ids"].([]any)
	nsID := int64(m["namespace_id"].(float64))

	for _, fuserID := range uIDs {
		uID := int64(fuserID.(float64))

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
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	uID := int64(m["user_id"].(float64))
	nsID := int64(m["namespace_id"].(float64))

	err := db.DeleteAdministrator(uID, nsID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error deleting admin users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
