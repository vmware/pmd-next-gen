// SPDX-License-Identifier: Apache-2.0

package management

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/distro-management-api/pkg/web"
	"github.com/distro-management-api/plugins/management/group"
	"github.com/distro-management-api/plugins/management/hostname"
	"github.com/distro-management-api/plugins/management/login"
	"github.com/distro-management-api/plugins/management/user"
	"github.com/distro-management-api/plugins/network/netlink/address"
	"github.com/distro-management-api/plugins/network/netlink/route"
	"github.com/distro-management-api/plugins/network/networkd"
	"github.com/distro-management-api/plugins/network/timesyncd"
	"github.com/distro-management-api/plugins/systemd"
)

type Describe struct {
	Hostname          *hostname.Describe        `json:"Hostname"`
	Systemd           *systemd.Describe         `json:"Systemd"`
	NetworkDescribe   *networkd.NetworkDescribe `json:"NetworDescribe"`
	LinksDescribe     *networkd.LinksDescribe   `json:"LinksDescribe"`
	Addresses         []address.AddressInfo     `json:"Addresses"`
	Routes            []route.RouteInfo         `json:"Routes"`
	NTP               *timesyncd.NTPServer      `json:"NTP"`
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
		return
	}

	s.Systemd, err = systemd.ManagerDescribe(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.NetworkDescribe, err = networkd.AcquireNetworkState(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.LinksDescribe, err = networkd.AcquireLinks(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.Addresses, err = address.AcquireAddresses()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.Routes, err = route.AcquireRoutes()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.NTP, err = timesyncd.AcquireNTPServer("currentntpserver", r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.HostInfo, err = host.Info()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.UserStat, err = host.Users()
	if err != nil {
		web.JSONResponseError(err, w)
		return
	}

	s.VirtualMemoryStat, err = mem.VirtualMemoryWithContext(r.Context())
	if err != nil {
		web.JSONResponseError(err, w)
		return
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
