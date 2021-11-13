package route

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAddRoute(w http.ResponseWriter, r *http.Request) {
	route, err := decodeJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := route.Configure(); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerDeleteRoute(w http.ResponseWriter, r *http.Request) {
	route, err := decodeJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err = route.DeleteGateWay(); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireRoutes(w http.ResponseWriter, r *http.Request) {
	rt := Route{
		Link: mux.Vars(r)["link"],
	}

	if err := rt.AcquireRoutes(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterRoute(router *mux.Router) {
	s := router.PathPrefix("/netlink").Subrouter().StrictSlash(false)

	s.HandleFunc("/route/{link}", routerAddRoute).Methods("POST")
	s.HandleFunc("/route/{link}", routerDeleteRoute).Methods("DELETE")
	s.HandleFunc("/route", routerGetRoute).Methods("GET")
	s.HandleFunc("/route/{link}", routerAcquireRoutes).Methods("GET")
}
