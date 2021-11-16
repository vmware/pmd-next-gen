// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/vishvananda/netlink"

	"github.com/pm-web/pkg/config"
	"github.com/pm-web/pkg/web"
)

type MatchSection struct {
	Name string `json:"Name"`
}
type NetworkSection struct {
	DHCP                string `json:"DHCP"`
	DNS                 string `json:"DNS"`
	Domains             string `json:"Domains"`
	NTP                 string `json:"NTP"`
	IPv6AcceptRA        string `json:"IPv6AcceptRA"`
	LinkLocalAddressing string `json:"LinkLocalAddressing"`
}

type Network struct {
	Link           string         `json:"Link"`
	MatchSection   MatchSection   `json:"MatchSection"`
	NetworkSection NetworkSection `json:"NetworkSection"`
}

func decodeJSONRequest(r *http.Request) (*Network, error) {
	n := Network{}
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		return &n, err
	}

	return &n, nil
}

func (n *Network) ConfigureNetworkSection(m *config.Meta) {
	if n.NetworkSection.DHCP != "" {
		m.SetKeySectionString("Network", "DHCP", n.NetworkSection.DHCP)
	}
}

func (n *Network) ConfigureNetwork(w http.ResponseWriter) error {
	link, err := netlink.LinkByName(n.Link)
	if err != nil {
		return err
	}

	network, err := CreateOrParseNetworkFile(link)
	if err != nil {
		return err
	}

	m, err := config.Load(path.Join("/etc/systemd/network", network))
	if err != nil {
		return err
	}

	n.ConfigureNetworkSection(m)
	m.Save()

	return web.JSONResponse("configured", w)
}
