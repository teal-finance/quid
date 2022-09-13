package db

import (
	"errors"

	"github.com/manifoldco/promptui"

	"github.com/teal-finance/quid/server"
	"github.com/teal-finance/quid/tokens"
)

func CreateQuidAdmin(username, password string) error {
	log.Data("Initializing Quid database")

	found, err := NamespaceExists("quid")
	if err != nil {
		return err
	}

	var nsID int64
	if found {
		nsID, err = SelectNsID("quid")
	} else {
		log.Data(`Creating the "quid" namespace`)
		algo := "HS256"
		accessKey := tokens.GenerateKeyHMAC(256)
		refreshKey := tokens.GenerateKeyHMAC(256)
		nsID, err = CreateNamespace("quid", "6m", "24h", algo, accessKey, refreshKey, false)
	}
	if err != nil {
		return err
	}

	found, err = GroupExists("quid_admin", nsID)
	if err != nil {
		return err
	}

	var gid int64
	if found {
		var g server.Group
		g, err = SelectGroup("quid_admin", nsID)
		gid = g.ID
	} else {
		log.Data(`Creating the "quid_admin" group`)
		gid, err = CreateGroup("quid_admin", nsID)
	}
	if err != nil {
		return err
	}

	n, err := CountUsersInGroup(gid)
	if err != nil {
		return err
	}

	if n > 0 {
		log.Data(`There are already %d users in "quid_admin" group => Do not create the Quid Admin user`)
		return nil
	}

	if username == "" {
		username, err = promptForUsername()
		if err != nil {
			log.ParamError(err)
			return err
		}
	}

	if password == "" {
		password, err = promptForPassword()
		if err != nil {
			log.ParamError(err)
			return err
		}
	}

	log.Dataf("Create the Quid Admin user usr=%q pwdLen=%d", username, len(password))
	u, err := CreateUser(username, password, nsID)
	if err != nil {
		return err
	}

	return AddUserInGroup(u.ID, gid)
}

func promptForUsername() (string, error) {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Username",
		Validate: validate,
	}

	return prompt.Run()
}

func promptForPassword() (string, error) {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("Password must have more than 6 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	return prompt.Run()
}
