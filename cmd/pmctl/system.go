package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/management/hostname"
	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/route"
	"github.com/pm-web/plugins/network/networkd"
	"github.com/pm-web/plugins/systemd"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Hostname struct {
	Success bool              `json:"success"`
	Message hostname.Describe `json:"message"`
	Errors  string            `json:"errors"`
}

type Systemd struct {
	Success bool             `json:"success"`
	Message systemd.Describe `json:"message"`
	Errors  string           `json:"errors"`
}

type HostInfo struct {
	Success bool          `json:"success"`
	Message host.InfoStat `json:"message"`
	Errors  string        `json:"errors"`
}

type UserStat struct {
	Success bool            `json:"success"`
	Message []host.UserStat `json:"message"`
	Errors  string          `json:"errors"`
}

type VMStat struct {
	Success bool                  `json:"success"`
	Message mem.VirtualMemoryStat `json:"message"`
	Errors  string                `json:"errors"`
}

func acquireHostname(host string, token map[string]string) (*hostname.Describe, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/system/hostname/describe", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/system/hostname/describe", nil)
	}
	if err != nil {
		fmt.Printf("Failed to fetch hostname: %v\n", err)
		return nil, err
	}

	h := Hostname{}
	if err := json.Unmarshal(resp, &h); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !h.Success {
		fmt.Printf("%v\n", h.Errors)
		return nil, errors.New(h.Errors)
	}

	return &h.Message, nil
}

func acquireSystemd(host string, token map[string]string) (*systemd.Describe, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/service/systemd/manager/describe", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/service/systemd/manager/describe", nil)
	}
	if err != nil {
		fmt.Printf("Failed to fetch hostname: %v\n", err)
		return nil, err
	}

	sd := Systemd{}
	if err := json.Unmarshal(resp, &sd); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !sd.Success {
		fmt.Printf("%v\n", sd.Errors)
		return nil, errors.New(sd.Errors)
	}

	return &sd.Message, nil
}

func acquireNetworkState(host string, token map[string]string) (*networkd.NetworkDescribe, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/networkd/network/describestate", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/networkd/network/describestate", nil)
	}
	if err != nil {
		fmt.Printf("Failed to fetch network state: %v\n", err)
		return nil, err
	}

	n := NetworkState{}
	if err := json.Unmarshal(resp, &n); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !n.Success {
		fmt.Printf("%v\n", n.Errors)
		return nil, errors.New(n.Errors)
	}

	return &n.Message, nil
}

func acquireHostInfo(host string, token map[string]string) (*host.InfoStat, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/proc/hostinfo", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/proc/hostinfo", nil)
	}
	if err != nil {
		fmt.Printf("Failed to host information: %v\n", err)
		return nil, err
	}

	h := HostInfo{}
	if err := json.Unmarshal(resp, &h); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !h.Success {
		fmt.Printf("%v\n", h.Errors)
		return nil, errors.New(h.Errors)
	}

	return &h.Message, nil
}

func acquireUserStat(host string, token map[string]string) ([]host.UserStat, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/proc/userstat", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/proc/userstat", nil)
	}
	if err != nil {
		fmt.Printf("Failed to userstat information: %v\n", err)
		return nil, err
	}

	h := UserStat{}
	if err := json.Unmarshal(resp, &h); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !h.Success {
		fmt.Printf("%v\n", h.Errors)
		return nil, errors.New(h.Errors)
	}

	return h.Message, nil
}

func acquireVirtualMemoryStat(host string, token map[string]string) (*mem.VirtualMemoryStat, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/proc/virtualmemory", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/proc/virtualmemory", nil)
	}
	if err != nil {
		fmt.Printf("Failed to userstat information: %v\n", err)
		return nil, err
	}

	h := VMStat{}
	if err := json.Unmarshal(resp, &h); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !h.Success {
		fmt.Printf("%v\n", h.Errors)
		return nil, errors.New(h.Errors)
	}

	return &h.Message, nil
}

func displayHostname(h *hostname.Describe) {
	fmt.Printf("              %v %v\n", color.HiBlueString("System Name:"), h.StaticHostname)
	fmt.Printf("                   %v %v (%v) %v\n", color.HiBlueString("Kernel:"), h.KernelName, h.KernelRelease, h.KernelVersion)
	fmt.Printf("                  %v %v\n", color.HiBlueString("Chassis:"), h.Chassis)
	if h.HardwareModel != "" {
		fmt.Printf("           %v %v\n", color.HiBlueString("Hardware Model:"), h.HardwareModel)
	}
	if h.HardwareVendor != "" {
		fmt.Printf("          %v %v\n", color.HiBlueString("Hardware Vendor:"), h.HardwareVendor)
	}
	if h.ProductUUID != "" {
		fmt.Printf("             %v %v\n", color.HiBlueString("Product UUID:"), h.ProductUUID)
	}
	fmt.Printf("         %v %v\n", color.HiBlueString("Operating System:"), h.OperatingSystemPrettyName)
	if h.OperatingSystemHomeURL != "" {
		fmt.Printf("%v %v\n", color.HiBlueString("Operating System Home URL:"), h.OperatingSystemHomeURL)
	}
}

func displaySystemd(sd *systemd.Describe) {
	fmt.Printf("          %v %v\n", color.HiBlueString("Systemd Version:"), sd.Version)
	fmt.Printf("             %v %v\n", color.HiBlueString("Architecture:"), sd.Architecture)
	fmt.Printf("           %v %v\n", color.HiBlueString("Virtualization:"), sd.Virtualization)
}

func displayNetworkState(n *networkd.NetworkDescribe) {
	fmt.Printf("            %-10v %v (%v)\n", color.HiBlueString("Network State:"), n.OperationalState, n.CarrierState)
	if n.OnlineState != "" {
		fmt.Printf("     %-10v %v\n", color.HiBlueString("Network Online State:"), n.OnlineState)
	}
	if len(n.DNS) > 0 {
		fmt.Printf("                      %-10v %v\n", color.HiBlueString("DNS:"), strings.Join(n.DNS, " "))
	}
	if len(n.Domains) > 0 {
		fmt.Printf("                  %-10v %v\n", color.HiBlueString("Domains:"), strings.Join(n.Domains, " "))
	}
	if len(n.NTP) > 0 {
		fmt.Printf("                      %-10v %v\n", color.HiBlueString("NTP:"), strings.Join(n.NTP, " "))
	}
}

func displayNetworkAddresses(addInfo []address.AddressInfo) {
	fmt.Printf("                  %v", color.HiBlueString("Address:"))

	b := true
	for _, addrs := range addInfo {
		if addrs.Name == "lo" {
			continue
		}
		for _, a := range addrs.Addresses {
			if b {
				fmt.Printf(" %v/%v %v %v\n", a.IP, a.Mask, color.HiGreenString("on link"), addrs.Name)
				b = false
			} else {
				fmt.Printf("                           %v/%v %v %v\n", a.IP, a.Mask, color.HiGreenString("on link"), addrs.Name)
			}
		}
	}
}

func displayRoutes(linkRoutes []route.RouteInfo) {
	fmt.Printf("                  %v", color.HiBlueString("Gateway:"))

	b := true
	gws := share.NewSet()
	for _, rt := range linkRoutes {
		if rt.Gw != "" {
			if b {
				fmt.Printf(" %v %v %v\n", rt.Gw, color.HiGreenString("on link"), rt.LinkName)
				gws.Add(rt.LinkName)
				b = false
			} else {
				if !gws.Contains(rt.LinkName) {
					fmt.Printf("                           %v %v %v\n", rt.Gw, color.HiGreenString("on link"), rt.LinkName)
					gws.Add(rt.LinkName)
				}
			}
		}
	}
}

func displayHostInfo(h *host.InfoStat, u []host.UserStat) {
	t := time.Unix(int64(h.BootTime), 0)
	d, _ := share.SecondsToDuration(h.Uptime)
	fmt.Printf("                   %v %v (%v) %v (%v) %v (%v) %v (%v)\n", color.HiBlueString("Uptime:"), color.HiYellowString("Running Since"), d,
		color.HiYellowString("Booted"), t.Format(time.UnixDate), color.HiYellowString("Users"), len(u), color.HiYellowString("Proc"), h.Procs)
}

func displayVMStat(v *mem.VirtualMemoryStat) {
	fmt.Printf("                   %v %v (%v) %v (%v) %v (%v) %v (%v)\n", color.HiBlueString("Memory:"), color.HiYellowString("Total"), v.Total,
		color.HiYellowString("Used"), v.Total, color.HiYellowString("Free"), v.Free, color.HiYellowString("Available"), v.Available)
}

func acquireSystemStatus(host string, token map[string]string) {
	h, err := acquireHostname(host, token)
	if err != nil {
		return
	}

	sd, err := acquireSystemd(host, token)
	if err != nil {
		return
	}

	n, err := acquireNetworkState(host, token)
	if err != nil {
		return
	}

	addrs, err := acquireLinkAddresses(host, token)
	if err != nil {
		return
	}

	rts, err := acquireLinkRoutes(host, token)
	if err != nil {
		return
	}

	hInfo, err := acquireHostInfo(host, token)
	if err != nil {
		return
	}

	u, err := acquireUserStat(host, token)
	if err != nil {
		return
	}

	v, err := acquireVirtualMemoryStat(host, token)
	if err != nil {
		return
	}

	displayHostname(h)
	displaySystemd(sd)
	displayNetworkState(n)
	displayNetworkAddresses(addrs)
	displayRoutes(rts)
	displayHostInfo(hInfo, u)
	displayVMStat(v)
}

func SetHostname(hostName string, host string, token map[string]string) {
	var resp []byte
	var err error

	h := hostname.Hostname {
		Method: "SetStaticHostname",
		Value: hostName,
	}

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/system/hostname/method", token, h)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/system/hostname/method", h)
	}
	if err != nil {
		fmt.Printf("Failed set hostname: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed set hostname: %v\n", m.Errors)
	}

	fmt.Println(m.Message)
}