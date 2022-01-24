// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func displayOneLinkNTP(link string, ntp *timesyncd.NTPServer) {
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
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/describe", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/describe", nil)
	}

	if err != nil {
		fmt.Printf("Failed to network info: %v\n", err)
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

func acquireNetworkStatus(cmd string, host string, ifName string, token map[string]string) {
	var resp []byte
	var err error

	switch cmd {
	case "network":
		n, err := acquireNetworkDescribe(host, token)
		if err != nil {
			fmt.Printf("Failed to fetch network status: %v\n", err)
			return
		}

		displayNetworkStatus(ifName, n)

	case "iostat":
		if host != "" {
			resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/proc/netdeviocounters", token, nil)
		} else {
			resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/proc/netdeviocounters", nil)
		}

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
		if host != "" {
			resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/proc/interfaces", token, nil)
		} else {
			resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/proc/interfaces", nil)
		}

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

func networkConfigureDHCP(link string, dhcp string, host string, token map[string]string) {
	if !validator.IsDHCP(dhcp) {
		fmt.Printf("Failed to parse DHCP: %s\n", dhcp)
		return
	}

	n := networkd.Network{
		Link: link,
		NetworkSection: networkd.NetworkSection{
			DHCP: dhcp,
		},
	}

	var resp []byte
	var err error
	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/network/networkd/network/configure", token, n)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/network/networkd/network/configure", n)
	}
	if err != nil {
		fmt.Printf("Failed to configure DHCP: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to configure DHCP: %v\n", m.Errors)
	}
}

func networkCreateVLan(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "vlan",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		case "id":
			if validator.IsVLanId(argStrings[i+1]) {
				n.VLanSection.Id = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse VLan Id: %s\n", argStrings[i+1])
				return
			}
		}

		i++
	}

	if validator.IsEmpty(n.Link) || validator.IsEmpty(n.VLanSection.Id) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create VLan. Missing VLAN name, dev or id\n")
		return
	}

	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/network/networkd/network/configure", token, n)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/network/networkd/netdev/configure", n)
	}
	if err != nil {
		fmt.Printf("Failed to create VLan: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create VLan: %v\n", m.Errors)
	}
}
