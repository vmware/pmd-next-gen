// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/web"
)

func routerAcquireCommand(w http.ResponseWriter, r *http.Request) {
	var err error

	switch mux.Vars(r)["command"] {
	case "clean":
		err = AcquireClean(w)
	case "info":
		err = AcquireInfoList(w, "")
	case "list":
		err = AcquireList(w, "")
	case "makecache":
		err = AcquireMakeCache(w)
	case "repolist":
		err = AcquireRepoList(w)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireCommandPkg(w http.ResponseWriter, r *http.Request) {
	pkg := mux.Vars(r)["pkg"]

	var err error
	switch mux.Vars(r)["command"] {
	case "info":
		err = AcquireInfoList(w, pkg)
	case "list":
		err = AcquireList(w, pkg)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterTdnf(router *mux.Router) {
	n := router.PathPrefix("/tdnf").Subrouter().StrictSlash(false)

	n.HandleFunc("/{command}/{pkg}", routerAcquireCommandPkg).Methods("GET")
	n.HandleFunc("/{command}", routerAcquireCommand).Methods("GET")
}
