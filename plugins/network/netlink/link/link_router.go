// SPDX-License-Identifier: Apache-2.0

package link

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/distro-management-api/pkg/web"
)

func routerAcquireLink(w http.ResponseWriter, r *http.Request) {
	links, err := AcquireLinks()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	web.JSONResponse(links, w)
}

func RegisterRouterLink(router *mux.Router) {
	s := router.PathPrefix("/netlink").Subrouter().StrictSlash(false)

	s.HandleFunc("/link", routerAcquireLink).Methods("GET")
}
