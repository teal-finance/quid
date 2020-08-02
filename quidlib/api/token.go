package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/tokens"
)

// RequestAccessToken : request an access token from a refresh token
func RequestAccessToken(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	refreshToken, ok := m["refresh_token"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "provide a refresh_token parameter",
		})
	}
	namespace, ok := m["namespace"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "provide a namespace parameter",
		})
	}
	timeout := c.Param("timeout")

	// get the refresh key
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if !exists {
		emo.Error("The refresh key does not exist")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": true,
		})
	}
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	// verify the refresh token
	var username string
	token, err := jwt.ParseWithClaims(refreshToken, &tokens.StandardRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ns.RefreshKey), nil
	})
	if claims, ok := token.Claims.(*tokens.StandardRefreshClaims); ok && token.Valid {
		username = claims.Name
		fmt.Printf("%v %v", claims.Name, claims.StandardClaims.ExpiresAt)
	} else {
		emo.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	// get the user
	found, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if !found {
		emo.Error("User not found: " + username)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	// get the user group names
	groupNames, err := db.SelectGroupsNamesForUser(u.ID)
	if err != nil {
		emo.Error("Groups error")
		log.Fatal(err)
		return err
	}

	// generate the access token
	isAuth, t, err := tokens.GenAccessToken(ns, u.Name, groupNames, timeout, ns.MaxTokenTTL)
	if !isAuth {
		emo.Error("Timeout unauthorized")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": true,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// RequestRefreshToken : http login handler
func RequestRefreshToken(c echo.Context) error {

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	username := m["username"].(string)
	password := m["password"].(string)
	namespace := m["namespace"].(string)
	timeout := c.Param("timeout")

	// get the namespace
	exists, ns, err := db.SelectNamespaceFromName(namespace)
	if err != nil {
		return err
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "namespace does not exist",
		})
	}

	// check if the user password matches
	isAuthorized, u, err := checkUserPassword(username, password, ns.ID)
	if err != nil {
		return err
	}
	// Respond with unauthorized status
	if isAuthorized == false {
		fmt.Println(username, "unauthorized")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	// generate the token
	isAuth, t, err := tokens.GenRefreshToken(ns, u.Name, timeout)
	if err != nil {
		log.Fatal(err)
	}
	if !isAuth {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "max timeout exceeded",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
