package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/teal-finance/quid/quidlib/server"
	db "github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// AllNamespaces : get all namespaces.
func AllNamespaces(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllNamespaces()
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting namespaces")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, &data)
}

// SetNamespaceRefreshTokenMaxTTL : set a max refresh token ttl for a namespace.
func SetNamespaceRefreshTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))
	refreshMxTTL := m["refresh_max_ttl"].(string)

	err := db.UpdateNamespaceRefreshTokenMaxTTL(id, refreshMxTTL)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max ttl in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceTokenMaxTTL : set a max access token ttl for a namespace.
func SetNamespaceTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))
	ttl := m["max_ttl"].(string)

	err := db.UpdateNamespaceTokenMaxTTL(id, ttl)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max ttl in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// NamespaceInfo : info about a namespace.
func NamespaceInfo(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))

	nu, err := db.CountUsersForNamespace(id)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting users in namespace")
		return
	}

	g, err := db.SelectGroupsForNamespace(id)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting users in namespace")
		return
	}

	data := server.NamespaceInfo{
		NumUsers: nu,
		Groups:   g,
	}

	gw.WriteOK(w, r, http.StatusOK, &data)
}

// GetNamespaceKey : get the key for a namespace.
func GetNamespaceKey(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))

	found, data, err := db.SelectNamespaceKey(id)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding namespace key")
		return
	}
	if !found {
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace not found")
		return
	}

	gw.WriteOK(w, "key", data)
}

// FindNamespace : namespace creation http handler.
func FindNamespace(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	name := m["name"].(string)

	data, err := db.SelectNamespaceStartsWith(name)
	if err != nil {
		emo.QueryError(err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding namespace")
		return
	}

	gw.WriteErr(w, r, http.StatusOK, &data)
}

// DeleteNamespace : namespace creation http handler.
func DeleteNamespace(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))

	qRes := db.DeleteNamespace(id)
	if qRes.HasError {
		emo.QueryError(qRes.Error.Message)
		if qRes.Error.HasUserMessage {
			gw.WriteErr(w, r, http.StatusConflict, "error deleting namespace: "+qRes.Error.Message)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceEndpointAvailability :.
func SetNamespaceEndpointAvailability(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	id := int64(m["id"].(float64))
	enable := m["enable"].(bool)

	err := db.SetNamespaceEndpointAvailability(id, enable)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error updating namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateNamespace : namespace creation http handler.
func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	name := m["name"].(string)
	maxTTL := m["max_ttl"].(string)
	refreshMaxTTL := m["refresh_max_ttl"].(string)
	enableEndpoint := m["enable_endpoint"].(bool)
	key := tokens.GenKey()
	refreshKey := tokens.GenKey()

	nsID, exists, err := createNamespace(name, key, refreshKey, maxTTL, refreshMaxTTL, enableEndpoint)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error creating namespace")
		return
	}
	if exists {
		gw.WriteErr(w, r, http.StatusConflict, "namespace already exists")
		return
	}

	gw.WriteOK(w, "namespace_id", nsID)
}

// createNamespace : create a namespace.
func createNamespace(name, key, refreshKey, ttl, refreshMaxTTL string, endpoint bool) (int64, bool, error) {
	exists, err := db.NamespaceExists(name)
	if err != nil {
		return 0, false, err
	}
	if exists {
		return 0, true, nil
	}

	nsID, err := db.CreateNamespace(name, key, refreshKey, ttl, refreshMaxTTL, endpoint)
	if err != nil {
		return 0, false, err
	}

	return nsID, false, nil
}
