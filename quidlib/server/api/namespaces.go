package api

import (
	"net/http"

	"github.com/teal-finance/garcon"
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
	var m RefreshMaxTTLRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	refreshMxTTL := m.RefreshMaxTTL

	if p := garcon.Printable(refreshMxTTL); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	var m MaxTTLRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	ttl := m.MaxTTL

	if p := garcon.Printable(ttl); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	var m InfoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

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
	var m InfoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

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
	var m NameRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name

	if p := garcon.Printable(name); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	var m InfoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

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
	var m Availability
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	enable := m.Enable

	err := db.SetNamespaceEndpointAvailability(id, enable)
	if err != nil {
		gw.WriteErr(w, r, http.StatusConflict, "error updating namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateNamespace : namespace creation http handler.
func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	var m NamespaceCreation
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name
	maxTTL := m.MaxTTL
	refreshMaxTTL := m.RefreshMaxTTL
	enableEndpoint := m.EnableEndpoint

	if p := garcon.Printable(name, maxTTL, refreshMaxTTL); p >= 0 {
		emo.Warning("JSON contains a forbidden character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
