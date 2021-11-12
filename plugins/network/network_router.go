// SPDX-License-Identifier: Apache-2.0

package network

import (
	"github.com/gorilla/mux"

	"github.com/pm-web/plugins/network/netlink"
)

func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	netlink.RegisterRouterNetlink(n)
}
