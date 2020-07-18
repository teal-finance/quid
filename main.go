package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/synw/quid/quidlib"
	"github.com/synw/quid/quidlib/api"
	"github.com/synw/quid/quidlib/conf"
	"github.com/synw/quid/quidlib/db"
	"github.com/synw/quid/quidlib/tokens"
)

func main() {
	init := flag.Bool("init", false, "initialize and create a superuser")
	key := flag.Bool("key", false, "create a random key")
	isDevMode := flag.Bool("dev", false, "development mode")
	genConf := flag.Bool("conf", false, "generate a config file")
	flag.Parse()

	// key flag
	if *key {
		k := tokens.GenKey()
		fmt.Println(k)
		return
	}

	// gen conf flag
	if *genConf {
		fmt.Println("Generating config file")
		conf.Create()
		fmt.Println("Config file created: edit config.json to provide your database settings")
		return
	}

	found, err := conf.Init()
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		fmt.Println("No config file found. Use the -conf option to generate one")
		return
	}

	// db
	err = db.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	db.ExecSchema()

	// initialization flag
	if *init {
		quidlib.InitDbConf()
		return
	}

	// http server
	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	// public routes
	e.POST("/request_token", api.RequestToken)
	e.POST("/admin_login", api.AdminLogin)

	config := middleware.JWTConfig{
		Claims:     &tokens.StandardUserClaims{},
		SigningKey: []byte(conf.EncodingKey),
	}

	// admin routes
	a := e.Group("/admin")
	if !*isDevMode {
		a.Use(middleware.JWTWithConfig(config))
		a.Use(api.AdminMiddleware)
	}

	g := a.Group("/groups")
	g.POST("/add", api.CreateGroup)
	g.POST("/delete", api.DeleteGroup)
	g.POST("/info", api.GroupsInfo)
	g.GET("/all", func(c echo.Context) error {
		data, err := db.SelectGroups()
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(http.StatusOK, &data)
	})

	m := a.Group("/users")
	m.POST("/add", api.CreateUserHandler)
	m.POST("/delete", api.DeleteUser)
	m.POST("/info", api.UserInfo)
	m.GET("/all", func(c echo.Context) error {
		data, err := db.SelectAllUsers()
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(http.StatusOK, &data)
	})

	ns := a.Group("/namespaces")
	ns.POST("/add", api.CreateNamespace)
	ns.POST("/delete", api.DeleteNamespace)
	ns.POST("/find", api.FindNamespace)
	ns.POST("/info", api.NamespaceInfo)
	ns.POST("/key", api.GetNamespaceKey)
	ns.POST("/endpoint", api.SetNamespaceEndpointAvailability)
	ns.GET("/all", func(c echo.Context) error {
		data, err := db.SelectAllNamespaces()
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(http.StatusOK, &data)
	})

	if *isDevMode {
		fmt.Println("Running in development mode")
	}

	// run server
	e.Logger.Fatal(e.Start(":8082"))
}
