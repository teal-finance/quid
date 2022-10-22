package db

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/teal-finance/quid/server"
	"github.com/teal-finance/quid/tokens"
)

func CreateQuidAdmin(username, password string, forcePrompt bool) error {
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

	if forcePrompt || (username == "") {
		username, err = promptForUsername(username)
		if err != nil {
			return err
		}
	}

	if forcePrompt || (password == "") {
		password, err = promptForPassword(password, forcePrompt)
		if err != nil {
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

var (
	UsernameTooShort = errors.New("username must have more than 2 characters.")
	PasswordTooShort = errors.New("password must have more than 5 characters.")
)

func promptForUsername(username string) (string, error) {
	fmt.Println(`
Enter the Quid Admin username.
` + UsernameTooShort.Error())

	prompt := promptui.Prompt{
		Label:   "Username",
		Default: username,
		Validate: func(s string) error {
			if len(s) <= 2 {
				return UsernameTooShort
			}
			return nil
		},
	}

	username, err := prompt.Run()
	if err != nil {
		log.ParamError(err)
	}
	return username, err
}

func promptForPassword(password string, showPassword bool) (string, error) {
	fmt.Println(`
Enter the Quid Admin password.
` + PasswordTooShort.Error())

	mask := '*'
	if showPassword {
		mask = 0
	}

	prompt := promptui.Prompt{
		Label:   "Password",
		Default: password,
		Mask:    mask,
		Validate: func(input string) error {
			if len(input) <= 5 {
				return errors.New("password must have more than 6 characters")
			}
			return nil
		},
	}

	password, err := prompt.Run()
	if err != nil {
		log.ParamError(err)
	}
	return password, err
}
