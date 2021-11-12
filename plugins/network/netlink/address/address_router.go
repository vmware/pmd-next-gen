package address

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireAddress(w http.ResponseWriter, r *http.Request) {
	link := Address{
		Link: mux.Vars(r)["link"],
	}

	if err := link.AcquireAddresses(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterAddress(router *mux.Router) {
	s := router.PathPrefix("/netlink").Subrouter().StrictSlash(false)

	s.HandleFunc("/address", routerAcquireAddress).Methods("GET")
	s.HandleFunc("/address/{link}", routerAcquireAddress).Methods("GET")
}
