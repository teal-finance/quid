package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"

	"github.com/teal-finance/garcon"
	db "github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// RequestAccessToken : request an access token from a refresh token.
func RequestAccessToken(w http.ResponseWriter, r *http.Request) {
	var m accessTokenRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RequestAccessToken:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := m.RefreshToken
	namespace := m.Namespace

	if p := garcon.Printable(refreshToken, namespace); p >= 0 {
		log.Warn("RequestAccessToken: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		log.QueryError("RequestAccessToken SelectNamespaceFromName:", err)
		w.WriteHeader(http.StatusInternalServerError)
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
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if err != nil {
		log.QueryError("RequestAccessToken SelectNonDisabledUser:", err)
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
		t, err = tokens.NewAccessToken(timeout, ns.MaxTokenTTL, u.Name, groupNames, orgsNames, ns.SigningAlgo, ns.AccessKey)
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

// RequestRefreshToken : http login handler.
func RequestRefreshToken(w http.ResponseWriter, r *http.Request) {
	var m passwordRequest
	if err := garcon.UnmarshalJSONRequest(w, r, &m); err != nil {
		log.ParamError("RequestRefreshToken:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := m.Username
	password := m.Password
	namespace := m.Namespace

	if p := garcon.Printable(username, password, namespace); p >= 0 {
		log.ParamError("RequestRefreshToken: JSON contains a forbidden character at p=", p)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// timeout
	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		log.QueryError("RequestRefreshToken SelectNamespaceFromName:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Data("RequestRefreshToken: namespace does not exist")
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist")
		return
	}

	// check if the endpoint is available
	if !ns.PublicEndpointEnabled {
		log.Data("RequestRefreshToken: public endpoint unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user password matches
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		return
	}

	// Respond with unauthorized status
	if !isAuthorized {
		log.Info("RequestRefreshToken: u=" + username + " unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// generate the token
	t, err := tokens.GenRefreshToken(timeout, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		log.Error("RequestRefreshToken GenRefreshToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if t == "" {
		log.Info("RequestRefreshToken: max timeout exceeded")
		gw.WriteErr(w, r, http.StatusUnauthorized, "max timeout exceeded")
		return
	}

	log.RefreshToken("RequestRefreshToken: user=" + u.Name + " t=" + timeout + " TTL=" + ns.MaxRefreshTokenTTL + " ns=" + ns.Name)
	gw.WriteOK(w, "token", t)
}
