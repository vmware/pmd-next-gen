// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/systemd"
)

type UnitStatus struct {
	Success bool `json:"success"`
	Message struct {
		Property               string `json:"property"`
		Unit                   string `json:"unit"`
		Name                   string `json:"Name"`
		Description            string `json:"Description"`
		MainPid                int    `json:"MainPid"`
		LoadState              string `json:"LoadState"`
		ActiveState            string `json:"ActiveState"`
		SubState               string `json:"SubState"`
		Followed               string `json:"Followed"`
		Path                   string `json:"Path"`
		JobID                  int    `json:"JobId"`
		JobType                string `json:"JobType"`
		JobPath                string `json:"JobPath"`
		UnitFileState          string `json:"UnitFileState"`
		StateChangeTimestamp   int64  `json:"StateChangeTimestamp"`
		InactiveExitTimestamp  int64  `json:"InactiveExitTimestamp"`
		ActiveEnterTimestamp   int64  `json:"ActiveEnterTimestamp"`
		ActiveExitTimestamp    int64  `json:"ActiveExitTimestamp"`
		InactiveEnterTimestamp int64  `json:"InactiveEnterTimestamp"`
	} `json:"message"`
	Errors string `json:"errors"`
}

type NetDevIOCounters struct {
	Success bool `json:"success"`
	Message []struct {
		Name        string `json:"name"`
		BytesSent   int    `json:"bytesSent"`
		BytesRecv   int    `json:"bytesRecv"`
		PacketsSent int    `json:"packetsSent"`
		PacketsRecv int    `json:"packetsRecv"`
		Errin       int    `json:"errin"`
		Errout      int    `json:"errout"`
		Dropin      int    `json:"dropin"`
		Dropout     int    `json:"dropout"`
		Fifoin      int    `json:"fifoin"`
		Fifoout     int    `json:"fifoout"`
	} `json:"message"`
	Errors string `json:"errors"`
}
type Interface struct {
	Success bool `json:"success"`
	Message []struct {
		Index        int      `json:"index"`
		Mtu          int      `json:"mtu"`
		Name         string   `json:"name"`
		HardwareAddr string   `json:"hardwareAddr"`
		Flags        []string `json:"flags"`
		Addrs        []struct {
			Addr string `json:"addr"`
		} `json:"addrs"`
	} `json:"message"`
	Errors string `json:"errors"`
}

type LinkStatus struct {
	Success bool `json:"success"`
	Message struct {
		Interfaces []struct {
			AddressState     string        `json:"AddressState"`
			AlternativeNames []interface{} `json:"AlternativeNames"`
			CarrierState     string        `json:"CarrierState"`
			Driver           interface{}   `json:"Driver"`
			IPv4AddressState string        `json:"IPv4AddressState"`
			IPv6AddressState string        `json:"IPv6AddressState"`
			Index            int           `json:"Index"`
			LinkFile         interface{}   `json:"LinkFile"`
			Model            interface{}   `json:"Model"`
			Name             string        `json:"Name"`
			OnlineState      interface{}   `json:"OnlineState"`
			OperationalState string        `json:"OperationalState"`
			Path             interface{}   `json:"Path"`
			SetupState       string        `json:"SetupState"`
			Type             string        `json:"Type"`
			Vendor           interface{}   `json:"Vendor"`
			NetworkFile      string        `json:"NetworkFile,omitempty"`
		} `json:"Interfaces"`
	} `json:"message"`
	Errors string `json:"errors"`
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
		fmt.Printf("             MTU: %v\n", n.Mtu)

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

func displayNetworkStatus(l *LinkStatus) {
	for _, n := range l.Message.Interfaces {
		fmt.Printf("            Name: %v\n", n.Name)
		fmt.Printf("           Index: %v\n", n.Index)
		fmt.Printf("           NetworkFile: %v\n", n.NetworkFile)

		fmt.Printf("\n\n")
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
			fmt.Printf("Failed to fetch networks status: %v\n", err)
			return
		}

		n := LinkStatus{}
		if err := json.Unmarshal(resp, &n); err != nil {
			fmt.Printf("Failed to decode json message: %v\n", err)
			return
		}

		if n.Success {
			displayNetworkStatus(&n)
			return
		}
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
			return
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
			return
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
