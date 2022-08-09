package api

import (
	"fmt"
	"log"

	color "github.com/acmacalister/skittles"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/quidlib/conf"
	"github.com/teal-finance/quid/quidlib/tokens"
)

// AdminNsKey : store the Quid namespace key for admin
var AdminNsKey = []byte("")

var echoServer = echo.New()

// var Garcon *garcon.Garcon
var Incorruptible *incorruptible.Incorruptible

// RunServer : configure and run the server.
func RunServer(adminNsKey string, port int) {
	AdminNsKey = []byte(adminNsKey)

	// TODO FIXME
	if !conf.IsDevMode {
		echoServer.Use(middleware.Recover())
		echoServer.Use(middleware.Secure())
	}

	g := garcon.New(
		garcon.WithDocURL("/doc"),
		garcon.WithServerHeader("Quid"),
		garcon.WithIncorruptible(conf.EncodingKey, 3600*24, true),
		garcon.WithLimiter(20, 30),
		garcon.WithProm(9193, "Quid"),
		garcon.WithDev(conf.IsDevMode))

	session, ok := g.TokenChecker().(*incorruptible.Incorruptible)
	if !ok {
		emo.Error("Garcon.Checker is not Incorruptible")
		log.Panic("Garcon.Checker is not Incorruptible")
	}
	// Garcon = g
	Incorruptible = session

	r := chi.NewRouter()

	// serve static files
	ws := garcon.NewStaticWebServer("ui/dist", g.Writer)
	r.NotFound(ws.ServeFile("index.html", "text/html; charset=utf-8")) // catches index.html and other Vue sub-folders
	r.Get("/favicon.ico", ws.ServeFile("favicon.ico", "image/x-icon"))
	r.Get("/assets/*", ws.ServeAssets())
	r.Get("/version", garcon.ServeVersion())

	// HTTP Routes
	// public routes
	r.Post("/token/refresh/{timeout}", RequestRefreshToken)
	r.Post("/token/access/{timeout}", RequestAccessToken)
	r.Post("/admin_login", AdminLogin)
	r.Post("/admin_token/access/", RequestAdminAccessToken)

	// admin routes
	r.Route("/admin", func(r chi.Router) {
		a := echoServer.Group("/admin")
		config := middleware.JWTConfig{
			Claims:     &tokens.AdminAccessClaim{},
			SigningKey: []byte(adminNsKey),
		}
		a.Use(middleware.JWTWithConfig(config))
		// FIXME a.Use(AdminMiddleware)

		// HTTP API
		r.Get("/logout", AdminLogout)
		r.Route("/groups", func(r chi.Router) {
			r.Post("/add", CreateGroup)
			r.Post("/delete", DeleteGroup)
			r.Post("/info", GroupsInfo)
			r.Post("/add_user", AddUserInGroup)
			r.Post("/remove_user", RemoveUserFromGroup)
			// r.Get("/all", AllGroups) // TODO: remove when old frontend is disabled
			r.Post("/nsall", AllGroupsForNamespace)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/add", CreateUserHandler)
			r.Post("/delete", DeleteUser)
			r.Post("/groups", UserGroupsInfo)
			r.Post("/orgs", UserOrgsInfo)
			// r.Get("/all", AllUsers) // TODO: remove when old frontend is disabled
			r.Post("/nsall", AllUsersInNamespace)
			// r.Post("/search", SearchForUsersInNamespace)
		})

		r.Route("/namespaces", func(r chi.Router) {
			r.Post("/add", CreateNamespace)
			r.Post("/delete", DeleteNamespace)
			r.Post("/find", FindNamespace)
			r.Post("/info", NamespaceInfo)
			r.Post("/key", GetNamespaceKey)
			r.Post("/max-ttl", SetNamespaceTokenMaxTTL)
			r.Post("/max-refresh-ttl", SetNamespaceRefreshTokenMaxTTL)
			r.Post("/groups", GroupsForNamespace)
			r.Post("/endpoint", SetNamespaceEndpointAvailability)
			r.Get("/all", AllNamespaces)
		})

		r.Route("/orgs", func(r chi.Router) {
			r.Get("/all", AllOrgs)
			r.Post("/add", CreateOrg)
			r.Post("/delete", DeleteOrg)
			r.Post("/find", FindOrg)
			r.Post("/add_user", AddUserInOrg)
			r.Post("/remove_user", RemoveUserFromOrg)
		})

		r.Route("/nsadmin", func(r chi.Router) {
			r.Post("/add", CreateAdministrators)
			r.Post("/nsall", AllAdministratorsInNamespace)
			r.Post("/delete", DeleteAdministrator)
			r.Post("/search/nonadmins", SearchForNonAdminUsersInNamespace)
		})

		// Namespace admin endpoints
		r.Route("/ns", func(r chi.Router) {
			// TODO nsadm.Use(middleware.JWTWithConfig(config))
			// TODO nsadm.Use(NsAdminMiddleware)

			// nsadmin users
			r.Route("/users", func(r chi.Router) {
				r.Post("/add", CreateUserHandler)
				r.Post("/delete", DeleteUser)
				r.Post("/groups", UserGroupsInfo)
				r.Post("/nsall", AllUsersInNamespace)
			})

			// nsadmin groups
			r.Route("/groups", func(r chi.Router) {
				r.Post("/add", CreateGroup)
				r.Post("/delete", DeleteGroup)
				r.Post("/info", GroupsInfo)
				r.Post("/add_user", AddUserInGroup)
				r.Post("/remove_user", RemoveUserFromGroup)
				r.Post("/nsall", AllGroupsForNamespace)
			})
		})
	})

	if conf.IsDevMode {
		fmt.Println(color.BoldRed("Running in development mode"))
	}

	server := g.Server(r, port)

	log.Print("Server listening on http://localhost", server.Addr)
	log.Fatal(server.ListenAndServe())
}
