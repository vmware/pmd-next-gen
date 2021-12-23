// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 VMware, Inc.

package login

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/distro-management-api/pkg/web"
)

func routerAcquireUsers(w http.ResponseWriter, r *http.Request) {
	if err := AcquireUsersFromLogin(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireSessions(w http.ResponseWriter, r *http.Request) {
	if err := AcquireSessionsFromLogin(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterLogin(router *mux.Router) {
	s := router.PathPrefix("/login").Subrouter().StrictSlash(false)

	s.HandleFunc("/listusers", routerAcquireUsers).Methods("GET")
	s.HandleFunc("/listsessions", routerAcquireSessions).Methods("GET")
}
