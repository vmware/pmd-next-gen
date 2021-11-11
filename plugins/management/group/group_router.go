// SPDX-License-Identifier: Apache-2.0

package group

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerGroupAdd(w http.ResponseWriter, r *http.Request) {
	g := new(Group)
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := g.GroupAdd(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerGroupModify(w http.ResponseWriter, r *http.Request) {
	g := new(Group)
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := g.GroupModify(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerGroupRemove(w http.ResponseWriter, r *http.Request) {
	g := new(Group)
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := g.GroupRemove(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterGroup(router *mux.Router) {
	s := router.PathPrefix("/group").Subrouter().StrictSlash(false)

	s.HandleFunc("/add", routerGroupAdd).Methods("POST")
	s.HandleFunc("/delete", routerGroupRemove).Methods("DELETE")
	s.HandleFunc("/modify", routerGroupModify).Methods("PUT")
}