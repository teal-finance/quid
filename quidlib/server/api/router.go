package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	color "github.com/logrusorgru/aurora"

	"github.com/synw/quid/quidlib/conf"
	"github.com/synw/quid/quidlib/tokens"
)

// SessionsStore : the session cookies store
var SessionsStore = sessions.NewCookieStore([]byte(conf.EncodingKey))

var echoServer = echo.New()

// RunServer : configure and run the server
func RunServer(adminNsKey string) {

	echoServer.Use(middleware.Logger())
	if !conf.IsDevMode {
		echoServer.Use(middleware.Recover())
		echoServer.Use(middleware.Secure())
	} else {
		echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:8080"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
			//AllowCredentials: true,
		}))
	}

	echoServer.Use(session.MiddlewareWithConfig(session.Config{Store: SessionsStore}))

	// serve static files in production
	if !conf.IsDevMode {
		echoServer.File("/", "quidui/dist/index.html")
		echoServer.Static("/js", "quidui/dist/js")
		echoServer.Static("/css", "quidui/dist/css")
	}

	// SetHTTPRoutes :
	// public routes
	echoServer.POST("/token/refresh/:timeout", RequestRefreshToken)
	echoServer.POST("/token/access/:timeout", RequestAccessToken)
	echoServer.POST("/admin_login", AdminLogin)

	// admin routes
	a := echoServer.Group("/admin")
	config := middleware.JWTConfig{
		Claims:     &tokens.StandardAccessClaims{},
		SigningKey: []byte(adminNsKey),
	}
	a.Use(middleware.JWTWithConfig(config))
	a.Use(AdminMiddleware)
	a.GET("/logout", AdminLogout)

	g := a.Group("/groups")
	g.POST("/add", CreateGroup)
	g.POST("/delete", DeleteGroup)
	g.POST("/info", GroupsInfo)
	g.POST("/add_user", AddUserInGroup)
	g.POST("/remove_user", RemoveUserFromGroup)
	g.GET("/all", AllGroups)

	m := a.Group("/users")
	m.POST("/add", CreateUserHandler)
	m.POST("/delete", DeleteUser)
	m.POST("/info", UserInfo)
	m.GET("/all", AllUsers)

	ns := a.Group("/namespaces")
	ns.POST("/add", CreateNamespace)
	ns.POST("/delete", DeleteNamespace)
	ns.POST("/find", FindNamespace)
	ns.POST("/info", NamespaceInfo)
	ns.POST("/key", GetNamespaceKey)
	ns.POST("/max-ttl", SetNamespaceTokenMaxTTL)
	ns.POST("/max-refresh-ttl", SetNamespaceRefreshTokenMaxTTL)
	ns.POST("/groups", GroupsForNamespace)
	ns.POST("/endpoint", SetNamespaceEndpointAvailability)
	ns.GET("/all", AllNamespaces)

	if conf.IsDevMode {
		fmt.Println(color.Bold(color.Red("Running in development mode")))
	}

	echoServer.Logger.Fatal(echoServer.Start(":8082"))
}
