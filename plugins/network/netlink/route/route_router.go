package route

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAddRoute(w http.ResponseWriter, r *http.Request) {
	rt, err := decodeJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := rt.Configure(); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerDeleteRoute(w http.ResponseWriter, r *http.Request) {
	rt, err := decodeJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err = rt.RemoveGateWay(); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireRoute(w http.ResponseWriter, r *http.Request) {
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
	s.HandleFunc("/route", routerAcquireRoute).Methods("GET")
	s.HandleFunc("/route/{link}", routerAcquireRoute).Methods("GET")
}