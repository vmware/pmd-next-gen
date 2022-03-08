// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package firewall

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/web"
)

func routerShowTables(w http.ResponseWriter, r *http.Request) {
	t, err := decodeNftJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := t.ShowTable(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerSaveTables(w http.ResponseWriter, r *http.Request) {
	t, err := decodeNftJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := t.SaveTable(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAddTables(w http.ResponseWriter, r *http.Request) {
	t, err := decodeNftJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := t.AddTable(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerRemoveTables(w http.ResponseWriter, r *http.Request) {
	t, err := decodeNftJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := t.RemoveTable(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterNft(router *mux.Router) {
	n := router.PathPrefix("/firewall/nft/").Subrouter().StrictSlash(false)

	n.HandleFunc("/tables/add", routerAddTables).Methods("POST")
	n.HandleFunc("/tables/show", routerShowTables).Methods("GET")
	n.HandleFunc("/tables/save", routerSaveTables).Methods("PUT")
	n.HandleFunc("/tables/remove", routerRemoveTables).Methods("DELETE")
}
