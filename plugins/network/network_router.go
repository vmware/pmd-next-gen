// SPDX-License-Identifier: Apache-2.0

package network

import (
	"github.com/gorilla/mux"

	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/link"
	"github.com/pm-web/plugins/network/netlink/route"
)

func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	link.RegisterRouterLink(n)
	address.RegisterRouterAddress(n)
	route.RegisterRouterRoute(n)
}
