// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package networkd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jaypipes/ghw"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pmd-nextgen/pkg/configfile"
	"github.com/pmd-nextgen/pkg/share"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
)

type MatchSection struct {
	Name string `json:"Name"`
}

type LinkSection struct {
	MTUBytes   string `json:"MTUBytes"`
	MACAddress string `json:"MACAddress"`
	Unmanaged  string `json:"Unmanaged"`
}

type NetworkSection struct {
	DHCP                string   `json:"DHCP"`
	Address             string   `json:"Address"`
	Gateway             string   `json:"Gateway"`
	DNS                 []string `json:"DNS"`
	Domains             []string `json:"Domains"`
	NTP                 []string `json:"NTP"`
	IPv6AcceptRA        string   `json:"IPv6AcceptRA"`
	LinkLocalAddressing string   `json:"LinkLocalAddressing"`
	MulticastDNS        string   `json:"MulticastDNS"`

	VLAN string `json:"VLAN"`
}
type AddressSection struct {
	Address string `json:"Address"`
	Peer    string `json:"Peer"`
	Label   string `json:"Label"`
	Scope   string `json:"Scope"`
}

type RouteSection struct {
	Gateway         string `json:"Gateway"`
	GatewayOnlink   string `json:"GatewayOnlink"`
	Destination     string `json:"Destination"`
	Source          string `json:"Source"`
	PreferredSource string `json:"PreferredSource"`
	Table           string `json:"Table"`
	Scope           string `json:"Scope"`
}

type DHCPv4Section struct {
	ClientIdentifier      string `json:"ClientIdentifier"`
	VendorClassIdentifier string `json:"VendorClassIdentifier"`
	RequestOptions        string `json:"RequestOptions"`
	SendOption            string `json:"SendOption"`
	UseDNS                string `json:"UseDNS"`
	UseNTP                string `json:"UseNTP"`
	UseHostname           string `json:"UseHostname"`
	UseDomains            string `json:"UseDomains"`
	UseRoutes             string `json:"UseRoutes"`
	UseMTU                string `json:"UseMTU"`
	UseGateway            string `json:"UseGateway"`
	UseTimezone           string `json:"UseTimezone"`
	IAID                  string `json:"IAID"`
}

type Network struct {
	Link            string           `json:"Link"`
	LinkSection     LinkSection      `json:"LinkSection"`
	MatchSection    MatchSection     `json:"MatchSection"`
	NetworkSection  NetworkSection   `json:"NetworkSection"`
	DHCPv4Section   DHCPv4Section    `json:"DHCPv4Section"`
	AddressSections []AddressSection `json:"AddressSections"`
	RouteSections   []RouteSection   `json:"RouteSections"`
}

type LinkDescribe struct {
	AddressState     string   `json:"AddressState"`
	AlternativeNames []string `json:"AlternativeNames"`
	CarrierState     string   `json:"CarrierState"`
	Driver           string   `json:"Driver"`
	IPv4AddressState string   `json:"IPv4AddressState"`
	IPv6AddressState string   `json:"IPv6AddressState"`
	Index            int      `json:"Index"`
	LinkFile         string   `json:"LinkFile"`
	Model            string   `json:"Model"`
	Name             string   `json:"Name"`
	OnlineState      string   `json:"OnlineState"`
	OperationalState string   `json:"OperationalState"`
	Path             string   `json:"Path"`
	SetupState       string   `json:"SetupState"`
	Type             string   `json:"Type"`
	Vendor           string   `json:"Vendor"`
	Manufacturer     string   `json:"Manufacturer"`
	NetworkFile      string   `json:"NetworkFile,omitempty"`
}

type LinksDescribe struct {
	Interfaces []LinkDescribe
}

type NetworkDescribe struct {
	AddressState     string   `json:"AddressState"`
	CarrierState     string   `json:"CarrierState"`
	OperationalState string   `json:"OperationalState"`
	OnlineState      string   `json:"OnlineState"`
	IPv4AddressState string   `json:"IPv4AddressState"`
	IPv6AddressState string   `json:"IPv6AddressState"`
	DNS              []string `json:"DNS"`
	Domains          []string `json:"Domains"`
	RouteDomains     []string `json:"RouteDomains"`
	NTP              []string `json:"NTP"`
}

func decodeNetworkJSONRequest(r *http.Request) (*Network, error) {
	n := Network{}
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		return nil, err
	}

	return &n, nil
}

func fillOneLink(link netlink.Link) LinkDescribe {
	l := LinkDescribe{
		Index: link.Attrs().Index,
		Name:  link.Attrs().Name,
		Type:  link.Attrs().EncapType,
	}

	l.AddressState, _ = ParseLinkAddressState(link.Attrs().Index)
	l.IPv4AddressState, _ = ParseLinkIPv4AddressState(link.Attrs().Index)
	l.IPv6AddressState, _ = ParseLinkIPv6AddressState(link.Attrs().Index)
	l.CarrierState, _ = ParseLinkCarrierState(link.Attrs().Index)
	l.OnlineState, _ = ParseLinkOnlineState(link.Attrs().Index)
	l.OperationalState, _ = ParseLinkOperationalState(link.Attrs().Index)
	l.SetupState, _ = ParseLinkSetupState(link.Attrs().Index)
	l.NetworkFile, _ = ParseLinkNetworkFile(link.Attrs().Index)

	c, err := configfile.ParseKeyFromSectionString(path.Join("/sys/class/net", link.Attrs().Name, "device/uevent"), "", "PCI_SLOT_NAME")
	if err == nil {
		pci, err := ghw.PCI()
		if err == nil {
			dev := pci.GetDevice(c)

			l.Model = dev.Product.Name
			l.Vendor = dev.Vendor.Name
			l.Path = "pci-" + dev.Address
		}
	}

	driver, err := configfile.ParseKeyFromSectionString(path.Join("/sys/class/net", link.Attrs().Name, "device/uevent"), "", "DRIVER")
	if err == nil {
		l.Driver = driver
	}

	return l
}

func buildLinkMessageFallback() (*LinksDescribe, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	linkDesc := LinksDescribe{}
	for _, l := range links {
		linkDesc.Interfaces = append(linkDesc.Interfaces, fillOneLink(l))
	}

	return &linkDesc, nil
}

func AcquireLinks(ctx context.Context) (*LinksDescribe, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	links, err := c.DBusLinkDescribe(ctx)
	if err != nil {
		return buildLinkMessageFallback()
	}

	return links, nil
}

func AcquireNetworkState(ctx context.Context) (*NetworkDescribe, error) {
	n := NetworkDescribe{}
	n.AddressState, _ = ParseNetworkAddressState()
	n.IPv4AddressState, _ = ParseNetworkIPv4AddressState()
	n.IPv6AddressState, _ = ParseNetworkIPv6AddressState()
	n.CarrierState, _ = ParseNetworkCarrierState()
	n.OnlineState, _ = ParseNetworkOnlineState()
	n.OperationalState, _ = ParseNetworkOperationalState()
	n.DNS, _ = ParseNetworkDNS()
	n.Domains, _ = ParseNetworkDomains()
	n.RouteDomains, _ = ParseNetworkRouteDomains()
	n.NTP, _ = ParseNetworkNTP()

	return &n, nil
}

func (n *Network) buildNetworkSection(m *configfile.Meta) error {
	if !validator.IsEmpty(n.NetworkSection.DHCP) {
		if validator.IsDHCP(n.NetworkSection.DHCP) {
			m.SetKeySectionString("Network", "DHCP", n.NetworkSection.DHCP)
		} else {
			log.Errorf("Failed to parse DHCP='%s'", n.NetworkSection.DHCP)
			return fmt.Errorf("invalid DHCP='%s'", n.NetworkSection.DHCP)
		}
	}

	if !validator.IsEmpty(n.NetworkSection.Address) {
		if validator.IsIP(n.NetworkSection.Address) {
			m.SetKeySectionString("Network", "Address", n.NetworkSection.Address)
		} else {
			log.Errorf("Failed to parse Address='%s'", n.NetworkSection.Address)
			return fmt.Errorf("invalid Address='%s'", n.NetworkSection.Address)
		}
	}

	if !validator.IsEmpty(n.NetworkSection.Gateway) {
		if validator.IsIP(n.NetworkSection.Gateway) {
			m.SetKeySectionString("Network", "Gateway", n.NetworkSection.Gateway)
		} else {
			log.Errorf("Failed to parse Gateway='%s'", n.NetworkSection.Gateway)
			return fmt.Errorf("invalid Gateway='%s'", n.NetworkSection.Gateway)
		}
	}

	if !validator.IsEmpty(n.NetworkSection.IPv6AcceptRA) && validator.IsBool(n.NetworkSection.IPv6AcceptRA) {
		m.SetKeySectionString("Network", "IPv6AcceptRA", n.NetworkSection.IPv6AcceptRA)
	}

	if !validator.IsEmpty(n.NetworkSection.LinkLocalAddressing) {
		if validator.IsLinkLocalAddressing(n.NetworkSection.LinkLocalAddressing) {
			m.SetKeySectionString("Network", "LinkLocalAddressing", n.NetworkSection.LinkLocalAddressing)
		} else {
			log.Errorf("Failed to parse LinkLocalAddressin='%s'", n.NetworkSection.LinkLocalAddressing)
			return fmt.Errorf("invalid LinkLocalAddressin='%s'", n.NetworkSection.LinkLocalAddressing)
		}
	}

	if !validator.IsEmpty(n.NetworkSection.MulticastDNS) && validator.IsBool(n.NetworkSection.MulticastDNS) {
		m.SetKeySectionString("Network", "MulticastDNS", n.NetworkSection.MulticastDNS)
	}

	if !validator.IsArrayEmpty(n.NetworkSection.Domains) {
		s := m.GetKeySectionString("Network", "Domains")
		t := share.UniqueString(strings.Split(s, " "), n.NetworkSection.NTP)
		m.SetKeySectionString("Network", "Domains", strings.Join(t[:], " "))
	}

	if !validator.IsArrayEmpty(n.NetworkSection.DNS) {
		for _, dns := range n.NetworkSection.DNS {
			if !govalidator.IsDNSName(dns) {
				log.Errorf("Failed to parse DNS='%s'", dns)
				return fmt.Errorf("invalid DNS='%s'", dns)
			}
		}
		s := m.GetKeySectionString("Network", "DNS")
		t := share.UniqueString(strings.Split(s, " "), n.NetworkSection.NTP)
		m.SetKeySectionString("Network", "DNS", strings.Join(t[:], " "))
	}

	if !validator.IsArrayEmpty(n.NetworkSection.NTP) {
		s := m.GetKeySectionString("Network", "NTP")
		t := share.UniqueString(strings.Split(s, " "), n.NetworkSection.NTP)
		m.SetKeySectionString("Network", "NTP", strings.Join(t[:], " "))
	}

	return nil
}

func (n *Network) buildLinkSection(m *configfile.Meta) error {
	if !validator.IsEmpty(n.LinkSection.MTUBytes) {
		if validator.IsMtu(n.LinkSection.MTUBytes) {
			m.SetKeySectionString("Link", "MTUBytes", n.LinkSection.MTUBytes)
		} else {
			log.Errorf("Invalid MTU='%s'", n.LinkSection.MTUBytes)
			return fmt.Errorf("invalid MTU='%s'", n.LinkSection.MTUBytes)
		}
	}

	if !validator.IsEmpty(n.LinkSection.MACAddress) {
		if validator.IsNotMAC(n.LinkSection.MACAddress) {
			log.Errorf("Failed to parse Mac='%s'", n.LinkSection.MACAddress)
			return fmt.Errorf("invalid Address='%s'", n.LinkSection.MACAddress)

		} else {
			m.SetKeySectionString("Link", "MACAddress", n.LinkSection.MACAddress)
		}
	}

	if !validator.IsEmpty(n.LinkSection.Unmanaged) && validator.IsBool(n.LinkSection.Unmanaged) {
		m.SetKeySectionString("Link", "Unmanaged", n.LinkSection.Unmanaged)
	}

	return nil
}

func (n *Network) buildDHCPv4Section(m *configfile.Meta) error {
	if !validator.IsEmpty(n.DHCPv4Section.ClientIdentifier) {
		m.SetKeySectionString("DHCPv4", "ClientIdentifier", n.DHCPv4Section.ClientIdentifier)
	}

	if !validator.IsEmpty(n.DHCPv4Section.VendorClassIdentifier) {
		m.SetKeySectionString("DHCPv4", "VendorClassIdentifier", n.DHCPv4Section.VendorClassIdentifier)
	}

	if !validator.IsEmpty(n.DHCPv4Section.RequestOptions) {
		m.SetKeySectionString("DHCPv4", "RequestOptions", n.DHCPv4Section.RequestOptions)
	}

	if !validator.IsEmpty(n.DHCPv4Section.SendOption) {
		m.SetKeySectionString("DHCPv4", "SendOption", n.DHCPv4Section.SendOption)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseDNS) && validator.IsBool(n.DHCPv4Section.UseDNS) {
		m.SetKeySectionString("DHCPv4", "UseDNS", n.DHCPv4Section.UseDNS)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseDomains) && validator.IsBool(n.DHCPv4Section.UseDomains) {
		m.SetKeySectionString("DHCPv4", "UseDomains", n.DHCPv4Section.UseDomains)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseNTP) && validator.IsBool(n.DHCPv4Section.UseNTP) {
		m.SetKeySectionString("DHCPv4", "UseNTP", n.DHCPv4Section.UseNTP)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseMTU) && validator.IsBool(n.DHCPv4Section.UseMTU) {
		m.SetKeySectionString("DHCPv4", "UseMTU", n.DHCPv4Section.UseMTU)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseGateway) && validator.IsBool(n.DHCPv4Section.UseGateway) {
		m.SetKeySectionString("DHCPv4", "UseGateway", n.DHCPv4Section.UseGateway)
	}

	if !validator.IsEmpty(n.DHCPv4Section.UseTimezone) && validator.IsBool(n.DHCPv4Section.UseTimezone) {
		m.SetKeySectionString("DHCPv4", "UseTimezone", n.DHCPv4Section.UseTimezone)
	}

	if !validator.IsEmpty(n.DHCPv4Section.IAID) && validator.IsIaId(n.DHCPv4Section.IAID) {
		m.SetKeySectionString("DHCPv4", "IAID", n.DHCPv4Section.IAID)
	}

	return nil
}

func (n *Network) buildAddressSection(m *configfile.Meta) error {
	for _, a := range n.AddressSections {
		if err := m.NewSection("Address"); err != nil {
			return err
		}

		if !validator.IsEmpty(a.Address) {
			if validator.IsIP(a.Address) {
				m.SetKeyToNewSectionString("Address", a.Address)
			} else {
				log.Errorf("Failed to parse Address='%s'", a.Address)
				return fmt.Errorf("invalid Address='%s'", a.Address)
			}
		}

		if !validator.IsEmpty(a.Peer) {
			if validator.IsIP(a.Peer) {
				m.SetKeyToNewSectionString("Peer", a.Peer)
			} else {
				log.Errorf("Failed to parse Peer='%s'", a.Peer)
				return fmt.Errorf("invalid Peer='%s'", a.Peer)
			}
		}

		if !validator.IsEmpty(a.Label) {
			m.SetKeyToNewSectionString("Label", a.Label)
		}

		if !validator.IsEmpty(a.Scope) && validator.IsScope(a.Scope) {
			m.SetKeyToNewSectionString("Scope", a.Scope)
		}
	}

	return nil
}

func (n *Network) buildRouteSection(m *configfile.Meta) error {
	for _, rt := range n.RouteSections {
		if err := m.NewSection("Route"); err != nil {
			return err
		}

		if !validator.IsEmpty(rt.Gateway) {
			if validator.IsIP(rt.Gateway) {
				m.SetKeyToNewSectionString("Gateway", rt.Gateway)
			} else {
				log.Errorf("Failed to parse Peer='%s'", rt.Gateway)
				return fmt.Errorf("invalid Peer='%s'", rt.Gateway)
			}
		}

		if !validator.IsEmpty(rt.GatewayOnlink) {
			m.SetKeyToNewSectionString("GatewayOnlink", rt.GatewayOnlink)
		}

		if !validator.IsEmpty(rt.Destination) {
			if validator.IsIP(rt.Destination) {
				m.SetKeyToNewSectionString("Destination", rt.Destination)
			} else {
				log.Errorf("Failed to parse Destination='%s'", rt.Destination)
				return fmt.Errorf("invalid Destination='%s'", rt.Destination)
			}
		}

		if !validator.IsEmpty(rt.Source) {
			if validator.IsIP(rt.Source) {
				m.SetKeyToNewSectionString("Source", rt.Source)
			} else {
				log.Errorf("Failed to parse Source='%s'", rt.Source)
				return fmt.Errorf("invalid Source='%s'", rt.Source)
			}
		}

		if !validator.IsEmpty(rt.PreferredSource) {
			if validator.IsIP(rt.PreferredSource) {
				m.SetKeyToNewSectionString("Source", rt.PreferredSource)
			} else {
				log.Errorf("Failed to parse Source='%s'", rt.PreferredSource)
				return fmt.Errorf("invalid Source='%s'", rt.PreferredSource)
			}
		}

		if !validator.IsEmpty(rt.Table) && govalidator.IsInt(rt.Table) {
			m.SetKeyToNewSectionString("Table", rt.Table)
		}

		if !validator.IsEmpty(rt.Scope) && validator.IsScope(rt.Scope) {
			m.SetKeyToNewSectionString("Scope", rt.Scope)
		}
	}

	return nil
}

func (n *Network) removeAddressSection(m *configfile.Meta) error {
	for _, a := range n.AddressSections {
		if !validator.IsEmpty(a.Address) {
			if err := m.RemoveSection("Address", "Address", a.Address); err != nil {
				log.Errorf("Failed to remove Address='%s': %v", a.Address, err)
				return err
			}
		}
	}

	return nil
}

func (n *Network) removeRouteSection(m *configfile.Meta) error {
	for _, rt := range n.RouteSections {
		if !validator.IsEmpty(rt.Gateway) {
			if err := m.RemoveSection("Route", "Gateway", rt.Gateway); err != nil {
				log.Errorf("Failed to remove Gateway='%s': %v", rt.Gateway, err)
				return err
			}
		}

		if !validator.IsEmpty(rt.Destination) {
			if err := m.RemoveSection("Route", "Destination", rt.Destination); err != nil {
				log.Errorf("Failed to remove Destination='%s': %v", rt.Destination, err)
				return err
			}
		}
	}

	return nil
}

func (n *Network) ConfigureNetwork(ctx context.Context, w http.ResponseWriter) error {
	m, err := CreateOrParseNetworkFile(n.Link)
	if err != nil {
		log.Errorf("Failed to parse network file for link='%s': %v", n.Link, err)
		return err
	}

	if err := n.buildNetworkSection(m); err != nil {
		return err
	}
	if err := n.buildLinkSection(m); err != nil {
		return err
	}
	if err := n.buildDHCPv4Section(m); err != nil {
		return err
	}
	if err := n.buildAddressSection(m); err != nil {
		return err
	}
	if err := n.buildRouteSection(m); err != nil {
		return err
	}

	if err := m.Save(); err != nil {
		log.Errorf("Failed to update config file='%s': %v", m.Path, err)
		return err
	}

	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection with the system bus: %v", err)
		return err
	}
	defer c.Close()

	if err := c.DBusNetworkReload(ctx); err != nil {
		return err
	}

	return web.JSONResponse("configured", w)
}

func (n *Network) RemoveNetwork(ctx context.Context, w http.ResponseWriter) error {
	m, err := CreateOrParseNetworkFile(n.Link)
	if err != nil {
		log.Errorf("Failed to parse network file for link='%s': %v", n.Link, err)
		return err
	}

	if err := n.removeAddressSection(m); err != nil {
		log.Errorf("Failed to remove address section: %v", err)
		return err
	}

	if err := n.removeRouteSection(m); err != nil {
		log.Errorf("Failed to remove route section: %v", err)
		return err
	}

	if err := m.Save(); err != nil {
		log.Errorf("Failed to update config file='%s': %v", m.Path, err)
		return err
	}

	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection with the system bus: %v", err)
		return err
	}
	defer c.Close()

	if err := c.DBusNetworkReload(ctx); err != nil {
		return err
	}

	return web.JSONResponse("removed", w)
}
