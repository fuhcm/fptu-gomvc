package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gosu-team/fptu-api/controllers"
	"github.com/gosu-team/fptu-api/lib"
	"github.com/gosu-team/fptu-api/middlewares"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	res.SendNotFound()
}

func privateRoute(controller http.HandlerFunc) http.Handler {
	return middlewares.JWTMiddleware().Handler(http.HandlerFunc(controller))
}

// NewRouter ...
func NewRouter() *mux.Router {

	// Create main router
	mainRouter := mux.NewRouter().StrictSlash(true)
	mainRouter.KeepContext = true

	// Handle 404
	mainRouter.NotFoundHandler = http.HandlerFunc(notFound)

	/**
	 * meta-data
	 */
	mainRouter.Methods("GET").Path("/api/info").HandlerFunc(controllers.GetAPIInfo)

	/**
	 * /users
	 */
	// usersRouter.HandleFunc("/", l.Use(c.GetAllUsersHandler, m.SaySomething())).Methods("GET")

	// API Version
	apiPath := "/api"
	apiVersion := "/v1"
	apiPrefix := apiPath + apiVersion

	// Auth routes
	mainRouter.Methods("POST").Path("/auth/login").HandlerFunc(controllers.LoginHandler)
	mainRouter.Methods("POST").Path("/auth/login_facebook").HandlerFunc(controllers.LoginHandlerWithoutPassword)

	// User routes
	mainRouter.Methods("GET").Path(apiPrefix + "/users").Handler(privateRoute(controllers.GetAllUsersHandler))
	// mainRouter.Methods("POST").Path(apiPrefix + "/users").Handler(privateRoute(controllers.CreateUserHandler))
	// mainRouter.Methods("POST").Path(apiPrefix + "/users").HandlerFunc(controllers.CreateUserHandler)
	mainRouter.Methods("GET").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.GetUserByIDHandler))
	mainRouter.Methods("PUT").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.UpdateUserHandler))
	mainRouter.Methods("DELETE").Path(apiPrefix + "/users/{id}").Handler(privateRoute(controllers.DeleteUserHandler))

	// Confession routes
	mainRouter.Methods("GET").Path(apiPrefix + "/admincp/confessions").Handler(privateRoute(controllers.GetAllConfessionsHandler))
	mainRouter.Methods("POST").Path(apiPrefix + "/confessions").HandlerFunc(controllers.CreateConfessionHandler)
	mainRouter.Methods("POST").Path(apiPrefix + "/myconfess").HandlerFunc(controllers.GetConfessionsBySenderHandler)
	mainRouter.Methods("GET").Path(apiPrefix + "/confessions/approved").HandlerFunc(controllers.GetApprovedConfessionsHandler)
	mainRouter.Methods("GET").Path(apiPrefix + "/confessions/overview").HandlerFunc(controllers.GetConfessionsOverviewHandler)
	mainRouter.Methods("PUT").Path(apiPrefix + "/admincp/confessions/approve").Handler(privateRoute(controllers.ApproveConfessionHandler))
	mainRouter.Methods("PUT").Path(apiPrefix + "/admincp/confessions/rollback_approve").Handler(privateRoute(controllers.RollbackApproveConfessionHandler))
	mainRouter.Methods("PUT").Path(apiPrefix + "/admincp/confessions/reject").Handler(privateRoute(controllers.RejectConfessionHandler))
	mainRouter.Methods("GET").Path(apiPrefix + "/confessions/search").HandlerFunc(controllers.SearchConfessionsHandler)

	// Get NextID
	mainRouter.Methods("GET").Path(apiPrefix + "/next_confession_id").HandlerFunc(controllers.GetNextConfessionNextIDHandler)

	// Crawl
	mainRouter.Methods("GET").Path("/crawl/{name}").HandlerFunc(controllers.GetHomeFeedHandler)
	mainRouter.Methods("GET").Path("/crawl/{name}/{id}").HandlerFunc(controllers.GetPostFeedHandler)

	return mainRouter
}
