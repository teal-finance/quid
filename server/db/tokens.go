package db

import (
	"github.com/teal-finance/quid/tokens"
)

// GenNsAdminToken : generate a refresh token for an admin user and namespace
// Deprecated because this function is not used.
func GenNsAdminToken(userName, nsName string) (string, error) {
	log.Info("Generating NS Admin token for", userName, nsName)

	// get the namespace
	ns, err := SelectNsFromName(nsName)
	if err != nil {
		return "", err
	}

	uid, err := selectEnabledUsrID(userName)
	if err != nil {
		return "", err
	}

	// check admin perms
	adminType, err := GetUserType(nsName, ns.ID, uid)
	if err != nil {
		return "", err
	}
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
			return "", log.Warn("the user is not a namespace admin").Err()
		}
	}

	log.Encrypt("Gen token", ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, userName, []byte(ns.RefreshKey))

	token, err := tokens.GenRefreshToken(ns.MaxRefreshTokenTTL, ns.MaxRefreshTokenTTL, ns.Name, userName, []byte(ns.RefreshKey))
	if err != nil {
		return "", log.Error("Error generating refresh token", err).Err()
	}
	return token, nil
}
