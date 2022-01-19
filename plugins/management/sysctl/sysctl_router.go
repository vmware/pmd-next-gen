// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package sysctl

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmd-nextgen/pkg/web"
)

func routerSysctlGet(w http.ResponseWriter, r *http.Request) {
	s := Sysctl{}
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := s.Get(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerSysctlGetPattern(w http.ResponseWriter, r *http.Request) {
	s := Sysctl{}
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := s.GetPattern(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerSysctlGetAll(w http.ResponseWriter, r *http.Request) {
	s := Sysctl{
		Pattern: "",
	}

	if err := s.GetPattern(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerSysctlUpdate(w http.ResponseWriter, r *http.Request) {
	s := Sysctl{}
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := s.Update(); err != nil {
		web.JSONResponseError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func routerSysctlDelete(w http.ResponseWriter, r *http.Request) {
	s := Sysctl{}
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	s.Value = "Delete"
	if err := s.Update(); err != nil {
		web.JSONResponseError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func routerSysctlLoad(w http.ResponseWriter, r *http.Request) {
	s := new(Sysctl)
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := s.Load(); err != nil {
		web.JSONResponseError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RegisterRouterSysctl register with mux
func RegisterRouterSysctl(router *mux.Router) {
	s := router.PathPrefix("/sysctl").Subrouter().StrictSlash(false)

	s.HandleFunc("/configstatus", routerSysctlGet).Methods("GET")
	s.HandleFunc("/configallstatus", routerSysctlGetAll).Methods("GET")
	s.HandleFunc("/configpatternstatus", routerSysctlGetPattern).Methods("GET")
	s.HandleFunc("/configupdate", routerSysctlUpdate).Methods("POST")
	s.HandleFunc("/configdelete", routerSysctlDelete).Methods("DELETE")
	s.HandleFunc("/configload", routerSysctlLoad).Methods("POST")
}
