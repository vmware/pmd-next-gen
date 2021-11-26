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

type LinkStatus struct {
	Success bool `json:"success"`
	Message struct {
		Interfaces []networkd.LinkDescribe `json:"Interfaces"`
	} `json:"message"`
	Errors string `json:"errors"`
}

type Links struct {
	Success bool            `json:"success"`
	Message []link.LinkInfo `json:"message"`
	Errors  string          `json:"errors"`
}

type Addresses struct {
	Success bool                  `json:"success"`
	Message []address.AddressInfo `json:"message"`
	Errors  string                `json:"errors"`
}

type Routes struct {
	Success bool              `json:"success"`
	Message []route.RouteInfo `json:"message"`
	Errors  string            `json:"errors"`
}

type DNS struct {
	Success bool           `json:"success"`
	Message []resolved.DNS `json:"message"`
	Errors  string         `json:"errors"`
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

func acquireLinks(host string, token map[string]string) ([]link.LinkInfo, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/netlink/link", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/netlink/link", nil)
	}

	if err != nil {
		fmt.Printf("Failed to acquire links: %v\n", err)
		return nil, err
	}

	a := Links{}
	if err := json.Unmarshal(resp, &a); err != nil {
		fmt.Printf("Failed to decode link json message: %v\n", err)
		return nil, err
	}

	if a.Success {
		return a.Message, nil
	}

	return nil, errors.New(a.Errors)
}

func acquireLinkAddresses(host string, token map[string]string) ([]address.AddressInfo, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/netlink/address", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/netlink/address", nil)
	}

	if err != nil {
		fmt.Printf("Failed to acquire addresses: %v\n", err)
		return nil, err
	}

	a := Addresses{}
	if err := json.Unmarshal(resp, &a); err != nil {
		fmt.Printf("Failed to decode address json message: %v\n", err)
		return nil, err
	}

	if a.Success {
		return a.Message, nil
	}

	return nil, errors.New(a.Errors)
}

func acquireLinkRoutes(host string, token map[string]string) ([]route.RouteInfo, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/netlink/route", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/netlink/route", nil)
	}

	if err != nil {
		fmt.Printf("Failed to fetch routes: %v\n", err)
		return nil, err
	}

	rt := Routes{}
	if err := json.Unmarshal(resp, &rt); err != nil {
		fmt.Printf("Failed to decode route json message: %v\n", err)
		return nil, err
	}

	if rt.Success {
		return rt.Message, nil
	}

	return nil, errors.New(rt.Errors)
}

func acquireDNS(host string, token map[string]string) ([]resolved.DNS, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/resolved/dns", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/resolved/dns", nil)
	}

	if err != nil {
		fmt.Printf("Failed to fetch dns: %v\n", err)
		return nil, err
	}

	rt := DNS{}
	if err := json.Unmarshal(resp, &rt); err != nil {
		fmt.Printf("Failed to decode DNS json message: %v\n", err)
		return nil, err
	}

	if rt.Success {
		return rt.Message, nil
	}

	return nil, errors.New(rt.Errors)
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

func displayOneLinkDNS(link string, dns []resolved.DNS) {
	dnsServers := share.NewSet()
	for _, d := range dns {
		if d.Link == link {
			dnsServers.Add(d.DNS)
		}
	}

	if dnsServers.Length() > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("DNS:"), strings.Join(dnsServers.Values(), " "))
	}
}

func displayCurrentNTP(link string, ntp *timesyncd.NTPServer) {
	if ntp.ServerAddress != "" && ntp.ServerName != "" {
		fmt.Printf("              %v %v (%v)\n", color.HiBlueString("NTP:"), ntp.ServerName, ntp.ServerAddress)
	}
}

func displayOneLinkNTP(link string, ntp *timesyncd.NTPServer) {
	if len(ntp.LinkNTPServers) > 0 {
		fmt.Printf("              %v %v\n", color.HiBlueString("NTP:"), ntp.LinkNTPServers)
	}
}

func displayNetworkStatus(l *LinkStatus, link string, links []link.LinkInfo, linkAddresses []address.AddressInfo, linkRoutes []route.RouteInfo, dns []resolved.DNS, ntp *timesyncd.NTPServer) {
	for _, n := range l.Message.Interfaces {
		if link != "" && link != n.Name {
			continue
		}

		displayOneLinkNetworkStatus(&n)
		for _, k := range links {
			if k.Name == n.Name {
				displayOneLink(&k)
			}
		}
		for _, k := range linkAddresses {
			if k.Name == n.Name {
				displayOneLinkAddresses(&k)

			}
		}
		displayOneLinkRoutes(n.Index, linkRoutes)

		if n.Name != "lo" {
			if len(dns) > 0 {
				displayOneLinkDNS(n.Name, dns)
			}

			if ntp != nil {
				displayOneLinkNTP(n.Name, ntp)
			}
		}

		fmt.Printf("\n")
	}
}

func acquireNetworkStatus(cmd string, host string, link string, token map[string]string) {
	var resp []byte
	var err error

	switch cmd {
	case "network":
		if host != "" {
			resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/networkd/network/describe", token, nil)
		} else {
			resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/networkd/network/describe", nil)
		}
		if err != nil {
			fmt.Printf("Failed to fetch network status: %v\n", err)
			return
		}

		n := LinkStatus{}
		if err := json.Unmarshal(resp, &n); err != nil {
			fmt.Printf("Failed to decode json message: %v\n", err)
			return
		}

		if !n.Success {
			fmt.Printf("Failed to fetch links: %v\n", err)
			return
		}

		links, err := acquireLinks(host, token)
		if err != nil {
			return
		}

		addresses, err := acquireLinkAddresses(host, token)
		if err != nil {
			return
		}

		routes, err := acquireLinkRoutes(host, token)
		if err != nil {
			return
		}

		dns, _ := acquireDNS(host, token)
		ntp, _ := acquireNTP(host, token)

		displayNetworkStatus(&n, link, links, addresses, routes, dns, ntp)

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
