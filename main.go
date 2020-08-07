package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	. "github.com/logrusorgru/aurora"

	"github.com/synw/quid/quidlib/conf"
	"github.com/synw/quid/quidlib/server/api"
	"github.com/synw/quid/quidlib/server/db"
	"github.com/synw/quid/quidlib/tokens"
)

// SessionsStore : the session cookies store
var SessionsStore = sessions.NewCookieStore([]byte(conf.EncodingKey))

func main() {
	init := flag.Bool("init", false, "initialize and create a superuser")
	key := flag.Bool("key", false, "create a random key")
	isDevMode := flag.Bool("dev", false, "development mode")
	genConf := flag.Bool("conf", false, "generate a config file")
	//heroku := flag.Bool("heroku", false, "Run the app on Heroku")
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

	// init conf flag
	found, err := conf.Init(*isDevMode)
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		fmt.Println("No config file found. Use the -conf option to generate one")
		return
	}

	// db
	db.Init(*isDevMode)
	err = db.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	db.ExecSchema()

	// initialization flag
	if *init {
		db.InitDbConf()
		return
	}

	api.Init(*isDevMode)
	tokens.Init(*isDevMode)

	// get the admin namespace
	_, adminNS, err := db.SelectNamespaceFromName("quid")
	if err != nil {
		log.Fatal(err)
	}

	// http server
	e := echo.New()

	e.Use(middleware.Logger())
	if !*isDevMode {
		e.Use(middleware.Recover())
		e.Use(middleware.Secure())
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:8080"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
			//AllowCredentials: true,
		}))
	}

	e.Use(session.MiddlewareWithConfig(session.Config{Store: SessionsStore}))

	// serve static files in production
	if !conf.IsDevMode {
		e.File("/", "quidui/dist/index.html")
		e.Static("/js", "quidui/dist/js")
		e.Static("/css", "quidui/dist/css")
	}

	// public routes
	e.POST("/token/refresh/:timeout", api.RequestRefreshToken)
	e.POST("/token/access/:timeout", api.RequestAccessToken)
	e.POST("/admin_login", api.AdminLogin)

	// admin routes
	a := e.Group("/admin")
	config := middleware.JWTConfig{
		Claims:     &tokens.StandardAccessClaims{},
		SigningKey: []byte(adminNS.Key),
	}
	a.Use(middleware.JWTWithConfig(config))
	a.Use(api.AdminMiddleware)
	a.GET("/logout", api.AdminLogout)
	g := a.Group("/groups")
	g.POST("/add", api.CreateGroup)
	g.POST("/delete", api.DeleteGroup)
	g.POST("/info", api.GroupsInfo)
	g.POST("/add_user", api.AddUserInGroup)
	g.POST("/remove_user", api.RemoveUserFromGroup)
	g.GET("/all", func(c echo.Context) error {
		data, err := db.SelectAllGroups()
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
	ns.POST("/max-ttl", api.SetNamespaceTokenMaxTTL)
	ns.POST("/max-refresh-ttl", api.SetNamespaceRefreshTokenMaxTTL)
	ns.POST("/groups", api.GroupsForNamespace)
	ns.POST("/endpoint", api.SetNamespaceEndpointAvailability)
	ns.GET("/all", func(c echo.Context) error {
		data, err := db.SelectAllNamespaces()
		if err != nil {
			log.Fatalln(err)
		}
		return c.JSON(http.StatusOK, &data)
	})

	if conf.IsDevMode {
		fmt.Println(Bold(Red("Running in development mode")))
	}

	// run server
	e.Logger.Fatal(e.Start(":8082"))
}
