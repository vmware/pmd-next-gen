// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
//	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/web"
)

func routerAcquireCommand(w http.ResponseWriter, r *http.Request) {
	var err error

	switch mux.Vars(r)["command"] {
	case "list":
		err = AcquireList(w)
	case "repolist":
		err = AcquireRepoList(w)
	default:
		err = errors.New("not found")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterTdnf(router *mux.Router) {
	n := router.PathPrefix("/tdnf").Subrouter().StrictSlash(false)

	n.HandleFunc("/{command}", routerAcquireCommand).Methods("GET")
}
