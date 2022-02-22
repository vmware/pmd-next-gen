// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/fatih/color"
	"github.com/pmd-nextgen/pkg/share"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network"
	"github.com/pmd-nextgen/plugins/network/netlink/address"
	"github.com/pmd-nextgen/plugins/network/netlink/link"
	"github.com/pmd-nextgen/plugins/network/netlink/route"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/pmd-nextgen/plugins/network/resolved"
	"github.com/pmd-nextgen/plugins/network/timesyncd"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/urfave/cli/v2"
)

type NetDevIOCounters struct {
	Success bool                 `json:"success"`
	Message []net.IOCountersStat `json:"message"`
	Errors  string               `json:"errors"`
}

type Interface struct {
	Success bool                `json:"success"`
	Message []net.InterfaceStat `json:"message"`
	Errors  string              `json:"errors"`
}

type NetworkDescribe struct {
	Success bool             `json:"success"`
	Message network.Describe `json:"message"`
	Errors  string           `json:"errors"`
}

type ResolveDescribe struct {
	Success bool              `json:"success"`
	Message resolved.Describe `json:"message"`
	Errors  string            `json:"errors"`
}

func displayInterfaces(i *Interface) {
	for _, n := range i.Message {
		fmt.Printf("            %v %v\n", color.HiBlueString("Name:"), n.Name)
		fmt.Printf("           %v %v\n", color.HiBlueString("Index:"), n.Index)
		fmt.Printf("             %v %v\n", color.HiBlueString("MTU:"), n.MTU)

		fmt.Printf("           %v", color.HiBlueString("Flags:"))
		for _, j := range n.Flags {
			fmt.Printf(" %v", j)
		}
		fmt.Printf("\n")

		fmt.Printf("%v %v\n", color.HiBlueString("Hardware Address:"), n.HardwareAddr)

		fmt.Printf("       %v", color.HiBlueString("Addresses:"))
		for _, j := range n.Addrs {
			fmt.Printf(" %v", j.Addr)
		}
		fmt.Printf("\n\n")
	}
}

func displayNetDevIOStatistics(netDev *NetDevIOCounters) {
	for _, n := range netDev.Message {
		fmt.Printf("            %v %v\n", color.HiBlueString("Name:"), n.Name)
		fmt.Printf("%v %v\n", color.HiBlueString("Packets received:"), n.PacketsRecv)
		fmt.Printf("%v %v\n", color.HiBlueString("  Bytes received:"), n.PacketsSent)
		fmt.Printf("%v %v\n", color.HiBlueString("      Bytes sent:"), n.PacketsSent)
		fmt.Printf("%v %v\n", color.HiBlueString("         Drop in:"), n.PacketsSent)
		fmt.Printf("%v %v\n", color.HiBlueString("        Drop out:"), n.Dropin)
		fmt.Printf("%v %v\n", color.HiBlueString("        Error in:"), n.Dropout)
		fmt.Printf("%v %v\n", color.HiBlueString("       Error out:"), n.Errout)
		fmt.Printf("%v %v\n", color.HiBlueString("         Fifo in:"), n.Fifoin)
		fmt.Printf("%v %v\n\n", color.HiBlueString("        Fifo out:"), n.Fifoout)
	}
}

func displayOneLinkNetworkStatus(l *networkd.LinkDescribe) {
	fmt.Printf("             %v %v\n", color.HiBlueString("Name:"), l.Name)
	if len(l.AlternativeNames) > 0 {
		fmt.Printf("%v %v\n", color.HiBlueString("Alternative Names:"), strings.Join(l.AlternativeNames, " "))
	}
	fmt.Printf("            %v %v\n", color.HiBlueString("Index:"), l.Index)
	if l.LinkFile != "" {
		fmt.Printf("        %v %v\n", color.HiBlueString("Link File:"), l.LinkFile)
	}
	if l.NetworkFile != "" {
		fmt.Printf("     %v %v\n", color.HiBlueString("Network File:"), l.NetworkFile)
	}
	fmt.Printf("             %v %v\n", color.HiBlueString("Type:"), l.Type)
	fmt.Printf("            %v %v (%v)\n", color.HiBlueString("State:"), l.OperationalState, l.SetupState)
	if l.Driver != "" {
		fmt.Printf("           %v %v\n", color.HiBlueString("Driver:"), l.Driver)
	}
	if l.Vendor != "" {
		fmt.Printf("           %v %v\n", color.HiBlueString("Vendor:"), l.Vendor)
	}
	if l.Model != "" {
		fmt.Printf("            %v %v\n", color.HiBlueString("Model:"), l.Model)
	}
	if l.Path != "" {
		fmt.Printf("             %v %v\n", color.HiBlueString("Path:"), l.Path)
	}
	fmt.Printf("    %v %v\n", color.HiBlueString("Carrier State:"), l.CarrierState)

	if l.OnlineState != "" {
		fmt.Printf("     %v %v\n", color.HiBlueString("Online State:"), l.OnlineState)
	}
	if l.IPv4AddressState != "" {
		fmt.Printf("%v %v\n", color.HiBlueString("IPv4Address State:"), l.IPv4AddressState)
	}
	if l.IPv6AddressState != "" {
		fmt.Printf("%v %v\n", color.HiBlueString("IPv6Address State:"), l.IPv6AddressState)
	}
}

func displayOneLink(l *link.LinkInfo) {
	if l.HardwareAddr != "" {
		fmt.Printf("       %v %v\n", color.HiBlueString("HW Address:"), l.HardwareAddr)
	}
	fmt.Printf("              %v %v\n", color.HiBlueString("MTU:"), l.Mtu)
	fmt.Printf("        %v %v\n", color.HiBlueString("OperState:"), l.OperState)
	fmt.Printf("            %v %v\n", color.HiBlueString("Flags:"), l.Flags)
}

func displayOneLinkAddresses(addInfo *address.AddressInfo) {
	fmt.Printf("        %v", color.HiBlueString("Addresses:"))
	for _, a := range addInfo.Addresses {
		fmt.Printf(" %v/%v", a.IP, a.Mask)
	}
	fmt.Printf("\n")
}

func displayOneLinkRoutes(ifIndex int, linkRoutes []route.RouteInfo) {
	gws := share.NewSet()
	for _, rt := range linkRoutes {
		if rt.LinkIndex == ifIndex && rt.Gw != "" {
			gws.Add(rt.Gw)
		}
	}

	if gws.Length() > 0 {
		fmt.Printf("          %v %v\n", color.HiBlueString("Gateway:"), strings.Join(gws.Values(), " "))
	}
}

func displayOneLinkDnsAndDomains(link string, dns []resolved.Dns, domains []resolved.Domains) {
	dnsServers := share.NewSet()
	for _, d := range dns {
		if d.Link == link {
			dnsServers.Add(d.Dns)
		}
	}

	if dnsServers.Length() > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("DNS:"), strings.Join(dnsServers.Values(), " "))
	}

	domain := share.NewSet()
	for _, d := range domains {
		if d.Link == link {
			domain.Add(d.Domain)
		}
	}

	if domain.Length() > 0 {
		fmt.Printf("           %v %v\n", color.HiBlueString("Domains:"), strings.Join(dnsServers.Values(), " "))
	}
}

func displayDnsAndDomains(n *resolved.Describe) {
	fmt.Printf("%v\n\n", color.HiBlueString("Global"))
	if !validator.IsEmpty(n.CurrentDNS) {
		fmt.Printf("%v %v\n", color.HiBlueString("CurrentDNS: "), n.CurrentDNS)
	}

	fmt.Printf("%v", color.HiBlueString("        DNS: "))
	for _, d := range n.DnsServers {
		if validator.IsEmpty(d.Link) {
			fmt.Printf("%v ", d.Dns)
		}
	}
	fmt.Printf("\n%v", color.HiBlueString("DNS Domains: "))
	for _, d := range n.Domains {
		fmt.Printf("%v ", d.Domain)
	}

	type linkDns struct {
		Index int32
		Link  string
		Dns   []string
	}

	l := linkDns{}
	dns := make(map[int32]*linkDns)
	for _, d := range n.DnsServers {
		if !validator.IsEmpty(d.Link) {
			if dns[d.Index] != nil {
				l := dns[d.Index]
				l.Dns = append(l.Dns, d.Dns)
			} else {
				dns[d.Index] = &linkDns{
					Index: d.Index,
					Link:  d.Link,
					Dns:   append(l.Dns, d.Dns),
				}
			}
		}
	}

	for _, d := range dns {
		fmt.Printf("\n%v %v (%v)\n", color.HiBlueString("Link"), d.Index, d.Link)
		for _, c := range n.LinkCurrentDNS {
			if c.Index == d.Index {
				fmt.Printf("%v %v\n", color.HiBlueString("Current DNS Server: "), c.Dns)
			}
		}
		fmt.Printf("       %v %v\n", color.HiBlueString("DNS Servers: "), strings.Join(d.Dns, " "))
	}
}

func displayOneLinkNTP(link string, ntp *timesyncd.Describe) {
	if len(ntp.LinkNTPServers) > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("NTP:"), ntp.LinkNTPServers)
	}
}

func displayNetworkStatus(ifName string, network *network.Describe) {
	for _, link := range network.Links {
		if ifName != "" && link.Name != ifName {
			continue
		}

		for _, l := range network.LinksDescribe.Interfaces {
			if link.Name == l.Name {
				displayOneLinkNetworkStatus(&l)
			}
		}

		displayOneLink(&link)

		for _, l := range network.Addresses {
			if l.Name == link.Name {
				displayOneLinkAddresses(&l)
			}
		}

		displayOneLinkRoutes(link.Index, network.Routes)

		if link.Name != "lo" {
			if len(network.Dns) > 0 {
				displayOneLinkDnsAndDomains(link.Name, network.Dns, network.Domains)
			}
		}

		fmt.Printf("\n")
	}
}

func acquireNetworkDescribe(host string, token map[string]string) (*network.Describe, error) {
	resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/network/describe", token, nil)
	if err != nil {
		fmt.Printf("Failed to acquire network info: %v\n", err)
		return nil, err
	}

	n := NetworkDescribe{}
	if err := json.Unmarshal(resp, &n); err != nil {
		fmt.Printf("Failed to decode link json message: %v\n", err)
		return nil, err
	}

	if n.Success {
		return &n.Message, nil
	}

	return nil, errors.New(n.Errors)
}

func acquireResolveDescribe(host string, token map[string]string) error {
	resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/network/resolved/describe", token, nil)
	if err != nil {
		fmt.Printf("Failed to acquire resolve info: %v\n", err)
		return err
	}

	n := ResolveDescribe{}
	if err := json.Unmarshal(resp, &n); err != nil {
		fmt.Printf("Failed to decode link json message: %v\n", err)
		return err
	}

	if n.Success {
		displayDnsAndDomains(&n.Message)
	}

	return nil
}

func acquireNetworkStatus(cmd string, host string, ifName string, token map[string]string) {
	switch cmd {
	case "network":
		n, err := acquireNetworkDescribe(host, token)
		if err != nil {
			fmt.Printf("Failed to fetch network status: %v\n", err)
			return
		}

		displayNetworkStatus(ifName, n)

	case "iostat":
		resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/proc/netdeviocounters", token, nil)
		if err != nil {
			fmt.Printf("Failed to fetch networks device's iostat: %v\n", err)
			return
		}

		n := NetDevIOCounters{}
		if err := json.Unmarshal(resp, &n); err != nil {
			fmt.Printf("Failed to decode json message: %v\n", err)
			return
		}

		if n.Success {
			displayNetDevIOStatistics(&n)
		}
	case "interfaces":
		resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/proc/interfaces", token, nil)
		if err != nil {
			fmt.Printf("Failed to fetch networks devices: %v\n", err)
			return
		}

		n := Interface{}
		if err := json.Unmarshal(resp, &n); err != nil {
			fmt.Printf("Failed to decode json message: %v\n", err)
			return
		}

		if n.Success {
			displayInterfaces(&n)
		}
	}
}

func networkConfigure(network *networkd.Network, host string, token map[string]string) {
	var resp []byte
	var err error

	resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/network/configure", token, *network)
	if err != nil {
		fmt.Printf("Failed to configure network: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to configure network: %v\n", m.Errors)
	}
}

func networkConfigureDHCP(link string, dhcp string, host string, token map[string]string) {
	n := networkd.Network{
		Link: link,
		NetworkSection: networkd.NetworkSection{
			DHCP: dhcp,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureLinkLocalAddressing(link string, linkLocalAddr string, host string, token map[string]string) {
	if !validator.IsLinkLocalAddressing(linkLocalAddr) {
		fmt.Printf("Invalid LinkLocalAddressing: %s\n", linkLocalAddr)
		return
	}

	n := networkd.Network{
		Link: link,
		NetworkSection: networkd.NetworkSection{
			LinkLocalAddressing: linkLocalAddr,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureMulticastDNS(link string, mcastDns string, host string, token map[string]string) {
	if !validator.IsMulticastDNS(mcastDns) {
		fmt.Printf("Invalid MulticastDNS: %s\n", mcastDns)
		return
	}

	n := networkd.Network{
		Link: link,
		NetworkSection: networkd.NetworkSection{
			MulticastDNS: mcastDns,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureRoute(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	link := ""

	r := networkd.RouteSection{}
	for i := 0; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			link = argStrings[i+1]
		case "gw":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Failed to parse gw='%s'\n", argStrings[i+1])
				return
			}
			r.Gateway = argStrings[i+1]
		case "gwonlink":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse gwonlink='%s'\n", argStrings[i+1])
				return
			}
			r.GatewayOnlink = argStrings[i+1]
		case "dest":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Failed to parse dest='%s'\n", argStrings[i+1])
				return
			}
			r.Destination = argStrings[i+1]
		case "src":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Failed to parse src='%s'\n", argStrings[i+1])
				return
			}
			r.Source = argStrings[i+1]
		case "prefsrc":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Failed to parse prefsrc='%s'\n", argStrings[i+1])
				return
			}
			r.PreferredSource = argStrings[i+1]
		case "table":
			if !govalidator.IsInt(argStrings[i+1]) {
				fmt.Printf("Failed to parse table='%s'\n", argStrings[i+1])
				return
			}
			r.Table = argStrings[i+1]
		case "scope":
			if !validator.IsScope(argStrings[i+1]) {
				fmt.Printf("Failed to parse scope='%s'\n", argStrings[i+1])
				return
			}
			r.Scope = argStrings[i+1]
		}

		i++
	}

	n := networkd.Network{
		Link: link,
		RouteSections: []networkd.RouteSection{
			r,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureDHCP4ClientIdentifier(link string, identifier string, host string, token map[string]string) {
	if !validator.IsClientIdentifier(identifier) {
		fmt.Printf("Invalid DHCP4 Client Identifier: %s\n", identifier)
		return
	}

	n := networkd.Network{
		Link: link,
		DHCPv4Section: networkd.DHCPv4Section{
			ClientIdentifier: identifier,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureDHCPIAID(link string, iaid string, host string, token map[string]string) {
	n := networkd.Network{
		Link: link,
		DHCPv4Section: networkd.DHCPv4Section{
			IAID: iaid,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureMTU(link string, mtu string, host string, token map[string]string) {
	n := networkd.Network{
		Link: link,
		LinkSection: networkd.LinkSection{
			MTUBytes: mtu,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureMAC(link string, mac string, host string, token map[string]string) {
	n := networkd.Network{
		Link: link,
		LinkSection: networkd.LinkSection{
			MACAddress: mac,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureLinkGroup(link string, group string, host string, token map[string]string) {
	if !validator.IsEmpty(group) {
		if !validator.IsLinkGroup(group) {
			fmt.Printf("Failed to parse group: Invalid Group=%s\n", group)
			return
		}
	}

	n := networkd.Network{
		Link: link,
		LinkSection: networkd.LinkSection{
			Group: group,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureLinkRequiredFamilyForOnline(link string, rfonline string, host string, token map[string]string) {
	if !validator.IsEmpty(rfonline) {
		if !validator.IsAddressFamily(rfonline) {
			fmt.Printf("Failed to parse online family='%s'\n", rfonline)
			return
		}
	}

	n := networkd.Network{
		Link: link,
		LinkSection: networkd.LinkSection{
			RequiredFamilyForOnline: rfonline,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureLinkActivationPolicy(link string, policy string, host string, token map[string]string) {
	if !validator.IsEmpty(policy) {
		if !validator.IsLinkActivationPolicy(policy) {
			fmt.Printf("Failed to parse activation policy='%s'\n", policy)
			return
		}
	}

	n := networkd.Network{
		Link: link,
		LinkSection: networkd.LinkSection{
			ActivationPolicy: policy,
		},
	}

	networkConfigure(&n, host, token)
}

func networkConfigureMode(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	n := networkd.Network{}
	for i := 0; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		case "arp":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse arp='%s'\n", argStrings[i+1])
				return
			}
			n.LinkSection.ARP = validator.BoolToString(argStrings[i+1])
		case "mc":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse mc='%s'\n", argStrings[i+1])
				return
			}
			n.LinkSection.Multicast = validator.BoolToString(argStrings[i+1])
		case "amc":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse amc='%s'\n", argStrings[i+1])
				return
			}
			n.LinkSection.AllMulticast = validator.BoolToString(argStrings[i+1])
		case "pcs":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse pcs='%s'\n", argStrings[i+1])
				return
			}
			n.LinkSection.Promiscuous = validator.BoolToString(argStrings[i+1])
		case "rfo":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Failed to parse rfo='%s'\n", argStrings[i+1])
				return
			}
			n.LinkSection.RequiredForOnline = validator.BoolToString(argStrings[i+1])
		}

		i++
	}

	networkConfigure(&n, host, token)
}

func networkConfigureAddress(link string, args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	a := networkd.AddressSection{}
	for i := 1; i < args.Len()-1; {
		switch argStrings[i] {
		case "address":
			a.Address = argStrings[i+1]
			if !validator.IsIP(a.Address) {
				fmt.Printf("Invalid IP address: %v\n", a.Address)
				return
			}
		case "peer":
			a.Peer = argStrings[i+1]
			if !validator.IsIP(a.Peer) {
				fmt.Printf("Invalid Peer IP address: %v\n", a.Peer)
				return
			}
		case "label":
			a.Label = argStrings[i+1]
		case "scope":
			a.Scope = argStrings[i+1]
			if !validator.IsScope(a.Scope) {
				fmt.Printf("Invalid scope: %s", a.Scope)
				return
			}
		default:
		}
		i++
	}
	n := networkd.Network{
		Link: link,
		AddressSections: []networkd.AddressSection{
			a,
		},
	}
	networkConfigure(&n, host, token)
}

func networkAddRoutingPolicyRule(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	link := ""

	r := networkd.RoutingPolicyRuleSection{}
	for i := 0; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			link = argStrings[i+1]
		case "tos":
			if !validator.IsRoutingTypeOfService(argStrings[i+1]) {
				fmt.Printf("Invalid tos=%s\n", argStrings[i+1])
				return
			}
			r.TypeOfService = argStrings[i+1]
		case "from":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Invalid from=%s\n", argStrings[i+1])
				return
			}
			r.From = argStrings[i+1]
		case "to":
			if !validator.IsIP(argStrings[i+1]) {
				fmt.Printf("Invalid to=%s\n", argStrings[i+1])
				return
			}
			r.To = argStrings[i+1]
		case "fwmark":
			if !validator.IsRoutingFirewallMark(argStrings[i+1]) {
				fmt.Printf("Invalid fwmark=%s\n", argStrings[i+1])
				return
			}
			r.FirewallMark = argStrings[i+1]
		case "table":
			if !validator.IsRoutingTable(argStrings[i+1]) {
				fmt.Printf("Invalid table=%s\n", argStrings[i+1])
				return
			}
			r.Table = argStrings[i+1]
		case "prio":
			if !validator.IsRoutingPriority(argStrings[i+1]) {
				fmt.Printf("Invalid prio=%s\n", argStrings[i+1])
				return
			}
			r.Priority = argStrings[i+1]
		case "iif":
			if validator.IsEmpty(argStrings[i+1]) {
				fmt.Printf("Invalid iif=%s\n", argStrings[i+1])
				return
			}
			r.IncomingInterface = argStrings[i+1]
		case "oif":
			if validator.IsEmpty(argStrings[i+1]) {
				fmt.Printf("Invalid oif=%s\n", argStrings[i+1])
				return
			}
			r.OutgoingInterface = argStrings[i+1]
		case "srcport":
			if !validator.IsRoutingPort(argStrings[i+1]) {
				fmt.Printf("Invalid srcport=%s\n", argStrings[i+1])
				return
			}
			r.SourcePort = argStrings[i+1]
		case "destport":
			if !validator.IsRoutingPort(argStrings[i+1]) {
				fmt.Printf("Invalid destport=%s\n", argStrings[i+1])
				return
			}
			r.DestinationPort = argStrings[i+1]
		case "ipproto":
			if !validator.IsRoutingIPProtocol(argStrings[i+1]) {
				fmt.Printf("Invalid ipproto=%s\n", argStrings[i+1])
				return
			}
			r.IPProtocol = argStrings[i+1]
		case "invertrule":
			if !validator.IsBool(argStrings[i+1]) {
				fmt.Printf("Invalid invertrule=%s\n", argStrings[i+1])
				return
			}
			r.InvertRule = validator.BoolToString(argStrings[i+1])
		case "family":
			if !validator.IsAddressFamily(argStrings[i+1]) {
				fmt.Printf("Invalid family=%s\n", argStrings[i+1])
				return
			}
			r.Family = argStrings[i+1]
		case "usr":
			if !validator.IsRoutingUser(argStrings[i+1]) {
				fmt.Printf("Invalid usr=%s\n", argStrings[i+1])
				return
			}
			r.User = argStrings[i+1]
		case "suppressprefixlen":
			if !validator.IsRoutingSuppressPrefixLength(argStrings[i+1]) {
				fmt.Printf("Invalid suppressprefixlen=%s\n", argStrings[i+1])
				return
			}
			r.SuppressPrefixLength = argStrings[i+1]
		case "suppressifgrp":
			if !validator.IsRoutingSuppressInterfaceGroup(argStrings[i+1]) {
				fmt.Printf("Invalid suppressifgrp=%s\n", argStrings[i+1])
				return
			}
			r.SuppressInterfaceGroup = argStrings[i+1]
		case "type":
			if !validator.IsRoutingType(argStrings[i+1]) {
				fmt.Printf("Invalid type=%s\n", argStrings[i+1])
				return
			}
			r.Type = argStrings[i+1]
		}

		i++
	}

	n := networkd.Network{
		Link: link,
		RoutingPolicyRuleSections: []networkd.RoutingPolicyRuleSection{
			r,
		},
	}
	networkConfigure(&n, host, token)
}

func networkAddDns(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var dns []string
	for i, args := range argStrings {
		switch args {
		case "dev":
			dev = argStrings[i+1]
		case "dns":
			dns = strings.Split(argStrings[i+1], ",")
		}
	}

	if validator.IsArrayEmpty(dns) {
		fmt.Printf("Failed to add dns. Missing dns server\n")
		return
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := resolved.GlobalDns{
			DnsServers: dns,
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/resolved/add", token, n)
		if err != nil {
			fmt.Printf("Failed to add global Dns server: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				DNS: dns,
			},
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/network/configure", token, n)
		if err != nil {
			fmt.Printf("Failed to add link Dns server: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to add Dns server: %v\n", m.Errors)
	}
}

func networkRemoveDns(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var dns []string
	for i, args := range argStrings {
		switch args {
		case "dev":
			dev = argStrings[i+1]
		case "dns":
			dns = strings.Split(argStrings[i+1], ",")
		}
	}

	if validator.IsArrayEmpty(dns) {
		fmt.Printf("Failed to remove dns. Missing dns server\n")
		return
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := resolved.GlobalDns{
			DnsServers: dns,
		}
		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/resolved/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove global Dns server: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				DNS: dns,
			},
		}

		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/networkd/network/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove link Dns server: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to remove Dns server: %v\n", m.Errors)
	}
}

func networkAddDomains(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var domains []string
	for i, args := range argStrings {
		switch args {
		case "dev":
			dev = argStrings[i+1]
		case "domains":
			domains = strings.Split(argStrings[i+1], ",")
		}
	}

	if validator.IsArrayEmpty(domains) {
		fmt.Printf("Failed to add domains. Missing domains\n")
		return
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := resolved.GlobalDns{
			Domains: domains,
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/resolved/add", token, n)
		if err != nil {
			fmt.Printf("Failed to add global domains: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				Domains: domains,
			},
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/network/configure", token, n)
		if err != nil {
			fmt.Printf("Failed to add link  domains: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to add domains: %v\n", m.Errors)
	}
}

func networkRemoveDomains(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var domains []string
	for i, args := range argStrings {
		switch args {
		case "dev":
			dev = argStrings[i+1]
		case "domains":
			domains = strings.Split(argStrings[i+1], ",")
		}
	}

	if validator.IsArrayEmpty(domains) {
		fmt.Printf("Failed to remove domains. Missing domains server\n")
		return
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := resolved.GlobalDns{
			Domains: domains,
		}
		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/resolved/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove global Dns server: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				Domains: domains,
			},
		}

		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/networkd/network/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove link domains: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to remove domains: %v\n", m.Errors)
	}
}

func networkAddNTP(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var ntp []string
	for i := range argStrings {
		switch argStrings[i] {
		case "dev":
			dev = argStrings[i+1]
		case "ntp":
			ntp = strings.Split(argStrings[i+1], ",")
		}
		i++
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := timesyncd.NTP{
			NTPServers: ntp,
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/timesyncd/add", token, n)
		if err != nil {
			fmt.Printf("Failed to add global NTP server: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				NTP: ntp,
			},
		}
		resp, err = web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/network/configure", token, n)
		if err != nil {
			fmt.Printf("Failed to add link NTP server: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to add NTP server: %v\n", m.Errors)
	}
}

func networkRemoveNTP(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()

	var dev string
	var ntp []string
	for i := range argStrings {
		switch argStrings[i] {
		case "dev":
			dev = argStrings[i+1]
		case "ntp":
			ntp = strings.Split(argStrings[i+1], ",")
		}
		i++
	}

	var resp []byte
	var err error
	if validator.IsEmpty(dev) {
		n := timesyncd.NTP{
			NTPServers: ntp,
		}
		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/timesyncd/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove global NTP server: %v\n", err)
			return
		}
	} else {
		n := networkd.Network{
			Link: dev,
			NetworkSection: networkd.NetworkSection{
				NTP: ntp,
			},
		}
		resp, err = web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/networkd/network/remove", token, n)
		if err != nil {
			fmt.Printf("Failed to remove link NTP server: %v\n", err)
			return
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to remove NTP server: %v\n", m.Errors)
	}
}

func networkConfigureIPv6AcceptRA(link string, ipv6ara string, host string, token map[string]string) {
	if !validator.IsBool(ipv6ara) {
		fmt.Printf("Invalid IPv6AcceptRA: %s\n", ipv6ara)
		return
	}

	n := networkd.Network{
		Link: link,
		NetworkSection: networkd.NetworkSection{
			IPv6AcceptRA: ipv6ara,
		},
	}

	networkConfigure(&n, host, token)
}
