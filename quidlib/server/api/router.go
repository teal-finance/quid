package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	color "github.com/logrusorgru/aurora"

	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/tokens"
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
	}
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	echoServer.Use(session.MiddlewareWithConfig(session.Config{Store: SessionsStore}))

	// serve static files
	echoServer.File("/", "adminui/dist/index.html")
	echoServer.File("/favicon.ico", "adminui/dist/favicon.ico")
	echoServer.Static("/assets", "adminui/dist/assets")

	// HTTP Routes
	// public routes
	echoServer.POST("/token/refresh/:timeout", RequestRefreshToken)
	echoServer.POST("/token/access/:timeout", RequestAccessToken)
	echoServer.POST("/admin_login", AdminLogin)
	echoServer.POST("/admin_token/access/", RequestAdminAccessToken)

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
	g.GET("/all", AllGroups) // TODO: remove when old frontend is disabled
	g.POST("/nsall", AllGroupsForNamespace)

	m := a.Group("/users")
	m.POST("/add", CreateUserHandler)
	m.POST("/delete", DeleteUser)
	m.POST("/groups", UserGroupsInfo)
	m.POST("/orgs", UserOrgsInfo)
	m.GET("/all", AllUsers) // TODO: remove when old frontend is disabled
	m.POST("/nsall", AllUsersInNamespace)
	m.POST("/search", SearchForUsersInNamespace)

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

	org := a.Group("/orgs")
	org.GET("/all", AllOrgs)
	org.POST("/add", CreateOrg)
	org.POST("/delete", DeleteOrg)
	org.POST("/find", FindOrg)
	org.POST("/add_user", AddUserInOrg)
	org.POST("/remove_user", RemoveUserFromOrg)

	nsa := a.Group("/nsadmin")
	nsa.POST("/add", CreateAdministrators)
	nsa.POST("/nsall", AllAdministratorsInNamespace)
	nsa.POST("/delete", DeleteAdministrator)

	if conf.IsDevMode {
		fmt.Println(color.Bold(color.Red("Running in development mode")))
	}

	echoServer.Logger.Fatal(echoServer.Start(":" + conf.Port))
}
