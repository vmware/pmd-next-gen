// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/distro-management-api/pkg/web"
)

func routerAcquireNTPServers(w http.ResponseWriter, r *http.Request) {
	ntp, err := AcquireNTPServer(mux.Vars(r)["ntpserver"], r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	web.JSONResponse(ntp, w)
}

func RegisterRouterTimeSyncd(router *mux.Router) {
	n := router.PathPrefix("/timesyncd").Subrouter().StrictSlash(false)

	n.HandleFunc("/{ntpserver}", routerAcquireNTPServers).Methods("GET")
}
