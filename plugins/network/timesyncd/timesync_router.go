// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireCurrentNTPServer(w http.ResponseWriter, r *http.Request) {
	if err := AcquireCurrentNTPServer(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireSystemNTPServers(w http.ResponseWriter, r *http.Request) {
	if err := AcquireSystemNTPServers(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterTimeSyncd(router *mux.Router) {
	n := router.PathPrefix("/timesyncd").Subrouter().StrictSlash(false)

	n.HandleFunc("/currentserver", routerAcquireCurrentNTPServer).Methods("GET")
	n.HandleFunc("/systemservers", routerAcquireSystemNTPServers).Methods("GET")
}
