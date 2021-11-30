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
}

func routerDescribeNetwork(w http.ResponseWriter, r *http.Request) {
	var err error
	s := Describe{}

	s.NetworkDescribe, err = networkd.AcquireNetworkState(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.LinksDescribe, err = networkd.AcquireLinks(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Addresses, err = address.AcquireAddresses()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Routes, err = route.AcquireRoutes()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Links, err = link.AcquireLinks()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Dns, err = resolved.AcquireDns(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Domains, err = resolved.AcquireDomains(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	web.JSONResponse(s, w)
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
