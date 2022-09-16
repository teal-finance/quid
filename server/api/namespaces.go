package api

import (
	"net/http"
	"strings"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	db "github.com/teal-finance/quid/server/db"
	"github.com/teal-finance/quid/tokens"
)

// allNamespaces : get all namespaces.
func allNamespaces(w http.ResponseWriter, r *http.Request) {
	data, err := db.SelectAllNamespaces()
	if err != nil {
		log.QueryError("AllNamespaces: error SELECT namespaces:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT namespaces")
		return
	}

	gw.WriteOK(w, data)
}

// setRefreshMaxTTL : set a max refresh token ttl for a namespace.
func setRefreshMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m server.RefreshMaxTTLRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceRefreshTokenMaxTTL:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	refreshMxTTL := m.RefreshMaxTTL

	if p := gg.Printable(refreshMxTTL); p >= 0 {
		log.Warn("SetNamespaceRefreshTokenMaxTTL: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	err := db.UpdateNsRefreshMaxTTL(id, refreshMxTTL)
	if err != nil {
		log.QueryError("SetNamespaceRefreshTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// setTokenMaxTTL : set a max access token ttl for a namespace.
func setTokenMaxTTL(w http.ResponseWriter, r *http.Request) {
	var m server.MaxTTLRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceTokenMaxTTL:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	ttl := m.MaxTTL

	if p := gg.Printable(ttl); p >= 0 {
		log.Warn("SetNamespaceTokenMaxTTL: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	err := db.UpdateNsTokenMaxTTL(id, ttl)
	if err != nil {
		log.QueryError("SetNamespaceTokenMaxTTL: error updating tokens max TTL in namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error updating tokens max TTL in namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// namespaceInfo : info about a namespace.
func namespaceInfo(w http.ResponseWriter, r *http.Request) {
	var m server.InfoRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
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

	g, err := db.SelectNsGroups(id)
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

// getAccessVerificationKey : get the key for a namespace.
func getAccessVerificationKey(w http.ResponseWriter, r *http.Request) {
	var m server.InfoRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("GetNamespaceAccessKey:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	found, algo, keyDER, err := db.SelectVerificationKeyDER(m.ID)
	if err != nil {
		log.QueryError(err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error SELECT namespace access key", "namespace_id", m.ID)
		return
	}
	if !found {
		log.QueryError("GetNamespaceAccessKey: namespace not found")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace not found", "namespace_id", m.ID)
		return
	}

	isHex := strings.ToLower(m.EncodingForm) != "base64"
	keyTxt := gg.EncodeHexOrB64Bytes(keyDER, isHex)

	gw.WriteOK(w, server.PublicKeyResponse{Alg: algo, Key: keyTxt})
}

// findNamespace : namespace creation http handler.
func findNamespace(w http.ResponseWriter, r *http.Request) {
	var m server.NameRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("FindNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	name := m.Name

	if p := gg.Printable(name); p >= 0 {
		log.Warn("FindNamespace: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	data, err := db.SelectNsStartsWith(name)
	if err != nil {
		log.QueryError("FindNamespace: error SELECT namespace:", err)
		gw.WriteErr(w, r, http.StatusInternalServerError, "error SELECT namespace")
		return
	}

	gw.WriteErr(w, r, http.StatusOK, &data)
}

// deleteNamespace : namespace creation http handler.
func deleteNamespace(w http.ResponseWriter, r *http.Request) {
	var m server.InfoRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
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

// enableNsEndpoint :.
func enableNsEndpoint(w http.ResponseWriter, r *http.Request) {
	var m server.Availability
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("SetNamespaceEndpointAvailability:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	id := m.ID
	enable := m.Enable

	err := db.EnableNsEndpoint(id, enable)
	if err != nil {
		log.QueryError("SetNamespaceEndpointAvailability: error updating namespace:", err)
		gw.WriteErr(w, r, http.StatusConflict, "error updating namespace")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createNamespace : namespace creation http handler.
func createNamespace(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceCreation
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn("CreateNamespace:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if p := gg.Printable(m.Name, m.MaxTTL, m.RefreshMaxTTL); p >= 0 {
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

	nsID, exists, err := db.CreateNamespaceIfExist(m.Name, m.MaxTTL, m.RefreshMaxTTL, m.Algo, accessKey, refreshKey, m.EnableEndpoint)
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
