// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/network/netlink/address"
	"github.com/pm-web/plugins/network/netlink/route"
	"github.com/pm-web/plugins/network/networkd"
	"github.com/pm-web/plugins/systemd"
)

type UnitStatus struct {
	Success bool               `json:"success"`
	Message systemd.UnitStatus `json:"message"`
	Errors  string             `json:"errors"`
}

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
		Interfaces []networkd.LinkState `json:"Interfaces"`
	} `json:"message"`
	Errors string `json:"errors"`
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

func executeSystemdUnitCommand(command string, unit string, host string, token map[string]string) {
	c := systemd.UnitAction{
		Action: command,
		Unit:   unit,
	}

	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/service/systemd", token, c)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from remote host: %v\n", err)
			return
		}
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/service/systemd", c)
		if err != nil {
			fmt.Printf("Failed to execute '%s': %v\n", command, err)
			os.Exit(1)
		}
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to execute command: %v\n", m.Errors)
		os.Exit(1)
	}
}

func acquireSystemdUnitStatus(unit string, host string, token map[string]string) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/service/systemd/"+unit+"/status", token, nil)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from remote host: %v\n", err)
			return
		}
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/service/systemd/"+unit+"/status", nil)
		if err != nil {
			fmt.Printf("Failed to fetch unit status from unix domain socket: %v\n", err)
			return
		}
	}

	u := UnitStatus{}
	if err := json.Unmarshal(resp, &u); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if u.Success {
		fmt.Printf("                  Name: %+v \n", u.Message.Name)
		fmt.Printf("           Description: %+v \n", u.Message.Description)
		fmt.Printf("               MainPid: %+v \n", u.Message.MainPid)
		fmt.Printf("             LoadState: %+v \n", u.Message.LoadState)
		fmt.Printf("           ActiveState: %+v \n", u.Message.ActiveState)
		fmt.Printf("              SubState: %+v \n", u.Message.SubState)
		fmt.Printf("         UnitFileState: %+v \n", u.Message.UnitFileState)

		if u.Message.StateChangeTimestamp > 0 {
			t := time.Unix(int64(u.Message.StateChangeTimestamp), 0)
			fmt.Printf("  StateChangeTimeStamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("  StateChangeTimeStamp: %+v \n", 0)
		}

		if u.Message.ActiveEnterTimestamp > 0 {
			t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
			fmt.Printf("  ActiveEnterTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("  ActiveEnterTimestamp: %+v \n", 0)
		}

		if u.Message.ActiveEnterTimestamp > 0 {
			t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
			fmt.Printf(" InactiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf(" InactiveExitTimestamp: %+v \n", 0)
		}

		if u.Message.ActiveExitTimestamp > 0 {
			t := time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			fmt.Printf("   ActiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf("   ActiveExitTimestamp: %+v \n", 0)
		}

		if u.Message.InactiveExitTimestamp > 0 {
			t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
			fmt.Printf(" InactiveExitTimestamp: %+v \n", t.Format(time.UnixDate))
		} else {
			fmt.Printf(" InactiveExitTimestamp: %+v \n", 0)
		}

		switch u.Message.ActiveState {
		case "active", "reloading":
			if u.Message.ActiveEnterTimestamp > 0 {
				t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				fmt.Printf("                Active: %s (%s) since %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("                Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		case "inactive", "failed":
			if u.Message.ActiveExitTimestamp != 0 {
				t := time.Unix(int64(u.Message.InactiveExitTimestamp), 0)
				fmt.Printf("                Active: %s (%s) since %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("                Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		case "activating":
			var t time.Time

			if u.Message.ActiveExitTimestamp > 0 || u.Message.ActiveEnterTimestamp > 0 {
				if u.Message.ActiveExitTimestamp > 0 {
					t = time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				} else if u.Message.ActiveEnterTimestamp > 0 {
					t = time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
				}

				fmt.Printf("               Active: %s (%s) %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
			} else {
				fmt.Printf("               Active: %s (%s)\n", u.Message.ActiveState, u.Message.SubState)
			}
		default:
			t := time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			fmt.Printf("               Active: %s (%s) ago %v\n", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
		}
	} else {
		fmt.Println(u.Errors)
	}
}

func displayNetDevIOStatistics(netDev *NetDevIOCounters) {
	for _, n := range netDev.Message {
		fmt.Printf("            Name: %v\n", n.Name)
		fmt.Printf("Packets received: %v\n", n.PacketsRecv)
		fmt.Printf("    Packets sent: %v\n", n.PacketsSent)
		fmt.Printf("  Bytes received: %v\n", n.BytesRecv)
		fmt.Printf("      Bytes sent: %v\n", n.BytesSent)
		fmt.Printf("         Drop in: %v\n", n.Dropin)
		fmt.Printf("        Drop out: %v\n", n.Dropout)
		fmt.Printf("        Error in: %v\n", n.Errin)
		fmt.Printf("       Error out: %v\n", n.Errout)
		fmt.Printf("         Fifo in: %v\n", n.Fifoin)
		fmt.Printf("        Fifo out: %v\n\n", n.Fifoout)
	}
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

func acquireLinkAddresses(host string, token map[string]string) ([]address.AddressInfo, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/netlink/address", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/netlink/address", nil)
	}

	if err != nil {
		fmt.Printf("Failed to fetch addresses: %v\n", err)
		return nil, err
	}

	a := Addresses{}
	if err := json.Unmarshal(resp, &a); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
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
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if rt.Success {
		return rt.Message, nil
	}

	return nil, errors.New(rt.Errors)
}

func displayNetworkStatus(l *LinkStatus, linkAddresses []address.AddressInfo, linkRoutes []route.RouteInfo) {
	for _, n := range l.Message.Interfaces {
		fmt.Printf("             Name: %v\n", n.Name)
		if len(n.AlternativeNames) > 0 {
			fmt.Printf("Alternative Names: %v\n", strings.Join(n.AlternativeNames, " "))
		}
		fmt.Printf("            Index: %v\n", n.Index)
		if n.LinkFile != "" {
			fmt.Printf("        Link File: %v\n", n.LinkFile)
		}
		if n.NetworkFile != "" {
			fmt.Printf("     Network File: %v\n", n.NetworkFile)
		}
		fmt.Printf("             Type: %v\n", n.Type)
		fmt.Printf("            State: %v(%v)\n", n.OperationalState, n.SetupState)
		if n.Driver != "" {
			fmt.Printf("           Driver: %v\n", n.Driver)
		}
		if n.Vendor != "" {
			fmt.Printf("           Vendor: %v\n", n.Vendor)
		}
		if n.Model != "" {
			fmt.Printf("            Model: %v\n", n.Model)
		}
		if n.Path != "" {
			fmt.Printf("             Path: %v\n", n.Path)
		}
		fmt.Printf("    Carrier State: %v\n", n.CarrierState)
		if n.OnlineState != "" {
			fmt.Printf("     Online State: %v\n", n.OnlineState)
		}
		fmt.Printf("IPv4Address State: %v\n", n.IPv4AddressState)
		fmt.Printf("IPv6Address State: %v\n", n.IPv6AddressState)

		for _, k := range linkAddresses {
			if k.Name == n.Name {
				if k.Mac != "" {
					fmt.Printf("       HW Address: %v\n", k.Mac)
				}
				fmt.Printf("              MTU: %v\n", k.MTU)
				fmt.Printf("        Addresses: ")
				for _, j := range k.Addresses {
					fmt.Printf("%v/%v ", j.IP, j.Mask)
				}
				fmt.Printf("\n")
			}
		}
		for _, k := range linkRoutes {
			if k.LinkIndex == n.Index && k.Gw != "" {
				fmt.Printf("          Gateway: %v\n", k.Gw)
				break
			}
		}

		fmt.Printf("\n")
	}
}

func acquireNetworkStatus(cmd string, host string, token map[string]string) {
	var resp []byte
	var err error

	switch cmd {
	case "network":
		if host != "" {
			resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/network/networkd/network/links", token, nil)
		} else {
			resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/network/networkd/network/links", nil)
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
			fmt.Printf("Failed to fetch network status: %v\n", err)
			return
		}

		addresses, err := acquireLinkAddresses(host, token)
		if err != nil {
			fmt.Printf("Failed to fetch link addresses: %v\n", err)
			return
		}

		routes, err := acquireLinkRoutes(host, token)
		if err != nil {
			fmt.Printf("Failed to fetch link routes: %v\n", err)
			return
		}

		displayNetworkStatus(&n, addresses, routes)

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

func main() {
	log.SetOutput(ioutil.Discard)

	token, _ := web.BuildAuthTokenFromEnv()

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version=%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:    "pmctl",
		Version: "v0.1",
		Usage:   "Introspects and controls the system",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "http://localhost:5208",
		},
	}

	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "Control the systemd system and service manager",
			Subcommands: []*cli.Command{
				{
					Name:  "status",
					Usage: "Show terse runtime status information about one unit",

					Action: func(c *cli.Context) error {
						acquireSystemdUnitStatus(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Start (activate) one unit specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("start", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stop (deactivate) one specified on the command line.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("stop", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "restart",
					Usage: "Stop and then start one unit specified on the command line. If the unit is not running yet, it will be started.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("restart", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "mask",
					Usage: "Mask one unit, as specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("mask", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "unmask",
					Usage: "Unmask one unit file, as specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("unmask", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "try-restart",
					Usage: "Stop and then start one unit specified on the command line if the unit are running. This does nothing if unit is not running.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("try-restart", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:  "reload-or-restart",
					Usage: "Reload one unit if they support it. If not, stop and then start instead. If the unit is not running yet, it will be started.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("reload-or-restart", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "status",
			Aliases: []string{"n"},
			Usage:   "Introspects of system or network status",
			Subcommands: []*cli.Command{
				{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Introspects the network status",

					Action: func(c *cli.Context) error {
						acquireNetworkStatus("network", c.String("url"), token)
						return nil
					},
					Subcommands: []*cli.Command{
						{
							Name:  "iostat",
							Usage: "Show iostat of interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("iostat", c.String("url"), token)
								return nil
							},
						},
						{
							Name:  "interfaces",
							Usage: "Show network interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("interfaces", c.String("url"), token)
								return nil
							},
						},
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Failed to run cli: '%+v'", err)
	}
}
