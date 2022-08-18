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
		emo.QueryError("AllNamespaces: error selecting namespaces:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error selecting namespaces")
		return
	}

	gw.WriteOK(w, data)
}

// SetNamespaceRefreshTokenMaxTTL : set a max refresh token ttl for a namespace.
func SetNamespaceRefreshTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m refreshMaxTTLRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("SetNamespaceRefreshTokenMaxTTL:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	refreshMxTTL := m.RefreshMaxTTL

	if p := garcon.Printable(refreshMxTTL); p >= 0 {
		emo.Warning("SetNamespaceRefreshTokenMaxTTL: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := db.UpdateNamespaceRefreshTokenMaxTTL(id, refreshMxTTL)
	if err != nil {
		emo.QueryError("SetNamespaceRefreshTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceTokenMaxTTL : set a max access token ttl for a namespace.
func SetNamespaceTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m maxTTLRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("SetNamespaceTokenMaxTTL:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	ttl := m.MaxTTL

	if p := garcon.Printable(ttl); p >= 0 {
		emo.Warning("SetNamespaceTokenMaxTTL: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := db.UpdateNamespaceTokenMaxTTL(id, ttl)
	if err != nil {
		emo.QueryError("SetNamespaceTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// NamespaceInfo : info about a namespace.
func NamespaceInfo(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("NamespaceInfo:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

	nu, err := db.CountUsersForNamespace(id)
	if err != nil {
		emo.QueryError("NamespaceInfo: error counting users in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting users in namespace")
		return
	}

	g, err := db.SelectGroupsForNamespace(id)
	if err != nil {
		emo.QueryError("NamespaceInfo: error counting groups in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting groups in namespace")
		return
	}

	data := server.NamespaceInfo{
		NumUsers: nu,
		Groups:   g,
	}

	gw.WriteOK(w, data)
}

// GetNamespaceKey : get the key for a namespace.
func GetNamespaceKey(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("GetNamespaceKey:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

	found, data, err := db.SelectNamespaceKey(id)
	if err != nil {
		emo.QueryError("GetNamespaceKey: error finding namespace key:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding namespace key")
		return
	}
	if !found {
		emo.QueryError("GetNamespaceKey: namespace not found")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace not found")
		return
	}

	gw.WriteOK(w, "key", data)
}

// FindNamespace : namespace creation http handler.
func FindNamespace(w http.ResponseWriter, r *http.Request) {
	var m nameRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("FindNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name

	if p := garcon.Printable(name); p >= 0 {
		emo.Warning("FindNamespace: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := db.SelectNamespaceStartsWith(name)
	if err != nil {
		emo.QueryError("FindNamespace: error finding namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding namespace")
		return
	}

	gw.WriteErr(w, r, http.StatusOK, &data)
}

// DeleteNamespace : namespace creation http handler.
func DeleteNamespace(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("DeleteNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID

	qRes := db.DeleteNamespace(id)
	if qRes.HasError {
		emo.QueryError(qRes.Error.Message)
		if qRes.Error.HasUserMessage {
			emo.Warning("DeleteNamespace: error deleting namespace")
			gw.WriteErr(w, r, http.StatusConflict, "error deleting namespace: "+qRes.Error.Message)
			return
		}
		emo.Error("DeleteNamespace")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceEndpointAvailability :.
func SetNamespaceEndpointAvailability(w http.ResponseWriter, r *http.Request) {
	var m availability
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("SetNamespaceEndpointAvailability:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := m.ID
	enable := m.Enable

	err := db.SetNamespaceEndpointAvailability(id, enable)
	if err != nil {
		emo.Warning("SetNamespaceEndpointAvailability: error updating namespace:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error updating namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateNamespace : namespace creation http handler.
func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceCreation
	if err := garcon.DecodeJSONBody(r, &m); err != nil {
		emo.Warning("CreateNamespace:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := m.Name
	maxTTL := m.MaxTTL
	refreshMaxTTL := m.RefreshMaxTTL
	enableEndpoint := m.EnableEndpoint

	if p := garcon.Printable(name, maxTTL, refreshMaxTTL); p >= 0 {
		emo.Warning("CreateNamespace: JSON contains a forbidden character at p=", p)
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
		emo.QueryError("createNamespace NamespaceExists:", err)
		return 0, false, err
	}
	if exists {
		emo.QueryError("createNamespace: already exist")
		return 0, true, nil
	}

	nsID, err := db.CreateNamespace(name, key, refreshKey, ttl, refreshMaxTTL, endpoint)
	if err != nil {
		emo.QueryError("createNamespace:", err)
		return 0, false, err
	}

	return nsID, false, nil
}
