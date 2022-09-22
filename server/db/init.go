package db

import (
	"errors"

	"github.com/manifoldco/promptui"

	"github.com/teal-finance/quid/server"
	"github.com/teal-finance/quid/tokens"
)

func CreateQuidAdmin(username, password string) error {
	exist, err := NamespaceExists("quid")
	if err != nil {
		return err
	}

	var nsID int64
	if exist {
		nsID, err = SelectNsID("quid")
	} else {
		log.V().Data(`Creating the "quid" namespace`)
		algo := "HS256"
		accessKey := tokens.GenerateKeyHMAC(256)
		refreshKey := tokens.GenerateKeyHMAC(256)
		nsID, err = CreateNamespace("quid", "6m", "24h", algo, accessKey, refreshKey, false)
	}
	if err != nil {
		return err
	}

	exist, err = GroupExists("quid_admin", nsID)
	if err != nil {
		return err
	}

	var gid int64
	if exist {
		var g server.Group
		g, err = SelectGroup("quid_admin", nsID)
		gid = g.ID
	} else {
		log.V().Data(`Creating the "quid_admin" group`)
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
		log.Dataf(`Quid Admin already created (%d user in "quid_admin" group)`, n)
		return nil
	}

	if username == "" {
		log.V().Input("Enter the Quid Admin username")
		username, err = promptForUsername()
		if err != nil {
			log.ParamError(err)
			return err
		}
	}

	if password == "" {
		log.V().Input("Enter the Quid Admin password")
		password, err = promptForPassword()
		if err != nil {
			log.ParamError(err)
			return err
		}
	}

	log.V().Dataf("Create the Quid Admin user=%q pwd=%d bytes", username, len(password))
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
