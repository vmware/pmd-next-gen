// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireHostname(w http.ResponseWriter, r *http.Request) {
	if err := AcquireHostnameProperties(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func routerSetHostname(w http.ResponseWriter, r *http.Request) {
	hostname := Hostname{}
	if err := json.NewDecoder(r.Body).Decode(&hostname); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := hostname.SetHostname(r.Context()); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterHostname(router *mux.Router) {
	s := router.PathPrefix("/hostname").Subrouter().StrictSlash(false)

	s.HandleFunc("/property", routerAcquireHostname).Methods("GET")
	s.HandleFunc("/method", routerSetHostname).Methods("POST")
}
