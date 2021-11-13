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
	rt := Route{}
	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (rt *Route) AddDefaultGateWay() error {
	link, err := netlink.LinkByName(rt.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %v", err, rt.Link)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(rt.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", rt.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(rt.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", rt.OnLink, err)
	} else {
		if b {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	route := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	if err := netlink.RouteAdd(route); err != nil {
		log.Errorf("Failed to add default GateWay address %s: %v", rt.Gateway, err)
		return err
	}

	return nil
}

func (rt *Route) ReplaceDefaultGateWay() error {
	link, err := netlink.LinkByName(rt.Link)
	if err != nil {
		return err
	}

	ipAddr, _, err := net.ParseCIDR(rt.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %v", rt.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(rt.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %v", rt.OnLink, err)
	} else {
		if b {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	route := &netlink.Route{
		Scope:     netlink.SCOPE_LINK,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	if err := netlink.RouteReplace(route); err != nil {
		log.Errorf("Failed to replace default GateWay address %s: %v", rt.Gateway, err)
		return err
	}

	return nil
}

func (rt *Route) RemoveGateWay() error {
	link, err := netlink.LinkByName(rt.Link)
	if err != nil {
		log.Errorf("Failed to delete default gateway %s: %v", link, err)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(rt.Gateway)
	if err != nil {
		return err
	}

	switch rt.Action {
	case "remove-default-gw":
		rt := &netlink.Route{
			Scope:     netlink.SCOPE_LINK,
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

func (rt *Route) AcquireRoutes(rw http.ResponseWriter) error {
	if rt.Link != "" {
		link, err := netlink.LinkByName(rt.Link)
		if err != nil {
			return err
		}

		routes, err := netlink.RouteList(link, netlink.FAMILY_ALL)
		if err != nil {
			return err
		}
		return web.JSONResponse(routes, rw)
	}

	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	return web.JSONResponse(routes, rw)
}

func (rt *Route) Configure() error {
	switch rt.Action {
	case "add-default-gw":
		return rt.AddDefaultGateWay()
	case "replace-default-gw":
		return rt.ReplaceDefaultGateWay()
	}

	return nil
}
