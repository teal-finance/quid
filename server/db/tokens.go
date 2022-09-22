package db

import (
	"errors"

	"github.com/teal-finance/quid/tokens"
)

// GenNsAdminTokenForUser : generate a refresh token for an admin user and namespace
func GenNsAdminTokenForUser(userName, nsName string) (string, error) {
	log.Info("Generating ns admin token for", userName, nsName)

	// get the namespace
	ns, err := SelectNsFromName(nsName)
	if err != nil {
		return "", err
	}

	exists, uid, err := SelectEnabledUsrID(userName)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("user does not exist")
	}

	// check admin perms
	adminType, err := GetUserType(nsName, ns.ID, uid)
	if err != nil {
		return "", err
	}
	// emo.Debug("USER", userName, uid, "NS", nsName, "NSa", isNsAdmin)
	if adminType == UserNoAdmin {
		qid, err := SelectNsID("quid")
		if err != nil {
			return "", err
		}
		isAdmin, err := IsUserInAdminGroup(uid, qid)
		if err != nil {
			return "", err
		}
		if !isAdmin {
			return "", errors.New("the user is not namespace admin")
		}
	}

	// get the refresh token
	log.Encrypt("Gen token", ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, userName, []byte(ns.RefreshKey))
	token, err := tokens.GenRefreshToken(ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, userName, []byte(ns.RefreshKey))
	if err != nil {
		msg := "Error generating refresh token"
		log.Error(msg, err)
		return "", err
	}
	return token, nil
}
