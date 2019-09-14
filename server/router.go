package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"webserver/chatsocket"
	"webserver/controllers"
	"webserver/lib"
	"webserver/middlewares"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	res.SendNotFound()
}

func privateRoute(controller http.HandlerFunc) http.Handler {
	return middlewares.JWTMiddleware().Handler(http.HandlerFunc(controller))
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	hub := chatsocket.NewHub()
	go hub.Run()
	chatsocket.ServeWs(hub, w, r)
}

// NewRouter ...
func NewRouter() *mux.Router {

	// Create main router
	router := mux.NewRouter().StrictSlash(true)
	router.KeepContext = true

	// Handle 404
	router.NotFoundHandler = http.HandlerFunc(notFound)

	/**
	 * meta-data
	 */
	router.Methods("GET").Path("/api/info").HandlerFunc(controllers.GetAPIInfo)

	// API Version
	apiPath := "/api"
	apiVersion := "/v1"
	apiPrefix := apiPath + apiVersion

	// Auth routes
	router.Methods("POST").Path("/auth/login").HandlerFunc(controllers.LoginHandler)
	router.Methods("POST").Path("/auth/login_facebook").HandlerFunc(controllers.LoginHandlerWithoutPassword)

	// User routes
	router.Methods("GET").Path(apiPrefix + "/users").Handler(privateRoute(controllers.GetAllUsersHandler))
	// router.Methods("POST").Path(apiPrefix + "/users").Handler(privateRoute(controllers.CreateUserHandler))
	// router.Methods("POST").Path(apiPrefix + "/users").HandlerFunc(controllers.CreateUserHandler)
	router.Methods("GET").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.GetUserByIDHandler))
	router.Methods("PUT").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.UpdateUserHandler))
	router.Methods("DELETE").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.DeleteUserHandler))
	router.Methods("GET").Path(apiPrefix + "/users").Handler(privateRoute(controllers.GetAllUsersHandler))

	// Confession routes
	router.Methods("GET").Path(apiPrefix + "/admincp/confessions").Handler(privateRoute(controllers.GetAllConfessionsHandler))
	router.Methods("POST").Path(apiPrefix + "/confessions").HandlerFunc(controllers.CreateConfessionHandler)
	router.Methods("POST").Path(apiPrefix + "/myconfess").HandlerFunc(controllers.GetConfessionsBySenderHandler)
	router.Methods("GET").Path(apiPrefix + "/confessions/approved").HandlerFunc(controllers.GetApprovedConfessionsHandler)
	router.Methods("GET").Path(apiPrefix + "/confessions/overview").HandlerFunc(controllers.GetConfessionsOverviewHandler)
	router.Methods("PUT").Path(apiPrefix + "/admincp/confessions/approve").Handler(privateRoute(controllers.ApproveConfessionHandler))
	router.Methods("PUT").Path(apiPrefix + "/admincp/confessions/reject").Handler(privateRoute(controllers.RejectConfessionHandler))
	router.Methods("GET").Path(apiPrefix + "/confessions/search").HandlerFunc(controllers.SearchConfessionsHandler)
	router.Methods("GET").Path(apiPrefix + "/radios").HandlerFunc(controllers.GetRadio)
	router.Methods("POST").Path(apiPrefix + "/radios").Handler(privateRoute(controllers.SetRadio))

	// Crawl
	router.Methods("GET").Path("/crawl/{name}").HandlerFunc(controllers.GetHomeFeedHandler)
	router.Methods("GET").Path("/crawl/{name}/{id}").HandlerFunc(controllers.GetPostFeedHandler)

	// Push
	router.Methods("POST").Path(apiPrefix + "/push/sync").HandlerFunc(controllers.SyncPushIDHandler)

	// Github Gist
	router.Methods("GET").Path("/gist").HandlerFunc(controllers.GetResolveGithubGist)

	router.Path("/ws").HandlerFunc(handleSocket)

	return router
}
