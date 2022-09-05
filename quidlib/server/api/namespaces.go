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
		log.QueryError("AllNamespaces: error SELECT namespaces:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT namespaces")
		return
	}

	gw.WriteOK(w, data)
}

// SetNamespaceRefreshTokenMaxTTL : set a max refresh token ttl for a namespace.
func SetNamespaceRefreshTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m refreshMaxTTLRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceRefreshTokenMaxTTL:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	refreshMxTTL := m.RefreshMaxTTL

	if p := garcon.Printable(refreshMxTTL); p >= 0 {
		log.Warn("SetNamespaceRefreshTokenMaxTTL: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	err := db.UpdateNamespaceRefreshTokenMaxTTL(id, refreshMxTTL)
	if err != nil {
		log.QueryError("SetNamespaceRefreshTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceTokenMaxTTL : set a max access token ttl for a namespace.
func SetNamespaceTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m maxTTLRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceTokenMaxTTL:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	ttl := m.MaxTTL

	if p := garcon.Printable(ttl); p >= 0 {
		log.Warn("SetNamespaceTokenMaxTTL: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	err := db.UpdateNamespaceTokenMaxTTL(id, ttl)
	if err != nil {
		log.QueryError("SetNamespaceTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// NamespaceInfo : info about a namespace.
func NamespaceInfo(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("NamespaceInfo:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID

	nu, err := db.CountUsersForNamespace(id)
	if err != nil {
		log.QueryError("NamespaceInfo: error counting users in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting users in namespace")
		return
	}

	g, err := db.SelectGroupsForNamespace(id)
	if err != nil {
		log.QueryError("NamespaceInfo: error counting groups in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error counting groups in namespace")
		return
	}

	data := server.NamespaceInfo{
		NumUsers: nu,
		Groups:   g,
	}

	gw.WriteOK(w, data)
}

// GetNamespaceAccessVerificationKey : get the key for a namespace.
func GetNamespaceAccessVerificationKey(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("GetNamespaceAccessKey:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	found, algo, key, err := db.SelectNamespaceAccessVerificationKey(m.ID)
	if err != nil {
		log.QueryError(err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error finding namespace access key", "namespace_id", m.ID)
		return
	}
	if !found {
		log.QueryError("GetNamespaceAccessKey: namespace not found")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace not found", "namespace_id", m.ID)
		return
	}

	gw.WriteOK(w, "alg", algo, "key", key)
}

// FindNamespace : namespace creation http handler.
func FindNamespace(w http.ResponseWriter, r *http.Request) {
	var m nameRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("FindNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name

	if p := garcon.Printable(name); p >= 0 {
		log.Warn("FindNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	data, err := db.SelectNamespaceStartsWith(name)
	if err != nil {
		log.QueryError("FindNamespace: error finding namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error finding namespace")
		return
	}

	gw.WriteErr(w, r, http.StatusOK, &data)
}

// DeleteNamespace : namespace creation http handler.
func DeleteNamespace(w http.ResponseWriter, r *http.Request) {
	var m infoRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("DeleteNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID

	qRes := db.DeleteNamespace(id)
	if qRes.HasError {
		log.QueryError(qRes.Error.Message)
		if qRes.Error.HasUserMessage {
			log.Warn("DeleteNamespace: error deleting namespace")
			gw.WriteErr(w, r, http.StatusConflict, "error deleting namespace: "+qRes.Error.Message)
			return
		}
		log.Error("DeleteNamespace")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetNamespaceEndpointAvailability :.
func SetNamespaceEndpointAvailability(w http.ResponseWriter, r *http.Request) {
	var m availability
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceEndpointAvailability:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	enable := m.Enable

	err := db.SetNamespaceEndpointAvailability(id, enable)
	if err != nil {
		log.QueryError("SetNamespaceEndpointAvailability: error updating namespace:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error updating namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateNamespace : namespace creation http handler.
func CreateNamespace(w http.ResponseWriter, r *http.Request) {
	var m namespaceCreation
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if p := garcon.Printable(m.Name, m.MaxTTL, m.RefreshMaxTTL); p >= 0 {
		log.Warn("CreateNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	if m.Algo == "" {
		m.Algo = "HS256"
		log.Param("No signing algo provided, defaults to " + m.Algo)
	}

	refreshKey := tokens.GenerateKeyHMAC(256)
	accessKey, err := tokens.GenerateSigningKey(m.Algo)
	if err != nil {
		log.Warn("Generate AccessKey algo=" + m.Algo + " err: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nsID, exists, err := createNamespace(m.Name, m.MaxTTL, m.RefreshMaxTTL, m.Algo, accessKey, refreshKey, m.EnableEndpoint)
	if err != nil {
		gw.WriteErr(w, r, http.StatusInternalServerError, "error creating namespace", "namespace", m.Name)
		return
	}
	if exists {
		gw.WriteErr(w, r, http.StatusConflict, "namespace already exists", "namespace", m.Name)
		return
	}

	gw.WriteOK(w, "namespace_id", nsID)
}

// createNamespace : create a namespace.
func createNamespace(name, ttl, refreshMaxTTL, algo string, accessKey, refreshKey []byte, endpoint bool) (int64, bool, error) {
	exists, err := db.NamespaceExists(name)
	if err != nil {
		log.QueryError("createNamespace NamespaceExists:", err)
		return 0, false, err
	}
	if exists {
		log.QueryError("createNamespace: already exist")
		return 0, true, nil
	}

	nsID, err := db.CreateNamespace(name, ttl, refreshMaxTTL, algo, accessKey, refreshKey, endpoint)
	if err != nil {
		log.QueryError("createNamespace:", err)
		return 0, false, err
	}

	return nsID, false, nil
}
