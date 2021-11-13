// SPDX-License-Identifier: Apache-2.0

package route

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/web"
)

type Route struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Gateway string `json:"gateway"`
	OnLink  string `json:"onlink"`
}

func decodeJSONRequest(r *http.Request) (*Route, error) {
	route := Route{}
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (route *Route) AddDefaultGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", err, route.Link)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(route.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", route.OnLink, err)
	} else {
		if b {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	rt := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	if err := netlink.RouteAdd(rt); err != nil {
		log.Errorf("Failed to add default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	return nil
}

func (route *Route) ReplaceDefaultGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(route.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", route.OnLink, err)
	} else {
		if b {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	rt := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	if err := netlink.RouteReplace(rt); err != nil {
		log.Errorf("Failed to replace default GateWay address %s: %v", route.Gateway, err)
		return err
	}

	return nil
}

func (route *Route) DeleteGateWay() error {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to delete default gateway %s: %v", link, err)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		return err
	}

	switch route.Action {
	case "del-default-gw":
		rt := &netlink.Route{
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: link.Attrs().Index,
			Gw:        ipAddr,
		}

		if err = netlink.RouteDel(rt); err != nil {
			log.Errorf("Failed to delete default GateWay address %s: %v", ipAddr, err)
			return err
		}
	}

	return nil
}

func (route *Route) AcquireRoutes(rw http.ResponseWriter) error {
	routes, err := netlink.RouteList(nil, 0)
	if err != nil {
		return err
	}

	return web.JSONResponse(routes, rw)
}

func (route *Route) Configure() error {
	switch route.Action {
	case "add-default-gw":
		return route.AddDefaultGateWay()
	case "replace-default-gw":
		return route.ReplaceDefaultGateWay()
	}

	return nil
}
