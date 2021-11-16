// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerConfigureNetwork(w http.ResponseWriter, r *http.Request) {
	n, err := decodeJSONRequest(r)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if err := n.ConfigureNetwork(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterNetworkd(router *mux.Router) {
	n := router.PathPrefix("/networkd").Subrouter().StrictSlash(false)

	n.HandleFunc("/network", routerConfigureNetwork).Methods("POST")
}
