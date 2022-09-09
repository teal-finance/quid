package api

import (
	"net/http"

	color "github.com/acmacalister/skittles"
	"github.com/go-chi/chi/v5"

	"github.com/teal-finance/garcon"
	"github.com/teal-finance/incorruptible"
	"github.com/teal-finance/quid/crypt"
)

var Incorruptible *incorruptible.Incorruptible

var gw garcon.Writer

// RunServer : configure and run the server.
func RunServer(port int, devMode bool) {
	server := newServer(port, devMode)

	if devMode {
		log.Info("" + color.BoldRed("Running in development mode"))
	}

	log.Print("Server listening on " + color.UnderlineBlue("http://localhost"+server.Addr))
	log.Fatal(garcon.ListenAndServe(&server))
}

func newServer(port int, devMode bool) http.Server {
	g := garcon.New(
		garcon.WithServerName("Quid"),
		garcon.WithDev(devMode))

	gw = g.Writer

	maxAge := 3600 * 3 // three hours
	if devMode {
		maxAge = 3600 * 24 * 365 // one year
	}
	Incorruptible = g.IncorruptibleCheckerBin(crypt.EncodingKey[:16], maxAge, true)

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

	// public routes: not protected by login cookie
	r.Post("/token/refresh/{timeout}", requestRefreshToken)
	r.Post("/token/access/{timeout}", requestAccessToken)
	r.Post("/token/public", getAccessPublicKey)
	r.Post("/admin_login", adminLogin)
	r.Get("/logout", adminLogout)
	r.Get("/status", status)

	// Quid admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Use(quidAdminMiddleware)

		// HTTP API
		r.Route("/groups", func(r chi.Router) {
			r.Post("/add", createGroup)
			r.Post("/delete", deleteGroup)
			r.Post("/info", groupsInfo)
			r.Post("/add_user", addUserInGroup)
			r.Post("/remove_user", removeUserFromGroup)
			r.Post("/nsall", allNsGroups)
		})

		// only admin can see the Git version & commit date.
		r.Get("/version", garcon.ServeVersion())

		r.Route("/users", func(r chi.Router) {
			r.Post("/add", createUser)
			r.Post("/delete", deleteUser)
			r.Post("/groups", userGroupsInfo)
			r.Post("/orgs", userOrgsInfo)
			r.Post("/nsall", allNsUsers)
		})

		r.Route("/namespaces", func(r chi.Router) {
			r.Post("/add", createNamespace)
			r.Post("/delete", deleteNamespace)
			r.Post("/find", findNamespace)
			r.Post("/info", namespaceInfo)
			r.Post("/key", getAccessVerificationKey)
			r.Post("/max-ttl", setTokenMaxTTL)
			r.Post("/max-refresh-ttl", setRefreshMaxTTL)
			r.Post("/groups", nsGroups)
			r.Post("/endpoint", enableNsEndpoint)
			r.Get("/all", allNamespaces)
		})

		r.Route("/orgs", func(r chi.Router) {
			r.Get("/all", allOrgs)
			r.Post("/add", createOrg)
			r.Post("/delete", deleteOrg)
			r.Post("/find", findOrg)
			r.Post("/add_user", addUserInOrg)
			r.Post("/remove_user", removeUserFromOrg)
		})

		r.Route("/nsadmin", func(r chi.Router) {
			r.Post("/add", createAdministrators)
			r.Post("/nsall", allNsAdministrators)
			r.Post("/delete", deleteAdministrator)
			r.Post("/search/nonadmins", listNonAdminUsersInNs)
		})
	})

	// Namespace admin endpoints
	r.Route("/ns", func(r chi.Router) {
		r.Use(nsAdminMiddleware)

		r.Post("/valid", validAccessToken)

		// nsadmin users
		r.Route("/users", func(r chi.Router) {
			r.Post("/add", createUser)
			r.Post("/delete", deleteUser)
			r.Post("/groups", userGroupsInfo)
			r.Post("/nsall", allNsUsers)
		})

		// nsadmin groups
		r.Route("/groups", func(r chi.Router) {
			r.Post("/add", createGroup)
			r.Post("/delete", deleteGroup)
			r.Post("/info", groupsInfo)
			r.Post("/add_user", addUserInGroup)
			r.Post("/remove_user", removeUserFromGroup)
			r.Post("/nsall", allNsGroups)
		})
	})

	return r
}
