package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"

	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/server"
	db "github.com/teal-finance/quid/server/db"
	"github.com/teal-finance/quid/tokens"
)

// requestAccessToken : request an access token from a refresh token.
func requestAccessToken(w http.ResponseWriter, r *http.Request) {
	var m server.AccessTokenRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RequestAccessToken:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	refreshToken := m.RefreshToken
	namespace := m.Namespace

	if p := gg.Printable(refreshToken, namespace); p >= 0 {
		log.Warn("RequestAccessToken: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNsFromName(namespace)
	if err != nil {
		log.QueryError("RequestAccessToken SelectNsFromName:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "DB error SELECT namespace", "namespace", namespace)
		return
	}
	if !exists {
		log.Data("RequestAccessToken: the namespace does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the endpoint is available
	if !ns.PublicEndpointEnabled {
		log.Warn("RequestAccessToken: Public endpoint unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// verify the refresh token
	var username string
	token, err := jwt.ParseWithClaims(refreshToken, &tokens.RefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ns.RefreshKey), nil
	})
	if err != nil {
		log.Warn("RequestAccessToken ParseWithClaims:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !token.Valid {
		log.Warn("RequestAccessToken: invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*tokens.RefreshClaims)
	if !ok {
		log.Error("RequestAccessToken: cannot convert to RefreshClaims")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username = claims.UserName
	log.AccessToken("RequestAccessToken:", claims.UserName, claims.ExpiresAt)

	// get the user
	found, u, err := db.SelectEnabledUser(username, ns.ID)
	if err != nil {
		log.QueryError(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	if !found {
		log.Warn("RequestAccessToken: user not found: " + username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get the user groups names
	groupNames, err := db.SelectGroupsNamesForUser(u.ID)
	if err != nil {
		log.QueryError("RequestAccessToken: Groups error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get the user orgs names
	orgsNames, err := db.SelectOrgsNamesForUser(u.ID)
	if err != nil {
		log.QueryError("RequestAccessToken: Orgs error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate the access token
	var t string
	if ns.SigningAlgo == "HS256" {
		t, err = tokens.GenAccessToken(timeout, ns.MaxTokenTTL, u.Name, groupNames, orgsNames, []byte(ns.AccessKey))
	} else {
		t, err = tokens.GenAccessTokenWithAlgo(ns.SigningAlgo, timeout, ns.MaxTokenTTL, u.Name, groupNames, orgsNames, ns.AccessKey)
	}

	if err != nil {
		log.Error("RequestAccessToken GenAccessToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if t == "" {
		log.Warn("RequestAccessToken: Timeout unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "error", "unauthorized")
		return
	}

	log.AccessToken("RequestAccessToken: user="+u.Name+" t="+timeout+" TTL="+ns.MaxTokenTTL+" grp=", groupNames, "org=", orgsNames)
	gw.WriteOK(w, "token", t)
}

func getAccessPublicKey(w http.ResponseWriter, r *http.Request) {
	var m server.NamespaceRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.Warn(err)
		gw.WriteErr(w, r, http.StatusBadRequest, "cannot decode JSON")
		return
	}

	if p := gg.Printable(m.Namespace); p >= 0 {
		log.Warn(`JSON {"name":....} has forbidden character at p=`, p)
		gw.WriteErr(w, r, http.StatusBadRequest, "forbidden character", "position", p)
		return
	}

	// get the namespace
	exists, ns, err := db.SelectNsFromName(m.Namespace)
	if err != nil {
		log.QueryError(err)
		gw.WriteErr(w, r, http.StatusBadRequest, "DB error SELECT namespace", "namespace", m.Namespace)
		return
	}
	if !exists {
		log.ParamError("Namespace", m.Namespace, "does not exist")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist", "namespace", m.Namespace)
		return
	}

	switch ns.SigningAlgo {
	case "RS256", "RS384", "RS512": // OK
	case "PS256", "PS384", "PS512": // OK
	case "ES256", "ES384", "ES512": // OK
	case "EdDSA": // OK
	default: // "HS256", "HS384", "HS512"
		log.ParamError("Namespace", m.Namespace, "has algo", ns.SigningAlgo, "without public key")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace signing algo has no public key", "algo", ns.SigningAlgo)
	}

	publicDER, err := tokens.DecryptVerificationKeyDER(ns.SigningAlgo, ns.AccessKey)
	if err != nil {
		log.Error(err)
	}

	keyB64 := base64.RawURLEncoding.EncodeToString(publicDER)
	gw.WriteOK(w, server.PublicKeyResponse{Alg: ns.SigningAlgo, Key: keyB64})
}

func validAccessToken(w http.ResponseWriter, r *http.Request) {
	var m server.AccessTokenValidationRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError(err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	if p := gg.Printable(m.AccessToken, m.Namespace); p >= 0 {
		log.Warn("JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	// get the namespace
	exists, ns, err := db.SelectNsFromName(m.Namespace)
	if err != nil {
		log.QueryError(err)
		gw.WriteErr(w, r, http.StatusBadRequest, "DB error SELECT namespace", "namespace", m.Namespace)
		return
	}
	if !exists {
		log.ParamError("Namespace", m.Namespace, "does not exist")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist", "namespace", m.Namespace)
		return
	}

	verificationKeyDER, err := tokens.DecryptVerificationKeyDER(ns.SigningAlgo, ns.AccessKey)
	if err != nil {
		gw.WriteErr(w, r, http.StatusBadRequest, "error decrypting verification key", "namespace", m.Namespace)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = tokens.ValidAccessToken(m.AccessToken, ns.SigningAlgo, verificationKeyDER)
	if err != nil {
		log.Result("Invalid AccessToken:", err)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"valid":false}`))
	} else {
		log.Result("Valid AccessToken")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"valid":true}`))
	}
}

// requestRefreshToken : http login handler.
func requestRefreshToken(w http.ResponseWriter, r *http.Request) {
	var m server.PasswordRequest
	if err := gg.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RequestRefreshToken:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "cannot decode JSON")
		return
	}

	username := m.Username
	password := m.Password
	namespace := m.Namespace

	if p := gg.Printable(username, password, namespace); p >= 0 {
		log.ParamError("RequestRefreshToken: JSON contains a forbidden character at p=", p)
		gw.WriteErr(w, r, http.StatusUnauthorized, "forbidden character", "position", p)
		return
	}

	// timeout
	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNsFromName(namespace)
	if err != nil {
		log.QueryError("RequestRefreshToken SelectNsFromName:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "DB error SELECT namespace", "namespace", namespace)
		return
	}
	if !exists {
		log.Data("RequestRefreshToken: namespace does not exist")
		gw.WriteErr(w, r, http.StatusUnauthorized, "namespace does not exist")
		return
	}

	// check if the endpoint is available
	if !ns.PublicEndpointEnabled {
		log.Data("RequestRefreshToken: public endpoint unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "endpoint disabled", "namespace", namespace)
		return
	}

	// check if the user password matches
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		gw.WriteErr(w, r, http.StatusUnauthorized, "error while checking password", "namespace_id", ns.ID, "usr", username)
		return
	}
	if !isAuthorized {
		log.Info("RequestRefreshToken u=" + username + ": disabled user or bad password")
		gw.WriteErr(w, r, http.StatusUnauthorized, "disabled user or bad password", "namespace_id", ns.ID, "usr", username)
		return
	}

	// generate the token
	t, err := tokens.GenRefreshToken(timeout, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		log.Error("RequestRefreshToken GenRefreshToken:", err)
		gw.WriteErr(w, r, http.StatusUnauthorized, "error while generating a new RefreshToken", "namespace", ns.Name, "usr", u.Name)
		return
	}
	if t == "" {
		log.Info("RequestRefreshToken: max timeout exceeded")
		gw.WriteErr(w, r, http.StatusUnauthorized, "max timeout exceeded", "timeout", timeout, "max", ns.MaxRefreshTokenTTL)
		return
	}

	log.RefreshToken("RequestRefreshToken: user=" + u.Name + " t=" + timeout + " TTL=" + ns.MaxRefreshTokenTTL + " ns=" + ns.Name)
	gw.WriteOK(w, "token", t)
}
