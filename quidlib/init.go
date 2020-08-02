package quidlib

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/synw/quid/quidlib/api"
	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/tokens"
)

// InitDbConf : initialize the database content
func InitDbConf() {
	fmt.Println("Initializing Quid database")
	// check namespace
	nsexists, err := db.NamespaceExists("quid")
	if err != nil {
		log.Fatal(err)
	}
	var nsid int64
	if !nsexists {
		key := tokens.GenKey()
		refreshKey := tokens.GenKey()
		fmt.Println("Creating the quid namespace")
		nsid, err = db.CreateNamespace("quid", key, refreshKey, "20m", "24h", false)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		nsid, err = db.SelectNamespaceID("quid")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check base admin group
	exists, err := db.GroupExists("quid_admin", nsid)
	if err != nil {
		log.Fatal(err)
	}
	var gid int64
	if !exists {
		fmt.Println("Creating the quid admin group")
		gid, err = db.CreateGroup("quid_admin", nsid)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		group, err := db.SelectGroup("quid_admin", nsid)
		if err != nil {
			log.Fatal(err)
		}
		gid = group.ID
	}
	// check superuser
	n, err := db.CountUsersInGroup(gid)
	if n == 0 {
		fmt.Println("Create a superuser")
		username, err := promptForUsername()
		if err != nil {
			log.Fatal(err)
		}
		pwd, err := promptForPassword()
		if err != nil {
			log.Fatal(err)
		}
		u, err := api.CreateUser(username, pwd, nsid)
		if err != nil {
			log.Fatal(err)
		}
		err = db.AddUserInGroup(u.ID, gid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Superuser", username, "created")
	}
	fmt.Println("Initialization complete")
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
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return result, nil
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
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
