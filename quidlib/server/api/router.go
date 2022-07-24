package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	color "github.com/acmacalister/skittles"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/quid/quidlib/conf"
)

// SessionsStore : the session cookies store.
var SessionsStore = sessions.NewCookieStore([]byte(conf.EncodingKey))

// AdminNsKey : store the Quid namespace key for admin
var AdminNsKey = []byte("")

// RunServer : configure and run the server.
func RunServer(adminNsKey, address string) {
	AdminNsKey = []byte(adminNsKey)

	//TODO echoServer.Use(middleware.Logger())
	//TODO
	//TODO if !conf.IsDevMode {
	//TODO 	echoServer.Use(middleware.Recover())
	//TODO 	echoServer.Use(middleware.Secure())
	//TODO }
	//TODO
	//TODO echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//TODO 	AllowOrigins:     []string{"*"},
	//TODO 	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
	//TODO 	AllowMethods:     []string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete},
	//TODO 	AllowCredentials: true,
	//TODO }))

	g, err := garcon.New(
		garcon.WithURLs(address),
		garcon.WithDocURL("/doc"),
		garcon.WithServerHeader("Quid"),
		garcon.WithJWT(adminNsKey, 0, true),
		garcon.WithLimiter(20, 30),
		garcon.WithProm(9193, address),
		garcon.WithDev(conf.IsDevMode))
	if err != nil {
		log.Panic("garcon.New: ", err)
	}

	r := chi.NewRouter()

	//TODO echoServer.Use(session.MiddlewareWithConfig(session.Config{Store: SessionsStore}))

	// serve static files
	ws := garcon.NewStaticWebServer("ui/dist", g.Writer)
	r.NotFound(ws.ServeFile("index.html", "text/html; charset=utf-8")) // catches index.html and other Vue sub-folders
	r.Get("/favicon.ico", ws.ServeFile("favicon.ico", "image/x-icon"))
	r.Get("/assets/*", ws.ServeAssets())
	r.Get("/version", g.ServeVersion())

	// HTTP Routes
	// public routes
	r.Post("/token/refresh/{timeout}", RequestRefreshToken)
	r.Post("/token/access/{timeout}", RequestAccessToken)
	r.Post("/admin_login", AdminLogin)
	r.Post("/admin_token/access/", RequestAdminAccessToken)

	// admin routes
	r.Route("/admin", func(r chi.Router) {
		//TODO a := echoServer.Group("/admin")
		//TODO config := middleware.JWTConfig{
		//TODO 	Claims:     &tokens.AdminAccessClaim{},
		//TODO 	SigningKey: []byte(adminNsKey),
		//TODO }
		//TODO a.Use(middleware.JWTWithConfig(config))
		//TODO a.Use(AdminMiddleware)

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
			//TODO nsadm.Use(middleware.JWTWithConfig(config))
			//TODO nsadm.Use(NsAdminMiddleware)

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

	server := http.Server{
		Addr:              address,
		Handler:           g.Middlewares.Then(r),
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Minute, // Garcon.Limiter postpones response, attacker should wait long time.
		IdleTimeout:       time.Second,
		ConnState:         g.ConnState,
		ErrorLog:          log.Default(),
	}

	log.Print("Server listening on http://localhost", server.Addr)
	log.Fatal(server.ListenAndServe())
}
