package netlink

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/network/netlink/link"
)

func routerAcquireLinkGet(w http.ResponseWriter, r *http.Request) {
	link := link.Link{
		Name: mux.Vars(r)["link"],
	}

	if err := link.AcquireLink(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterNetlink(router *mux.Router) {
	s := router.PathPrefix("/netlink").Subrouter().StrictSlash(false)

	s.HandleFunc("/link", routerAcquireLinkGet).Methods("GET")
	s.HandleFunc("/link/{link}", routerAcquireLinkGet).Methods("GET")
}
