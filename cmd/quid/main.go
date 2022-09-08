package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/teal-finance/emo"
	"github.com/teal-finance/quid/server/api"
	"github.com/teal-finance/quid/server/db"
	"github.com/teal-finance/quid/tokens"
)

var log = emo.NewZone("quid")

func main() {
	init := flag.Bool("init", false, "initialize and create the QuidAdmin")
	key := flag.String("key", "", "create a random key among HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512 and Ed25519")
	env := flag.Bool("env", false, "init from environment variables not config file")
	isDevMode := flag.Bool("dev", false, "development mode")
	isVerbose := flag.Bool("v", false, "verbose (info and debug logs)")
	genConf := flag.Bool("conf", false, "generate a config file")
	genDevQuidToken := flag.Bool("dev-quid-token", false, "generate a QuidAdmin JWT (was required to test the frontend)")
	genDevNsToken := flag.Bool("dev-ns-token", false, "generate a NamespaceAdmin JWT (was required to test the frontend)")
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
		if err := create(); err != nil {
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
		conn, port = initFromEnv(*isDevMode)
		autoConfDb = (adminUser != "") && (adminPassword != "")
	} else {
		// init conf flag
		conn, port = initFromFile(*isDevMode)
	}
	isCmd := *genDevQuidToken || *genDevNsToken

	// Database
	db.Init(*isVerbose, *isDevMode, isCmd)

	if err := db.Connect(conn); err != nil {
		log.Fatal(err)
	}

	if err := db.ExecSchema(); err != nil {
		log.Fatal(err)
	}

	// gen dev token flag
	if *genDevQuidToken {
		if *env {
			log.Fatal("This command is not allowed when initializing from environment variables")
		}

		username := os.Args[2]
		err := writeQuidAdminToken(username)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("QuidAdmin JWT generated in env file")
		return
	}

	// gen namespace dev token flag
	if *genDevNsToken {
		if *env {
			log.Fatal("This command is not allowed when initializing from environment variables")
		}

		username := os.Args[2]
		namespace := os.Args[3]
		err := writeNsAdminToken(username, namespace)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("NamespaceAdmin JWT generated in env file for user", username, "and namespace", namespace)
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
		db.InitDbAutoConf(adminUser, adminPassword)
	}

	printOnlyErrors := !*isVerbose && !*isDevMode
	if printOnlyErrors {
		emo.GlobalVerbosity(false)
	}

	api.Init(*isVerbose, *isDevMode)
	tokens.Init(*isVerbose, *isDevMode, isCmd)

	// http server
	api.RunServer(port, *isDevMode)
}
