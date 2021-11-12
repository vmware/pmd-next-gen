// SPDX-License-Identifier: Apache-2.0

package address

import (
	"encoding/json"
	"net/http"

	"github.com/vishvananda/netlink"

	"github.com/pm-web/pkg/web"
)

type Address struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Address string `json:"address"`
	Label   string `json:"label"`
}

func decodeJSONRequest(r *http.Request) (*Address, error) {
	address := Address{}

	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		return &address, err
	}

	return &address, nil
}

func (a *Address) Add() error {
	link, err := netlink.LinkByName(a.Link)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(a.Address)
	if err != nil {
		return err
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		return err
	}

	return nil
}

func (a *Address) Remove() error {
	link, err := netlink.LinkByName(a.Link)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(a.Address)
	if err != nil {
		return err
	}

	if err = netlink.AddrDel(link, addr); err != nil {
		return err
	}

	return nil
}

func (a *Address) AcquireAddresses(rw http.ResponseWriter) error {
	if a.Link != "" {
		link, err := netlink.LinkByName(a.Link)
		if err != nil {
			return err
		}

		addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			return err
		}

		return web.JSONResponse(addrs, rw)
	}

	addrs, err := netlink.AddrList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	return web.JSONResponse(addrs, rw)
}
