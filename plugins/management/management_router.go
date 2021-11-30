// SPDX-License-Identifier: Apache-2.0

package management

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/management/group"
	"github.com/pm-web/plugins/management/hostname"
	"github.com/pm-web/plugins/management/login"
	"github.com/pm-web/plugins/management/user"
	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/route"
	"github.com/pm-web/plugins/network/networkd"
	"github.com/pm-web/plugins/systemd"
)

type Describe struct {
	Hostname          *hostname.Describe        `json:"Hostname"`
	Systemd           *systemd.Describe         `json:"Systemd"`
	NetworkDescribe   *networkd.NetworkDescribe `json:"NetworDescribe"`
	LinksDescribe     *networkd.LinksDescribe   `json:"LinksDescribe"`
	Addresses         []address.AddressInfo     `json:"Addresses"`
	Routes            []route.RouteInfo         `json:"Routes"`
	HostInfo          *host.InfoStat            `json:"HostInfo"`
	UserStat          []host.UserStat           `json:"UserStat"`
	VirtualMemoryStat *mem.VirtualMemoryStat    `json:"VirtualMemoryStat"`
}

func routerDescribeSystem(w http.ResponseWriter, r *http.Request) {
	var err error
	s := Describe{}

	s.Hostname, err = hostname.MethodDescribe(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.Systemd, err = systemd.ManagerDescribe(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

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

	s.HostInfo, err = host.Info()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.UserStat, err = host.Users()
	if err != nil {
		web.JSONResponseError(err, w)
	}

	s.VirtualMemoryStat, err = mem.VirtualMemoryWithContext(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
	}

	web.JSONResponse(s, w)
}

func RegisterRouterManagement(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	group.RegisterRouterGroup(n)
	user.RegisterRouterUser(n)

	hostname.RegisterRouterHostname(n)
	login.RegisterRouterLogin(n)

	n.HandleFunc("/describe", routerDescribeSystem).Methods("GET")
}
