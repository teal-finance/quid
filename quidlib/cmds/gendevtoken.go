package cmds

import (
	"errors"
	"log"
	"os"

	emolib "github.com/teal-finance/quid/quidlib/emo"
	"github.com/teal-finance/quid/quidlib/server/api"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

var out = emolib.NewZone("gentoken")

func GenDevAdminToken(username string) error {
	// get the namespace
	_, ns, err := db.SelectNamespaceFromName("quid")
	if err != nil {
		return err
	}
	_, u, err := db.SelectNonDisabledUser(username, ns.ID)
	if err != nil {
		return err
	}
	// check the user admin group
	isAdmin, err := api.IsUserInAdminGroup(u.ID, ns.ID)
	if err != nil {
		return err
	}
	if !isAdmin {
		log.Fatal("User is not admin")
	}
	// get the refresh token
	token, err := tokens.GenRefreshToken("24h", ns.MaxRefreshTokenTTL, ns.Name, u.UserName, []byte(ns.RefreshKey))
	if err != nil {
		msg := "Error generating refresh token"
		emo.Error(msg, err)
		return err
	}
	if token == "" {
		out.Info("Unauthorized: timeout max (", ns.MaxRefreshTokenTTL, ") for refresh token for namespace", ns.Name)
		return errors.New("")
	}

	// write the ui env file

	dir, err := os.Getwd()
	if err != nil {
		emo.Error(err)
		return err
	}
	filepath := dir + "/ui/.env.dev.local"
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	v := "VITE_DEV_TOKEN=\"" + token + "\""
	_, err = f.Write([]byte(v))
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	return nil
}
