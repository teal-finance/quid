package main

import (
	"flag"
	"fmt"
	"log"

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
	isVerbose := flag.Bool("v", false, "verbose mode")
	genConf := flag.Bool("conf", false, "generate a config file")
	flag.Parse()

	// key flag
	if *key {
		if *env {
			fmt.Println("The key command is not allowed when initializing from environment variables")
			log.Fatal()
		}
		k := tokens.GenKey()
		fmt.Println(k)
		return
	}

	// gen conf flag
	if *genConf {
		if *env {
			fmt.Println("The conf command is not allowed when initializing from environment variables")
			log.Fatal()
		}
		fmt.Println("Generating config file")
		conf.Create()
		fmt.Println("Config file created: edit config.json to provide your database settings")
		return
	}

	autoConfDb := false
	// env flag
	if *env {
		if *isVerbose {
			fmt.Println("Initializing from env")
		}
		autoConfDb = conf.InitFromEnv(*isDevMode)
	} else {
		// init conf flag
		found, err := conf.Init(*isDevMode)
		if err != nil {
			log.Fatal(err)
		}
		if !found {
			fmt.Println("No config file found. Use the -conf option to generate one")
			return
		}
	}

	// db
	db.Init(*isVerbose)
	if err := db.Connect(); err != nil {
		log.Fatalln(err)
	}
	db.ExecSchema()

	// initialization flag
	if *init {
		if *env {
			fmt.Println("The init command is not allowed when initializing from environment variables")
			log.Fatal()
		}
		db.InitDbConf()
		return
	}
	if autoConfDb {
		if *isVerbose {
			fmt.Println("Running autoconf")
		}
		db.InitDbAutoConf(conf.DefaultAdminUser, conf.DefaultAdminPassword)
	}

	api.Init(*isVerbose)
	tokens.Init(*isVerbose)

	// get the admin namespace
	_, adminNS, err := db.SelectNamespaceFromName("quid")
	if err != nil {
		log.Fatal(err)
	}

	// http server
	api.RunServer(adminNS.Key)
}
