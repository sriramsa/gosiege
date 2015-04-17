// Package listener provides routes used for handling the REST APIs
package listener

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Route Handler for using with http.Handler
type RouteHandler struct {
	function func(http.ResponseWriter, *http.Request)
}

func (handler RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.function(w, r)
}

// Registers the routes for the REST API calls
func regApiRoutes() {
	r := mux.NewRouter()

	regSessionRoutes(r)

	regNodeRoutes(r)

	http.Handle("/", r)

}

func regSessionRoutes(r *mux.Router) {
	sr := r.PathPrefix("/gosiege/sessions").Subrouter()

	sr.Path("/new").
		Methods("PUT").
		Handler(handlers.CombinedLoggingHandler(os.Stdout, RouteHandler{safelyDo(newSessHandler)}))

	sr.Path("/stop/{Id:[0-9]+}").
		Methods("GET", "PATCH").
		Handler(handlers.CombinedLoggingHandler(os.Stdout, RouteHandler{safelyDo(stopSessHandler)}))

	sr.Path("/update/{Id:[0-9]+}").
		Methods("GET", "PATCH", "PUT").
		Handler(handlers.CombinedLoggingHandler(os.Stdout, RouteHandler{safelyDo(updateSessHandler)}))
}

func regNodeRoutes(r *mux.Router) {
	nr := r.PathPrefix("/gosiege/nodes").Subrouter()

	// Node Routes
	nr.Path("/new").
		Methods("GET").
		Handler(handlers.CombinedLoggingHandler(os.Stdout, RouteHandler{newSessHandler}))
}
