package api

import (
	"log"
	"net/http"

	color "github.com/acmacalister/skittles"
	"github.com/go-chi/chi/v5"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/quidlib/conf"
)

var Incorruptible *incorruptible.Incorruptible

var gw garcon.Writer

// RunServer : configure and run the server.
func RunServer(port int) {
	server := newServer(port)

	if conf.IsDevMode {
		log.Print("INF " + color.BoldRed("Running in development mode"))
	}

	log.Print("INF Server listening on " + color.UnderlineBlue("http://localhost"+server.Addr))
	log.Fatal(garcon.ListenAndServe(&server))
}

func newServer(port int) http.Server {
	g := garcon.New(
		garcon.WithServerName("Quid"),
		garcon.WithDev(conf.IsDevMode))

	gw = g.Writer

	maxAge := 3600 * 3 // three hours
	if conf.IsDevMode {
		maxAge = 3600 * 24 * 365 // one year
	}
	Incorruptible = g.IncorruptibleChecker(conf.EncodingKey[:32], maxAge, true)

	middleware := garcon.NewChain(
		g.MiddlewareRejectUnprintableURI(),
		g.MiddlewareLogRequest(),
		g.MiddlewareRateLimiter(10, 30),
		g.MiddlewareCORSWithMethodsHeaders(
			[]string{http.MethodGet, http.MethodOptions, http.MethodPost, http.MethodDelete},
			[]string{"Origin", "Content-Type", "Authorization"},
		))

	router := newRouter(g)
	handler := middleware.Then(router)

	return garcon.Server(handler, port, nil)
}

func newRouter(g *garcon.Garcon) http.Handler {
	r := chi.NewRouter()

	// Static website: set the Incorruptible cookie only when visiting index.html
	ws := g.NewStaticWebServer("ui/dist")
	r.NotFound(ws.ServeFile("index.html", "text/html; charset=utf-8"))
	r.Get("/favicon.ico", ws.ServeFile("favicon.ico", "image/x-icon"))
	r.Get("/js/*", ws.ServeDir("text/javascript; charset=utf-8"))
	r.Get("/assets/*", ws.ServeAssets())

	// HTTP Routes
	// public routes
	r.Post("/token/refresh/{timeout}", RequestRefreshToken)
	r.Post("/token/access/{timeout}", RequestAccessToken)
	r.Post("/admin_login", AdminLogin)
	r.Get("/logout", AdminLogout)
	// r.With(Incorruptible.Chk).Post("/admin_token/access/", RequestAdminAccessToken)
	r.Get("/status", status)

	// admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(QuidAdminMiddleware)

		// HTTP API
		r.Route("/groups", func(r chi.Router) {
			r.Post("/add", CreateGroup)
			r.Post("/delete", DeleteGroup)
			r.Post("/info", GroupsInfo)
			r.Post("/add_user", AddUserInGroup)
			r.Post("/remove_user", RemoveUserFromGroup)
			// r.Get("/all", AllGroups) // TODO: remove when old frontend is disabled
			r.Post("/nsall", AllGroupsForNamespace)
		})

		// only admin can see the Git version & commit date.
		r.Get("/version", garcon.ServeVersion())

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
			r.Post("/key", GetNamespaceAccessKey)
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
	})

	// Namespace admin endpoints
	r.Route("/ns", func(r chi.Router) {
		r.Use(NsAdminMiddleware)

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

	return r
}
