// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/pm-web/pkg/systemd"
	"github.com/pm-web/pkg/web"
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

func PerformSystemdUnitCommand(command string, unit string) {
	action := systemd.UnitAction{
		Action: command,
		Unit:   unit,
	}
	resp, err := web.DispatchUnixDomainSocket("POST", "http://localhost/api/v1/service/systemd", action)
	if err != nil {
		fmt.Printf("Failed to execute '%s': %v", command, err)
		os.Exit(1)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v", err)
		os.Exit(0)
	}

	if m.Success {
		fmt.Printf("Command executed successfully")
	} else {
		var s string
		if m.Message== "" {
			s = "n/a"
		}
		
		fmt.Printf("Failed to execute command: %v(%v) ", s, m.Errors)
		os.Exit(1)
	}
}

func fetchSystemdUnitStatus(unit string) {
	resp, err := web.FetchUnixDomainSocket("http://localhost/api/v1/service/systemd/" + unit + "/status")
	if err != nil {
		fmt.Printf("Failed to fetch status: %v", err)
		return
	}

	u := UnitStatus{}
	if err := json.Unmarshal(resp, &u); err != nil {
		fmt.Printf("Failed to decode json message: %v", err)
		return
	}

	if u.Success {
		fmt.Printf("                Name: %+v \n", u.Message.Name)
		fmt.Printf("         Description: %+v \n", u.Message.Description)
		fmt.Printf("             MainPid: %+v \n", u.Message.MainPid)
		fmt.Printf("           LoadState: %+v \n", u.Message.LoadState)
		fmt.Printf("         ActiveState: %+v \n", u.Message.ActiveState)
		fmt.Printf("            SubState: %+v \n", u.Message.SubState)
		fmt.Printf("       UnitFileState: %+v \n", u.Message.UnitFileState)

		t := time.Unix(int64(u.Message.StateChangeTimestamp), 0)
		fmt.Printf("StateChangeTimeStamp: %+v \n", t.Format(time.UnixDate))

		switch u.Message.ActiveState {
		case "active", "reloading":

			t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
			fmt.Printf("              Active: %s (%s) since %v", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))

		case "inactive", "failed":

			t := time.Unix(int64(u.Message.ActiveEnterTimestamp), 0)
			fmt.Printf("             Active: %s (%s) since %v", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))

		case "activating":
			var t time.Time

			if u.Message.ActiveExitTimestamp != 0 {
				t = time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			} else {
				t = time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			}

			fmt.Printf("             Active: %s (%s) %v", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))

		default:
			t := time.Unix(int64(u.Message.ActiveExitTimestamp), 0)
			fmt.Printf("             Active: %s (%s) ago %v", u.Message.ActiveState, u.Message.SubState, t.Format(time.UnixDate))
		}
	} else {
		fmt.Println(u.Errors)
	}
}

func main() {
	log.SetOutput(ioutil.Discard)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version=%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:    "pmctl",
		Version: "v0.1",
		Usage:   "Introspects and manage system",
	}

	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "service",
			Aliases: []string{"s"},
			Usage:   "Perform systemd service commands",
			Subcommands: []*cli.Command{
				{
					Name:  "status",
					Usage: "Display systemd service status",
					Action: func(c *cli.Context) error {
						fetchSystemdUnitStatus(c.Args().First())
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Start systemd service",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("start", c.Args().First())
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stop systemd service",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("stop", c.Args().First())
						return nil
					},
				},
				{
					Name:  "restart",
					Usage: "Restart systemd service",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("restart", c.Args().First())
						return nil
					},
				},
				{
					Name:  "mask",
					Usage: "Mask systemd unit",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("mask", c.Args().First())
						return nil
					},
				},
				{
					Name:  "unmask",
					Usage: "Unmask systemd unit",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("unmask", c.Args().First())
						return nil
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
