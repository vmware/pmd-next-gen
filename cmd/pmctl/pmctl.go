// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/pm-web/pkg/web"
)

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
			Name:  "service",
			Usage: "Introspects and controls the systemd services",
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
			Aliases: []string{"s"},
			Usage:   "Introspects of system or network status",
			Subcommands: []*cli.Command{
				{
					Name:    "network",
					Aliases: []string{"n"},
					Usage:   "Introspects network status",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "interface", Aliases: []string{"i"}},
					},

					Action: func(c *cli.Context) error {
						acquireNetworkStatus("network", c.String("url"), c.String("interface"), token)
						return nil
					},
					Subcommands: []*cli.Command{
						{
							Name:  "iostat",
							Usage: "Show iostat of interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("iostat", c.String("url"), "", token)
								return nil
							},
						},
						{
							Name:  "interfaces",
							Usage: "Show network interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("interfaces", c.String("url"), "", token)
								return nil
							},
						},
					},
				},
				{
					Name:    "system",
					Aliases: []string{"s"},
					Usage:   "Introspects system status",

					Action: func(c *cli.Context) error {
						acquireSystemStatus(c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "create a new user or update user information",
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"n"},
					Usage:   "Add a new user",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "home-dir", Aliases: []string{"d"}},
						&cli.StringFlag{Name: "groups", Usage: "Separate by ,"},
						&cli.StringFlag{Name: "uid"},
						&cli.StringFlag{Name: "gid"},
						&cli.StringFlag{Name: "shell"},
						&cli.StringFlag{Name: "password"},
					},

					Action: func(c *cli.Context) error {
						userAdd(c.Args().First(), c.String("uid"), c.String("groups"), c.String("gid"), c.String("shell"), c.String("home-dir"), c.String("password"), c.String("gid"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:    "remove",
					Aliases: []string{"n"},
					Usage:   "Remove an existing user",

					Action: func(c *cli.Context) error {
						userRemove(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "create a new group or update group information",
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "Add a new group",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "gid"},
					},

					Action: func(c *cli.Context) error {
						groupAdd(c.Args().First(), c.String("gid"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:    "remove",
					Aliases: []string{"a"},
					Usage:   "Remove an existing group",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "gid"},
					},

					Action: func(c *cli.Context) error {
						groupRemove(c.Args().First(), c.String("gid"), c.String("url"), token)
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
