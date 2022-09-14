package main

import (
	"flag"

	"github.com/teal-finance/emo"
	"github.com/teal-finance/garcon/gg"
	"github.com/teal-finance/quid/crypt"
	"github.com/teal-finance/quid/server/api"
	"github.com/teal-finance/quid/server/db"
	"github.com/teal-finance/quid/tokens"
)

var log = emo.NewZone("quid")

const (
	// defaultKey is AES-128-bits (16 bytes) in hexadecimal form (32 digits).
	defaultKey  = "00112233445566778899aabbccddeeff"
	defaultUsr  = "admin"
	defaultPwd  = "my_password"
	defaultPort = 8090

	defaultDBUser = "pguser"
	defaultDBPass = "my_password"
	defaultDBHost = "localhost"
	defaultDBPort = "5432"
	defaultDBName = "quid"
)

var defaultDBurl = buildURL(defaultDBUser, defaultDBPass, defaultDBHost, defaultDBPort, defaultDBName)

func buildURL(usr, pwd, host, port, name string) string {
	return "postgres://" + usr + ":" + pwd + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
}

func main() {
	var (
		dev     = flag.Bool("dev", false, "Development mode")
		verbose = flag.Bool("v", false, "Verbose (enables the info and debug logs)")
		key     = flag.String("key", gg.EnvStr("QUID_KEY", defaultKey), "AES-128 key to encrypt the private keys of the refresh/access tokens in the database. Accept 32 hexadecimal digits or 22 Base64 characters. Env. var: QUID_KEY")
		admin   = flag.String("admin", gg.EnvStr("QUID_ADMIN_USR", defaultUsr), "The username of the Quid Administrator. Env. var: QUID_ADMIN_USR")
		pwd     = flag.String("pwd", gg.EnvStr("QUID_ADMIN_PWD", defaultPwd), "The password of the Quid Administrator. Env. var: QUID_ADMIN_PWD")
		conf    = flag.Bool("conf", false, `Generate a "config.json" file with a random AES-128 key`)
		dbUser  = flag.String("db-usr", gg.EnvStr("DB_USR", defaultDBUser), "Username to read/write the database. Env. var: DB_USR")
		dbPass  = flag.String("db-pwd", gg.EnvStr("DB_PWD", defaultDBPass), "Password of the database user. Env. var: DB_PWD")
		dbHost  = flag.String("db-host", gg.EnvStr("DB_HOST", defaultDBHost), "Network location of the Postgres server. Env. var: DB_HOST")
		dbPort  = flag.String("db-port", gg.EnvStr("DB_PORT", defaultDBPort), "TCP port of the Postgres server. Env. var: DB_PORT")
		dbName  = flag.String("db-name", gg.EnvStr("DB_NAME", defaultDBName), "Name of the Postgres database. Env. var: DB_NAME")
		dbURL   = flag.String("db-url", gg.EnvStr("DB_URL", defaultDBurl), "The endpoint of the PostgreSQL server. Env. var: DB_URL")
		www     = flag.String("www", gg.EnvStr("WWW_DIR", "ui/dist"), "Folder of the web static files. Env. var: WWW_DIR")
		port    = flag.Int("port", gg.EnvInt("PORT", defaultPort), "Listening port of the Quid server")
	)
	flag.Parse()

	cfgName, cfgUsr, cfgPwd, cfgKey := readConfigFile()
	if cfgName != "" && *dbName == defaultDBName {
		*dbName = cfgName
	}
	if cfgUsr != "" && *dbUser == defaultDBUser {
		*dbUser = cfgUsr
	}
	if cfgPwd != "" && *dbPass == defaultDBPass {
		*dbPass = cfgPwd
	}
	if cfgKey != "" && *key == defaultKey {
		*key = cfgKey
	}

	if *dbURL == defaultDBurl {
		*dbURL = buildURL(*dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	}

	if (!*dev) && (*key == defaultKey) && (*dbURL == defaultDBurl) {
		if !*verbose {
			log.Print("Default values for -dev, -key QUID_KEY and -db DB_URL => Enable -v verbose mode")
			*verbose = true
		}
		log.Print("Default values for -key QUID_KEY and -db DB_URL => Enable -dev mode")
		*dev = true
	}
	emo.GlobalVerbosity(*verbose)

	log.Param("-dev                  =", *dev)
	log.Param("-v                    =", *verbose)
	log.Param("-conf                 =", *conf)
	log.Param("-key   QUID_KEY       =", *key)
	log.Param("-admin QUID_ADMIN_USR =", *admin)
	log.Param("-pwd   QUID_ADMIN_PWD =", *pwd)
	log.Param("-conf                 =", *conf)
	log.Param("-db-usr  DB_USR       =", *dbUser)
	log.Param("-db-pwd  DB_PWD       =", *dbPass)
	log.Param("-db-host DB_HOST      =", *dbHost)
	log.Param("-db-port DB_PORT      =", *dbPort)
	log.Param("-db-name DB_NAME      =", *dbName)
	log.Param("-db-url  DB_URL =", *dbURL)
	log.Param("-www     WWW_DIR      =", *www)
	log.Param("-port    PORT         =", *port)

	if *conf {
		log.Info(`Generating "config.json" file with random AES-128 key`)
		if err := createConfigFile(*dbName, *dbUser, *dbPass); err != nil {
			log.Fatal(`Cannot create "config.json" file`, err)
		}
		log.State(`Config file created: edit "config.json" to provide your database settings`)
		return
	}

	crypt.EncodingKey = tokens.DecodeHexOrB64(*key, 16)
	if crypt.EncodingKey == nil {
		log.Panic("Want AES-128 key in hexadecimal (32 digits) or Base64 (unpadded 22 characters RFC 4648 §5), but got", len(*key), "bytes:", *key)
	}

	if err := db.Connect(*dbURL); err != nil {
		log.Fatal(err)
	}

	if err := db.ExecSchema(); err != nil {
		log.Fatal(err)
	}

	err := db.CreateQuidAdmin(*admin, *pwd)
	if err != nil {
		log.Fatal(err)
	}

	api.RunServer(*port, *dev, *www)
}
