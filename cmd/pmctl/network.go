package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/network"
	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/link"
	"github.com/pm-web/plugins/network/netlink/route"
	"github.com/pm-web/plugins/network/networkd"
	"github.com/pm-web/plugins/network/resolved"
	"github.com/pm-web/plugins/network/timesyncd"
	"github.com/shirou/gopsutil/v3/net"
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

type NTP struct {
	Success bool                `json:"success"`
	Message timesyncd.NTPServer `json:"message"`
	Errors  string              `json:"errors"`
}

func displayInterfaces(i *Interface) {
	for _, n := range i.Message {
		fmt.Printf("            Name: %v\n", n.Name)
		fmt.Printf("           Index: %v\n", n.Index)
		fmt.Printf("             MTU: %v\n", n.MTU)

		fmt.Printf("           Flags: ")
		for _, j := range n.Flags {
			fmt.Printf("%v ", j)
		}
		fmt.Printf("\n")

		fmt.Printf("Hardware Address: %v\n", n.HardwareAddr)

		fmt.Printf("       Addresses: ")
		for _, j := range n.Addrs {
			fmt.Printf("%v ", j.Addr)
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

func displayOneLinkDNS(link string, dns []resolved.Dns) {
	dnsServers := share.NewSet()
	for _, d := range dns {
		if d.Link == link {
			dnsServers.Add(d.Dns)
		}
	}

	if dnsServers.Length() > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("DNS:"), strings.Join(dnsServers.Values(), " "))
	}
}

func displayOneLinkNTP(link string, ntp *timesyncd.NTPServer) {
	if len(ntp.LinkNTPServers) > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("NTP:"), ntp.LinkNTPServers)
	}
}

func displayNetworkStatus(ifName string, network *network.Describe, ntp *timesyncd.NTPServer) {
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
				displayOneLinkDNS(link.Name, network.Dns)
			}

			if ntp != nil {
				displayOneLinkNTP(link.Name, ntp)
			}
		}

		fmt.Printf("\n")
	}
}

func acquireLNetwork(host string, token map[string]string) (*network.Describe, error) {
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

func acquireNTP(host string, token map[string]string) (*timesyncd.NTPServer, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/timesyncd/linkntpserver", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/timesyncd/linkntpserver", nil)
	}

	if err != nil {
		fmt.Printf("Failed to fetch NTP: %v\n", err)
		return nil, err
	}

	rt := NTP{}
	if err := json.Unmarshal(resp, &rt); err != nil {
		fmt.Printf("Failed to decode NTP json message: %v\n", err)
		return nil, err
	}

	if rt.Success {
		return &rt.Message, nil
	}

	return nil, errors.New(rt.Errors)
}

func acquireNetworkStatus(cmd string, host string, ifName string, token map[string]string) {
	var resp []byte
	var err error

	switch cmd {
	case "network":

		n, err := acquireLNetwork(host, token)
		if err != nil {
			fmt.Printf("Failed to fetch network status: %v\n", err)
			return
		}

		ntp, _ := acquireNTP(host, token)
		displayNetworkStatus(ifName, n, ntp)

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
