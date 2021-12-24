// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package resolved

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/distro-management-api/pkg/web"
)

func routerAcquireLinkDNS(w http.ResponseWriter, r *http.Request) {
	if err := AcquireLinkDns(r.Context(), mux.Vars(r)["link"], w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireLinkDomains(w http.ResponseWriter, r *http.Request) {
	if err := AcquireLinkDomains(r.Context(), mux.Vars(r)["link"], w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireDNS(w http.ResponseWriter, r *http.Request) {
	dns, err := AcquireDns(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	web.JSONResponse(dns, w)
}

func routerAcquireDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := AcquireDomains(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	web.JSONResponse(domains, w)
}

func RegisterRouterResolved(router *mux.Router) {
	n := router.PathPrefix("/resolved").Subrouter().StrictSlash(false)

	n.HandleFunc("/dns", routerAcquireDNS).Methods("GET")
	n.HandleFunc("/domains", routerAcquireDomains).Methods("GET")
	n.HandleFunc("/{link}/dns", routerAcquireLinkDNS).Methods("GET")
	n.HandleFunc("/{link}/domains", routerAcquireLinkDomains).Methods("GET")
}
