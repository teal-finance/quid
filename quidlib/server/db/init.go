package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// InitDbConf : initialize the database content
func InitDbConf() {
	initDbConf(true, "", "")
}

// InitDbAutoConf : initialize the database content
func InitDbAutoConf(username, password string) {
	initDbConf(false, username, password)
}

func initDbConf(prompt bool, username, password string) {
	if prompt {
		fmt.Println("Initializing Quid database")
	}
	// check namespace
	nsexists, err := NamespaceExists("quid")
	if err != nil {
		log.Fatal(err)
	}
	var nsid int64
	if !nsexists {
		key := tokens.GenKey()
		refreshKey := tokens.GenKey()
		fmt.Println("Creating the quid namespace")
		nsid, err = CreateNamespace("quid", key, refreshKey, "6m", "24h", false)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		nsid, err = SelectNamespaceID("quid")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check base admin group
	exists, err := GroupExists("quid_admin", nsid)
	if err != nil {
		log.Fatal(err)
	}
	var gid int64
	if !exists {
		fmt.Println("Creating the quid admin group")
		gid, err = CreateGroup("quid_admin", nsid)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		group, err := SelectGroup("quid_admin", nsid)
		if err != nil {
			log.Fatal(err)
		}
		gid = group.ID
	}
	// check superuser
	n, _ := CountUsersInGroup(gid)
	if n == 0 {
		var uname string
		if prompt {
			fmt.Println("Create a superuser")
			uname, err = promptForUsername()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			uname = username
		}
		var pwd string
		if prompt {
			pwd, err = promptForPassword()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			pwd = password
		}
		u, err := CreateUser(uname, pwd, nsid)
		if err != nil {
			log.Fatal(err)
		}
		err = AddUserInGroup(u.ID, gid)
		if err != nil {
			log.Fatal(err)
		}
		if prompt {
			fmt.Println("Superuser", username, "created")
		}
	}
	if prompt {
		fmt.Println("Initialization complete")
	}
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
