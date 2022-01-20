// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/pmd-nextgen/pkg/web"
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
					Name:        "status",
					Description: "Show terse runtime status information about one unit",

					Action: func(c *cli.Context) error {
						acquireSystemdUnitStatus(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "start",
					Description: "Start (activate) one unit specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("start", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "stop",
					Description: "Stop (deactivate) one specified on the command line.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("stop", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "restart",
					Description: "Stop and then start one unit specified on the command line. If the unit is not running yet, it will be started.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("restart", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "mask",
					Description: "Mask one unit, as specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("mask", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "unmask",
					Description: "Unmask one unit file, as specified on the command line",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("unmask", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "try-restart",
					Description: "Stop and then start one unit specified on the command line if the unit are running. This does nothing if unit is not running.",
					Action: func(c *cli.Context) error {
						executeSystemdUnitCommand("try-restart", c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "reload-or-restart",
					Description: "Reload one unit if they support it. If not, stop and then start instead. If the unit is not running yet, it will be started.",
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
					Name:        "network",
					Aliases:     []string{"n"},
					Description: "Introspects network status",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "interface", Aliases: []string{"i"}},
					},

					Action: func(c *cli.Context) error {
						acquireNetworkStatus("network", c.String("url"), c.String("interface"), token)
						return nil
					},
					Subcommands: []*cli.Command{
						{
							Name:        "iostat",
							Description: "Show iostat of interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("iostat", c.String("url"), "", token)
								return nil
							},
						},
						{
							Name:        "interfaces",
							Description: "Show network interfaces",

							Action: func(c *cli.Context) error {
								acquireNetworkStatus("interfaces", c.String("url"), "", token)
								return nil
							},
						},
					},
				},
				{
					Name:        "system",
					Aliases:     []string{"s"},
					Description: "Introspects system status",

					Action: func(c *cli.Context) error {
						acquireSystemStatus(c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "group",
					Aliases:     []string{"g"},
					Description: "Introspects group status",

					Action: func(c *cli.Context) error {
						acquireGroupStatus(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "user",
					Aliases:     []string{"u"},
					Description: "Introspects user status",

					Action: func(c *cli.Context) error {
						acquireUserStatus(c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "sysctl",
					Aliases:     []string{"s"},
					Description: "Introspects sysctl status",

					Action: func(c *cli.Context) error {
						acquireSysctlStatus("statusall", "", c.String("url"), token)
						return nil
					},
					Subcommands: []*cli.Command{
						{
							Name:        "key",
							Aliases:     []string{"k"},
							Description: "Show sysctl configuration based on key",

							Action: func(c *cli.Context) error {
								if c.NArg() < 1 {
									fmt.Printf("sysctl: No key is specified\n")
									return nil
								}

								acquireSysctlParamStatus(c.Args().First(), c.String("url"), token)
								return nil
							},
						},
						{
							Name:        "pattern",
							Aliases:     []string{"p"},
							Description: "Show sysctl configuration based on pattern",
							Flags: []cli.Flag{
								&cli.StringFlag{Name: "pattern"},
							},

							Action: func(c *cli.Context) error {
								if c.NArg() < 1 {
									fmt.Printf("sysctl: No pattern is specified\n")
									return nil
								}
								acquireSysctlStatus("statuspattern", c.Args().First(), c.String("url"), token)
								return nil
							},
						},
					},
				},
			},
		},
		{
			Name:    "system",
			Aliases: []string{"s"},
			Usage:   "Configures system",
			Subcommands: []*cli.Command{
				{
					Name:        "set-hostname",
					Description: "Set system hostname",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("No hostname suppplied\n")
							return nil
						}

						SetHostname(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "network",
			Aliases: []string{"n"},
			Usage:   "Configures network",
			Subcommands: []*cli.Command{
				{
					Name:        "set-dhcp",
					UsageText:   "set-dhcp [LINK] [DHCP-MODE {yes|no|ipv4|ipv6}]",
					Description: "Enables DHCPv4 and/or DHCPv6 client support. Accepts \"yes\", \"no\", \"ipv4\", or \"ipv6\".",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureDHCP(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "create-vlan",
					UsageText:   "create-vlan [VLAN name] dev [LINK MASTER] id [ID INTEGER]",
					Description: "Create vlan.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 5 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateVLan(c.Args(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "pkg",
			Aliases: []string{"p", "tdnf"},
			Usage:   "Package Management",
			Subcommands: []*cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Description: "List Packages",

					Action:  func(c *cli.Context) error {
						tdnfList(c.String("url"), token)
						return nil
					},
				},
				{
					Name:    "repolist",
					Aliases: []string{"rl"},
					Description: "List Repositories",

					Action:  func(c *cli.Context) error {
						tdnfRepoList(c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "Create a new user or update user information",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Aliases:     []string{"a"},
					Description: "Add a new user",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "home-dir", Aliases: []string{"d"}},
						&cli.StringFlag{Name: "groups", Aliases: []string{"grp"}, Usage: "Separate by ,"},
						&cli.StringFlag{Name: "uid", Aliases: []string{"u"}},
						&cli.StringFlag{Name: "gid", Aliases: []string{"g"}},
						&cli.StringFlag{Name: "shell", Aliases: []string{"s"}},
						&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
					},

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("No user name suppplied\n")
							return nil
						}
						userAdd(c.Args().First(), c.String("uid"), c.String("groups"), c.String("gid"), c.String("shell"), c.String("home-dir"), c.String("password"), c.String("gid"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove",
					Aliases:     []string{"r"},
					Description: "Remove an existing user",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("No user name suppplied\n")
							return nil
						}
						userRemove(c.Args().First(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "Create a new group or update group information",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Aliases:     []string{"a"},
					Description: "Add a new group",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "gid"},
					},

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("No group name suppplied\n")
							return nil
						}
						groupAdd(c.Args().First(), c.String("gid"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove",
					Aliases:     []string{"r"},
					Description: "Remove an existing group",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "gid"},
					},

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("No group name suppplied\n")
							return nil
						}
						groupRemove(c.Args().First(), c.String("gid"), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "sysctl",
			Aliases: []string{"s"},
			Usage:   "Add or Update, remove and load sysctl configuration",
			Subcommands: []*cli.Command{
				{
					Name:        "update",
					Aliases:     []string{"u"},
					Description: "Add or update sysctl cofiguration",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "key", Aliases: []string{"k"}},
						&cli.StringFlag{Name: "value", Aliases: []string{"v"}},
						&cli.StringFlag{Name: "filename", Aliases: []string{"f"}},
					},

					Action: func(c *cli.Context) error {
						sysctlUpdateConfig(c.String("key"), c.String("value"), c.String("filename"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove",
					Aliases:     []string{"r"},
					Description: "Remove an entry from sysctl configuration",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "key", Aliases: []string{"k"}},
						&cli.StringFlag{Name: "filename", Aliases: []string{"f"}},
					},

					Action: func(c *cli.Context) error {
						sysctlRemoveConfig(c.String("key"), c.String("filename"), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "load",
					Aliases:     []string{"l"},
					Description: "Load sysctl configuration from files",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "files", Aliases: []string{"f"}, Usage: "Separate by ,"},
					},

					Action: func(c *cli.Context) error {
						sysctlLoadConfig(c.String("files"), c.String("url"), token)
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
