package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"

	"github.com/teal-finance/quid/quidlib/tokens"
)

// InitDbConf : initialize the database content.
func InitDbConf() {
	initDbConf(true, "", "")
}

// InitDbAutoConf : initialize the database content.
func InitDbAutoConf(username, password string) {
	initDbConf(false, username, password)
}

func initDbConf(prompt bool, username, password string) {
	if prompt {
		fmt.Println("Initializing Quid database")
	}

	// check namespace
	var nsID int64

	nsExists, err := NamespaceExists("quid")
	if err != nil {
		log.Fatal(err)
	}

	if nsExists {
		nsID, err = SelectNamespaceID("quid")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Creating the quid namespace")

		algo := "HS256"
		accessKey := tokens.GenerateKeyHMAC(256)
		refreshKey := tokens.GenerateKeyHMAC(256)

		nsID, err = CreateNamespace("quid", "6m", "24h", algo, accessKey, refreshKey, false)
		if err != nil {
			log.Fatal(err)
		}
	}

	// check base admin group
	var gid int64

	exists, err := GroupExists("quid_admin", nsID)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		group, er := SelectGroup("quid_admin", nsID)
		if er != nil {
			log.Fatal(er)
		}
		gid = group.ID
	} else {
		fmt.Println("Creating the quid admin group")
		gid, err = CreateGroup("quid_admin", nsID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// check superuser
	if n, _ := CountUsersInGroup(gid); n == 0 {
		var name string
		if prompt {
			fmt.Println("Create a superuser")
			name, err = promptForUsername()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			name = username
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

		u, err := CreateUser(name, pwd, nsID)
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
