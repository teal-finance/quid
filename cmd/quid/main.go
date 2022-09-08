package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/teal-finance/emo"
	"github.com/teal-finance/quid/pkg/cmds"
	"github.com/teal-finance/quid/pkg/conf"
	"github.com/teal-finance/quid/pkg/server/api"
	"github.com/teal-finance/quid/pkg/server/db"
	"github.com/teal-finance/quid/pkg/tokens"
)

var log = emo.NewZone("main")

func main() {
	init := flag.Bool("init", false, "initialize and create the QuidAdmin")
	key := flag.String("key", "", "create a random key among HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512 and Ed25519")
	env := flag.Bool("env", false, "init from environment variables not config file")
	isDevMode := flag.Bool("dev", false, "development mode")
	isVerbose := flag.Bool("v", false, "verbose (info and debug logs)")
	genConf := flag.Bool("conf", false, "generate a config file")
	genDevToken := flag.Bool("devtoken", false, "generate a quid admin dev token for frontend")
	genDevNsToken := flag.Bool("devnstoken", false, "generate a namespace admin dev token for frontend")
	flag.Parse()

	// key flag
	if *key != "" {
		if *env {
			log.Fatal("The key command is not allowed when initializing from environment variables")
		}

		fmt.Println(tokens.GenerateSigningKey(*key))
		return
	}

	// gen conf flag
	if *genConf {
		log.Info("Generating config file")
		if *env {
			log.Fatal("This command is not allowed when initializing from environment variables")
		}
		if err := conf.Create(); err != nil {
			log.Fatal("Cannot create config file", err)
		}
		log.State("Config file created: edit config.json to provide your database settings")
		return
	}

	// Read configuration
	var (
		conn       string
		port       int
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
		log.Fatal(err)
	}

	if err := db.ExecSchema(); err != nil {
		log.Fatal(err)
	}

	// gen dev token flag
	if *genDevToken {
		if *env {
			log.Fatal("This command is not allowed when initializing from environment variables")
		}

		username := os.Args[2]
		err := cmds.WriteQuidAdminToken(username)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Dev token generated in env file")
		return
	}

	// gen namespace dev token flag
	if *genDevNsToken {
		if *env {
			log.Fatal("This command is not allowed when initializing from environment variables")
		}

		username := os.Args[2]
		namespace := os.Args[3]
		err := cmds.WriteNsAdminToken(username, namespace)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Dev nsadmin token generated in env file for user", username, "and namespace", namespace)
		return
	}

	// flag -init => initialize database
	if *init {
		if *env {
			log.Fatal("The init command is not allowed when initializing from environment variables")
		}

		db.InitDbConf()
		return
	}

	if autoConfDb {
		log.Info("Configure automatically the DB")
		db.InitDbAutoConf(conf.AdminUser, conf.AdminPassword)
	}

	printOnlyErrors := !*isVerbose && !*isDevMode
	if printOnlyErrors {
		emo.GlobalVerbosity(false)
	}

	api.Init(*isVerbose, *isDevMode)
	tokens.Init(*isVerbose, *isDevMode, isCmd)

	// http server
	api.RunServer(port)
}
