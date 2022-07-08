package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/teal-finance/quid/quidlib/cmds"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/api"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

func main() {
	init := flag.Bool("init", false, "initialize and create a superuser")
	key := flag.Bool("key", false, "create a random key")
	env := flag.Bool("env", false, "init from environment variables not config file")
	isDevMode := flag.Bool("dev", false, "development mode")
	isVerbose := flag.Bool("v", false, "verbose (info and debug logs)")
	genConf := flag.Bool("conf", false, "generate a config file")
	genDevToken := flag.Bool("devtoken", false, "generate a quid admin dev token for frontend")
	genDevNsToken := flag.Bool("devnstoken", false, "generate a namespace admin dev token for frontend")
	flag.Parse()

	// key flag
	if *key {
		if *env {
			fmt.Println("The key command is not allowed when initializing from environment variables")
			os.Exit(1)
		}

		fmt.Println(tokens.GenKey())
		return
	}

	// gen conf flag
	if *genConf {
		if *env {
			fmt.Println("This command is not allowed when initializing from environment variables")
			os.Exit(2)
		}

		cmds.GeNConf()
		return
	}

	// Read configuration
	var (
		conn, port string
		autoConfDb bool
	)
	if *env {
		// env flag
		conn, port = conf.InitFromEnv(*isDevMode)
		autoConfDb = (conf.AdminUser != "") && (conf.AdminPassword != "")
	} else {
		// init conf flag
		conn, port = conf.InitFromFile(*isDevMode)
	}
	isCmd := *genDevToken || *genDevNsToken

	// Database
	db.Init(*isVerbose, *isDevMode, isCmd)

	if err := db.Connect(conn); err != nil {
		log.Fatalln(err)
	}

	if err := db.ExecSchema(); err != nil {
		log.Fatalln(err)
	}

	// gen dev token flag
	if *genDevToken {
		if *env {
			fmt.Println("This command is not allowed when initializing from environment variables")
			os.Exit(2)
		}

		username := os.Args[2]
		err := cmds.WriteDevAdminToken(username)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Dev token generated in env file")
		return
	}

	// gen namespace dev token flag
	if *genDevNsToken {
		if *env {
			fmt.Println("This command is not allowed when initializing from environment variables")
			os.Exit(2)
		}

		username := os.Args[2]
		namespace := os.Args[3]
		err := cmds.WriteNsAdminToken(username, namespace)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Dev nsadmin token generated in env file for user", username, "and namespace", namespace)
		return
	}

	// flag -init => initialize database
	if *init {
		if *env {
			fmt.Println("The init command is not allowed when initializing from environment variables")
			os.Exit(5)
		}

		db.InitDbConf()
		return
	}

	if autoConfDb {
		fmt.Println("Configure automatically the DB")
		db.InitDbAutoConf(conf.AdminUser, conf.AdminPassword)
	}

	api.Init(*isVerbose, *isDevMode)
	tokens.Init(*isVerbose, *isDevMode, isCmd)

	// get the admin namespace
	_, adminNS, err := db.SelectNamespaceFromName("quid")
	if err != nil {
		log.Fatal(err)
	}

	// http server
	api.RunServer(adminNS.Key, ":"+port)
}
