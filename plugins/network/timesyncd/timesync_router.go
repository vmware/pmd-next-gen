// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timesyncd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/web"
)

func routerAcquireNTPServers(w http.ResponseWriter, r *http.Request) {
	ntp, err := AcquireNTPServer(mux.Vars(r)["ntpserver"], r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	web.JSONResponse(ntp, w)
}

func routerDescribeNTPServers(w http.ResponseWriter, r *http.Request) {
	ntp, err := DescribeNTPServers(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	web.JSONResponse(ntp, w)
}

func RegisterRouterTimeSyncd(router *mux.Router) {
	n := router.PathPrefix("/timesyncd").Subrouter().StrictSlash(false)

	n.HandleFunc("/describe", routerDescribeNTPServers).Methods("GET")
	n.HandleFunc("/{ntpserver}", routerAcquireNTPServers).Methods("GET")
}
