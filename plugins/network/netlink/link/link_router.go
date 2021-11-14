// SPDX-License-Identifier: Apache-2.0

package link

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireLink(w http.ResponseWriter, r *http.Request) {
	link := Link{
		Name: mux.Vars(r)["link"],
	}

	if err := link.AcquireLink(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterLink(router *mux.Router) {
	s := router.PathPrefix("/netlink").Subrouter().StrictSlash(false)

	s.HandleFunc("/link", routerAcquireLink).Methods("GET")
	s.HandleFunc("/link/{link}", routerAcquireLink).Methods("GET")
}
