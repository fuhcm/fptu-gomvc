package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gosu-team/cfapp-api/controllers"
	"github.com/gosu-team/cfapp-api/lib"
	"github.com/gosu-team/cfapp-api/middlewares"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	res.SendOK("Go server is running, this is default page, also a notfound page.")
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
	mainRouter.Methods("GET").Path(apiPrefix + "/confessions/overview").HandlerFunc(controllers.GetConfessionsOverviewHandler)
	mainRouter.Methods("PUT").Path(apiPrefix + "/admincp/confessions/approve").Handler(privateRoute(controllers.ApproveConfessionHandler))
	mainRouter.Methods("PUT").Path(apiPrefix + "/admincp/confessions/reject").Handler(privateRoute(controllers.RejectConfessionHandler))

	return mainRouter
}
