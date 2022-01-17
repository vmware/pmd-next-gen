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
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		if err := s.Get(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerSysctlGetPattern(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		if err := s.GetPattern(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerSysctlGetAll(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s := new(Sysctl)
		s.Pattern = ""

		if err := s.GetPattern(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerSysctlUpdate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		if err := s.Update(); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func routerSysctlDelete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		s := new(Sysctl)
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		s.Value = "Delete"
		if err := s.Update(); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func routerSysctlLoad(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT":
		s := new(Sysctl)
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Error decoding request", http.StatusBadRequest)
			return
		}

		if err := s.Load(); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// RegisterRouterSysctl register with mux
func RegisterRouterSysctl(router *mux.Router) {
	s := router.PathPrefix("/sysctl").Subrouter().StrictSlash(false)

	s.HandleFunc("/configstatus", routerSysctlGet)
	s.HandleFunc("/configallstatus", routerSysctlGetAll)
	s.HandleFunc("/configpatternstatus", routerSysctlGetPattern)
	s.HandleFunc("/configupdate", routerSysctlUpdate)
	s.HandleFunc("/configdelete", routerSysctlDelete)
	s.HandleFunc("/configload", routerSysctlLoad)
}
