// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireLinkDNS(w http.ResponseWriter, r *http.Request) {
	if err := AcquireLinkDNS(r.Context(), mux.Vars(r)["link"], w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireLinkDomains(w http.ResponseWriter, r *http.Request) {
	if err := AcquireLinkDomains(r.Context(), mux.Vars(r)["link"], w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireDNS(w http.ResponseWriter, r *http.Request) {
	if err := AcquireDNSFromResolveManager(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireDomains(w http.ResponseWriter, r *http.Request) {
	if err := AcquireDomainsFromResolveManager(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterResolved(router *mux.Router) {
	n := router.PathPrefix("/resolved").Subrouter().StrictSlash(false)

	n.HandleFunc("/dns", routerAcquireDNS).Methods("GET")
	n.HandleFunc("/domains", routerAcquireDomains).Methods("GET")
	n.HandleFunc("/{link}/dns", routerAcquireLinkDNS).Methods("GET")
	n.HandleFunc("/{link}/domains", routerAcquireLinkDomains).Methods("GET")
}
