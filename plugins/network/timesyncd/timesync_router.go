// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

func routerAcquireNTPServers(w http.ResponseWriter, r *http.Request) {
	if err := AcquireNTPServer(mux.Vars(r)["ntpserver"], r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterTimeSyncd(router *mux.Router) {
	n := router.PathPrefix("/timesyncd").Subrouter().StrictSlash(false)

	n.HandleFunc("/{ntpserver}", routerAcquireNTPServers).Methods("GET")
}
