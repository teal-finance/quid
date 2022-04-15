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
	if ns.Name == "quid" {
		// check if the user admin group
		isAdmin, err := IsUserInAdminGroup(u.ID, ns.ID)
		if err != nil {
			return "", err
		}
		if !isAdmin {
			return "", errors.New("the user is not quid admin")
		}
	} else {
		isAdmin, err := AdministratorExists(u.ID, ns.ID)
		if err != nil {
			return "", err
		}
		if !isAdmin {
			return "", errors.New("the user is not namespace admin")
		}
	}

	// get the refresh token
	token, err := tokens.GenRefreshToken(ns.MaxTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, u.Name, []byte(ns.RefreshKey))
	if err != nil {
		msg := "Error generating refresh token"
		emo.Error(msg, err)
		return "", err
	}
	return token, nil
}
