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

	"github.com/pm-web/plugins/systemd"
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

func PerformSystemdUnitCommand(command string, unit string, host string) {
	c := systemd.UnitAction{
		Action: command,
		Unit:   unit,
	}

	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/service/systemd", nil, c)
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

func fetchSystemdUnitStatus(unit string, host string) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/service/systemd/"+unit+"/status", nil, nil)
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

func main() {
	log.SetOutput(ioutil.Discard)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version=%s\n", c.App.Version)
	}

	app := &cli.App{
		Name:    "pmctl",
		Version: "v0.1",
		Usage:   "Introspects and Controls the system",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "http://localhost:8080",
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
						fetchSystemdUnitStatus(c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Start (activate) one unit specified on the command line",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("start", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stop (deactivate) one specified on the command line.",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("stop", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "restart",
					Usage: "Stop and then start one unit specified on the command line. If the unit is not running yet, it will be started.",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("restart", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "mask",
					Usage: "Mask one unit, as specified on the command line",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("mask", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "unmask",
					Usage: "Unmask one unit file, as specified on the command line",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("unmask", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "try-restart",
					Usage: "Stop and then start one unit specified on the command line if the unit are running. This does nothing if unit is not running.",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("try-restart", c.Args().First(), c.String("url"))
						return nil
					},
				},
				{
					Name:  "reload-or-restart",
					Usage: "Reload one unit if they support it. If not, stop and then start instead. If the unit is not running yet, it will be started.",
					Action: func(c *cli.Context) error {
						PerformSystemdUnitCommand("reload-or-restart", c.Args().First(), c.String("url"))
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
