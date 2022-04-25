package db

import (
	"errors"

	"github.com/teal-finance/quid/quidlib/tokens"
)

// GenNsAdminTokenForUser : generate a refresh token for an admin user and namespace
func GenNsAdminTokenForUser(userName string, nsName string) (string, error) {
	// get the namespace
	_, ns, err := SelectNamespaceFromName(nsName)
	if err != nil {
		return "", err
	}
	_, u, err := SelectNonDisabledUser(userName, ns.ID)
	if err != nil {
		return "", err
	}

	// check admin perms
	isAdmin, err := IsUserAdmin(nsName, ns.ID, u.ID)
	if err != nil {
		return "", err
	}
	if !isAdmin {
		return "", errors.New("the user is not namespace admin")
	}

	// get the refresh token
	emo.Encrypt("Gen token", ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	token, err := tokens.GenRefreshToken(ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		msg := "Error generating refresh token"
		emo.Error(msg, err)
		return "", err
	}
	return token, nil
}
