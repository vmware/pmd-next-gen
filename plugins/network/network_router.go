// SPDX-License-Identifier: Apache-2.0

package network

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/link"
	"github.com/pm-web/plugins/network/netlink/route"
	"github.com/pm-web/plugins/network/networkd"
	"github.com/pm-web/plugins/network/resolved"
	"github.com/pm-web/plugins/network/timesyncd"
)

type Describe struct {
	NetworkDescribe *networkd.NetworkDescribe `json:"NetworDescribe"`
	LinksDescribe   *networkd.LinksDescribe   `json:"LinksDescribe"`
	Links           []link.LinkInfo           `json:"links"`
	Addresses       []address.AddressInfo     `json:"Addresses"`
	Routes          []route.RouteInfo         `json:"Routes"`
	Dns             []resolved.Dns            `json:"Dns"`
	Domains         []resolved.Domains        `json:"Domains"`
	NTP             *timesyncd.NTPServer      `json:"NTP"`
}

func routerDescribeNetwork(w http.ResponseWriter, r *http.Request) {
	var err error
	n := Describe{}

	n.NetworkDescribe, err = networkd.AcquireNetworkState(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.LinksDescribe, err = networkd.AcquireLinks(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.Addresses, err = address.AcquireAddresses()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.Routes, err = route.AcquireRoutes()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.Links, err = link.AcquireLinks()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.Dns, err = resolved.AcquireDns(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.Domains, err = resolved.AcquireDomains(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	n.NTP, err = timesyncd.AcquireNTPServer("linkntpservers", r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	web.JSONResponse(n, w)
}

func RegisterRouterNetwork(router *mux.Router) {
	n := router.PathPrefix("/network").Subrouter()

	// netlink
	link.RegisterRouterLink(n)
	address.RegisterRouterAddress(n)
	route.RegisterRouterRoute(n)

	// systemd-networkd
	networkd.RegisterRouterNetworkd(n)
	// systemd-resolved
	resolved.RegisterRouterResolved(n)
	// systemd-timesynd
	timesyncd.RegisterRouterTimeSyncd(n)

	n.HandleFunc("/describe", routerDescribeNetwork).Methods("GET")
}
