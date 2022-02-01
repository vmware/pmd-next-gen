// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package networkd

import (
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pmd-nextgen/pkg/configfile"
	"github.com/pmd-nextgen/pkg/system"
)

func ParseLinkString(ifindex int, key string) (string, error) {
	path := "/run/systemd/netif/links/" + strconv.Itoa(ifindex)
	v, err := configfile.ParseKeyFromSectionString(path, "", key)
	if err != nil {
		return "", err
	}

	return v, nil
}

func ParseLinkSetupState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "ADMIN_STATE")
}

func ParseLinkCarrierState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "CARRIER_STATE")
}

func ParseLinkOnlineState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "ONLINE_STATE")
}

func ParseLinkActivationPolicy(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "ACTIVATION_POLICY")
}

func ParseLinkNetworkFile(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "NETWORK_FILE")
}

func ParseLinkOperationalState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "OPER_STATE")
}

func ParseLinkAddressState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "ADDRESS_STATE")
}

func ParseLinkIPv4AddressState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "IPV4_ADDRESS_STATE")
}

func ParseLinkIPv6AddressState(ifindex int) (string, error) {
	return ParseLinkString(ifindex, "IPV6_ADDRESS_STATE")
}

func ParseLinkDNS(ifindex int) ([]string, error) {
	s, err := ParseLinkString(ifindex, "DNS")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseLinkNTP(ifindex int) ([]string, error) {
	s, err := ParseLinkString(ifindex, "NTP")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseLinkDomains(ifindex int) ([]string, error) {
	s, err := ParseLinkString(ifindex, "DOMAINS")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseNetworkState(key string) (string, error) {
	v, err := configfile.ParseKeyFromSectionString("/run/systemd/netif/state", "", key)
	if err != nil {
		return "", err
	}

	return v, nil
}

func ParseNetworkOperationalState() (string, error) {
	return ParseNetworkState("OPER_STATE")
}

func ParseNetworkCarrierState() (string, error) {
	return ParseNetworkState("CARRIER_STATE")
}

func ParseNetworkAddressState() (string, error) {
	return ParseNetworkState("ADDRESS_STATE")
}

func ParseNetworkIPv4AddressState() (string, error) {
	return ParseNetworkState("IPV4_ADDRESS_STATE")
}

func ParseNetworkIPv6AddressState() (string, error) {
	return ParseNetworkState("IPV6_ADDRESS_STATE")
}

func ParseNetworkOnlineState() (string, error) {
	return ParseNetworkState("ONLINE_STATE")
}

func ParseNetworkDNS() ([]string, error) {
	s, err := ParseNetworkState("DNS")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseNetworkNTP() ([]string, error) {
	s, err := ParseNetworkState("NTP")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseNetworkDomains() ([]string, error) {
	s, err := ParseNetworkState("DOMAINS")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func ParseNetworkRouteDomains() ([]string, error) {
	s, err := ParseNetworkState("ROUTE_DOMAINS")
	if err != nil {
		return nil, err
	}

	return strings.Split(s, " "), nil
}

func CreateMatchSection(m *configfile.Meta, link string) error {
	m.NewSection("Match")
	m.SetKeyToNewSectionString("Name", link)

	return nil
}

func CreateNetworkFile(link string) (*configfile.Meta, error) {
	file := "10-" + link + ".network"

	f, err := os.Create(path.Join("/etc/systemd/network", file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := configfile.Load(path.Join("/etc/systemd/network", file))
	if err != nil {
		return nil, err
	}

	if err := CreateMatchSection(m, link); err != nil {
		return nil, err
	}
	return m, nil
}

func CreateOrParseNetworkFile(l string) (*configfile.Meta, error) {
	link, err := netlink.LinkByName(l)
	if err != nil {
		return nil, err
	}

	if _, err := ParseLinkSetupState(link.Attrs().Index); err != nil {
		m, err := CreateNetworkFile(link.Attrs().Name)
		if err != nil {
			return nil, err
		}

		system.ChangePermission("systemd-network", m.Path)
		return m, nil
	}

	n, err := ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		m, err := CreateNetworkFile(link.Attrs().Name)
		if err != nil {
			return nil, err
		}

		system.ChangePermission("systemd-network", m.Path)
		return m, nil
	}

	return configfile.Load(n)
}

func CreateNetDevFile(link string, kind string) (*configfile.Meta, error) {
	file := "10-" + link + "-" + kind + ".netdev"

	f, err := os.Create(path.Join("/etc/systemd/network", file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := configfile.Load(path.Join("/etc/systemd/network", file))
	if err != nil {
		return nil, err
	}

	system.ChangePermission("systemd-network", m.Path)
	return m, nil
}

func CreateNetDevNetworkFile(link string, kind string) error {
	file := "10-" + link + "-" + kind + ".network"

	f, err := os.Create(path.Join("/etc/systemd/network", file))
	if err != nil {
		return err
	}
	defer f.Close()

	m, err := configfile.Load(path.Join("/etc/systemd/network", file))
	if err != nil {
		return err
	}

	if err := CreateMatchSection(m, link); err != nil {
		return err
	}

	if err := m.Save(); err != nil {
		log.Errorf("Failed to update config file='%s': %v", m.Path, err)
		return err
	}

	system.ChangePermission("systemd-network", m.Path)
	return nil
}

func CreateOrParseLinkFile(link string) (*configfile.Meta, error) {
	file := "10-" + link + ".link"

	var m *configfile.Meta
	var err error
	if !system.PathExists(path.Join("/etc/systemd/network", file)) {
		f, err := os.Create(path.Join("/etc/systemd/network", file))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		m, err = configfile.Load(path.Join("/etc/systemd/network", file))
		if err != nil {
			return nil, err
		}

		if err := CreateMatchSection(m, link); err != nil {
			return nil, err
		}
		system.ChangePermission("systemd-network", m.Path)
	} else {
		m, err = configfile.Load(path.Join("/etc/systemd/network", file))
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}
