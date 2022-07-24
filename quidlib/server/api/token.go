package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	db "github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// RequestAccessToken : request an access token from a refresh token.
func RequestAccessToken(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	refreshToken, ok := m["refresh_token"].(string)
	if !ok {
		emo.ParamError("provide a refresh_token parameter")
		gw.WriteErr(w, r, http.StatusBadRequest, "provide a refresh_token parameter")
		return
	}

	namespace, ok := m["namespace"].(string)
	if !ok {
		emo.ParamError("provide a namespace parameter")
		gw.WriteErr(w, r, http.StatusBadRequest, "provide a namespace parameter")
		return
	}

	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if !exists {
		emo.Warning("The namespace does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emo.Error(err)
		log.Fatal(err)
		return
	}

	// check if the endpoint is available
	if !ns.PublicEndpointEnabled {
		emo.Warning("Public endpoint unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
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

	if claims, ok := token.Claims.(*tokens.RefreshClaims); ok && token.Valid {
		username = claims.UserName
		fmt.Printf("%v %v", claims.UserName, claims.ExpiresAt)
	} else {
		emo.Warning(err.Error())
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	// get the user
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if !found {
		emo.Warning("User not found: " + username)
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// get the user groups names
	groupNames, err := db.SelectGroupsNamesForUser(u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emo.Error("Groups error")
		log.Fatal(err)
		return
	}

	// get the user orgs names
	orgsNames, err := db.SelectOrgsNamesForUser(u.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		emo.Error("Groups error")
		log.Fatal(err)
		return
	}

	// generate the access token
	t, err := tokens.GenAccessToken(timeout, ns.MaxTokenTTL, u.Name, groupNames, orgsNames, []byte(ns.Key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	if t == "" {
		emo.Warning("Timeout unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "error", "unauthorized")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "token", t)
}

// RequestRefreshToken : http login handler.
func RequestRefreshToken(w http.ResponseWriter, r *http.Request) {
	m := echo.Map{}
	//TODO if err := c.Bind(&m); err != nil {
	//TODO 	return
	//TODO }

	// username
	usernameParam, ok := m["username"]
	var username string
	if ok {
		username = usernameParam.(string)
	} else {
		gw.WriteErr(w, r, http.StatusBadRequest, "provide a username")
		return
	}

	// password
	passwordParam, ok := m["password"]
	var password string
	if ok {
		password = passwordParam.(string)
	} else {
		gw.WriteErr(w, r, http.StatusBadRequest, "provide a password")
		return
	}

	// namespace
	nsParam, ok := m["namespace"]
	var namespace string
	if ok {
		namespace = nsParam.(string)
	} else {
		gw.WriteErr(w, r, http.StatusBadRequest, "provide a namespace")
		return
	}

	// timeout
	timeout := chi.URLParam(r, "timeout")

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		gw.WriteErr(w, r, http.StatusBadRequest, "namespace does not exist")
		return
	}
	if !ns.PublicEndpointEnabled {
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	// check if the endpoint is available
	if !ns.PublicEndpointEnabled {
		emo.Warning("Public endpoint unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	// check if the user password matches
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		return
	}

	// Respond with unauthorized status
	if !isAuthorized {
		fmt.Println(username, "unauthorized")
		gw.WriteErr(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	// generate the token
	t, err := tokens.GenRefreshToken(timeout, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		log.Fatal(err)
	}
	if t == "" {
		gw.WriteErr(w, r, http.StatusUnauthorized, "max timeout exceeded")
		return
	}

	gw.WriteOK(w, r, http.StatusOK, "token", t)
}
