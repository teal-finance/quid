package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/server/api"
	"github.com/teal-finance/quid/quidlib/server/db"
	"github.com/teal-finance/quid/quidlib/tokens"
)

func main() {
	init := flag.Bool("init", false, "initialize and create a superuser")
	key := flag.Bool("key", false, "create a random key")
	env := flag.Bool("env", false, "init from environment variables not config file")
	isVerbose := flag.Bool("v", false, "verbose mode")
	genConf := flag.Bool("conf", false, "generate a config file")
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
			fmt.Println("The conf command is not allowed when initializing from environment variables")
			os.Exit(2)
		}

		fmt.Println("Generating config file")
		if err := conf.Create(); err != nil {
			fmt.Println("Cannot create config file ", err)
			os.Exit(3)
		}

		fmt.Println("Config file created: edit config.json to provide your database settings")

		return
	}

	// Read configuration
	var (
		conn, port string
		autoConfDb bool
	)
	if *env {
		// env flag
		conn, port = conf.InitFromEnv()
		autoConfDb = (conf.AdminUser != "") && (conf.AdminPassword != "")
	} else {
		// init conf flag
		conn, port = conf.InitFromFile()
	}

	// Database
	db.Init(*isVerbose)

	if err := db.Connect(conn); err != nil {
		log.Fatalln(err)
	}

	if err := db.ExecSchema(); err != nil {
		log.Fatalln(err)
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

	api.Init(*isVerbose)
	tokens.Init(*isVerbose)

	// get the admin namespace
	_, adminNS, err := db.SelectNamespaceFromName("quid")
	if err != nil {
		log.Fatal(err)
	}

	// http server
	api.RunServer(adminNS.Key, ":"+port)
}
