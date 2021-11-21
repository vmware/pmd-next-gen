// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireLinkProperty(w http.ResponseWriter, r *http.Request) {
	if err := AcquireLinkDNS(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireManagerProperty(w http.ResponseWriter, r *http.Request) {
	if err := AcquireManagerDNS(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterResolved(router *mux.Router) {
	n := router.PathPrefix("/resolved").Subrouter().StrictSlash(false)

	n.HandleFunc("/dns", routerAcquireManagerProperty).Methods("GET")
	n.HandleFunc("/{link}/dns", routerAcquireLinkProperty).Methods("GET")
}
