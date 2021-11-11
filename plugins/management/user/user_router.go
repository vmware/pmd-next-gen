// SPDX-License-Identifier: Apache-2.0

package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAddUser(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := u.Add(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerModifyUser(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := u.Modify(w); err != nil {
		web.JSONResponseError(err, w)
		return
	}
}

func routerRemoveUser(w http.ResponseWriter, r *http.Request) {
	u := new(User)
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := u.Remove(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterUser(router *mux.Router) {
	s := router.PathPrefix("/user").Subrouter().StrictSlash(false)

	s.HandleFunc("/add", routerAddUser).Methods("POST")
	s.HandleFunc("/delete", routerRemoveUser).Methods("DELETE")
	s.HandleFunc("/modify", routerModifyUser).Methods("PUT")
}