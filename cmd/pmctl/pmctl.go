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

	"github.com/pmd-nextgen/pkg/validator"
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
							Name:        "dns",
							Description: "Show dns and domaains",

							Action: func(c *cli.Context) error {
								acquireResolveDescribe(c.String("url"), token)
								return nil
							},
						},
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
			Usage:   "Network device configuration",
			Subcommands: []*cli.Command{
				{
					Name:        "add-dns",
					UsageText:   "dev [LINK] dns [DNS]",
					Description: "Add Link or global DNS server address. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkAddDns(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove-dns",
					UsageText:   "dev [LINK] dns [DNS]",
					Description: "Remove Link or global DNS server address. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkRemoveDns(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "add-domain",
					UsageText:   "dev [LINK] domain [DNS]",
					Description: "Add Link or global DNS domain name. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkAddDomains(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove-domain",
					UsageText:   "dev [LINK] domain [DNS]",
					Description: "Remove Link or global DNS domain name. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 1 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkRemoveDomains(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "add-ntp",
					UsageText:   "dev [LINK] ntp [NTP]",
					Description: "Add Link or global NTP server address. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkAddNTP(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove-ntp",
					UsageText:   "dev [LINK] ntp [NTP]",
					Description: "Removes Link or global NTP server address. This option may be specified more than once separated by ,",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkRemoveNTP(c.Args(), c.String("url"), token)
						return nil
					},
				},
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
					Name:        "set-dhcp4-client-identifier",
					UsageText:   "set-dhcp4-client-identifier [LINK] [IDENTIFIER {mac|duid|duid-only}]",
					Description: "Configures Link DHCPv4 identifier.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureDHCP4ClientIdentifier(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-dhcp-iaid",
					UsageText:   "set-dhcp-iaid [LINK] [IAID]",
					Description: "Configures the DHCP Identity Association Identifier (IAID) for the interface, a 32-bit unsigned integer.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						if !validator.IsIaId(c.Args().Get(1)) {
							fmt.Printf("Invalid IAID.\n")
							return nil
						}

						networkConfigureDHCPIAID(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},

				{
					Name:        "set-mtu",
					UsageText:   "set-mtu [LINK] [MTU NUMBER]",
					Description: "Configures Link MTU.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						if !validator.IsMtu(c.Args().Get(1)) {
							fmt.Printf("MTU must be a valid value.\n")
							return nil
						}

						networkConfigureMTU(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-mac",
					UsageText:   "set-mac [LINK] [MAC]",
					Description: "Configures Link MAC address.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						if validator.IsNotMAC(c.Args().Get(1)) {
							fmt.Printf("Invalid MAC address: %v\n", c.Args().Get(1))
							return nil
						}

						networkConfigureMAC(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-group",
					UsageText:   "set-group [LINK] [GROUP INTEGER]",
					Description: "Configures Link Group.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkGroup(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-rf-online",
					UsageText:   "set-rf-online [LINK] [FAMILY STRING]",
					Description: "Configures Link RequiredFamilyForOnline.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkRequiredFamilyForOnline(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-active-policy",
					UsageText:   "set-active-policy [LINK] [POLICY STRING]",
					Description: "Configures Link ActivationPolicy.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkActivationPolicy(c.Args().First(), c.Args().Get(1), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-link-mode",
					UsageText:   "set-link-mode [LINK] mode [BOOLEAN] arp [BOOLEAN] mc [BOOLEAN] amc [BOOLEAN] pcs [BOOLEAN] rfo [BOOLEAN]",
					Description: "Set Link mode,arp,multicast,allmulticast,promiscuous and requiredforonline.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureMode(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "add-link-address",
					UsageText:   "add-link-address [LINK] address [ADDRESS] peer [ADDRESS] label [NUMBER] scope {global|link|host|NUMBER}]",
					Description: "Configures Link Address.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 2 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureAddress(c.Args().First(), c.Args(), c.String("url"), token)
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
				{
					Name:        "create-bond",
					UsageText:   "create-bond [BOND name] dev [LINK MASTER] mode [STRING] thp [STRING] ltr [string] mms [STRING]",
					Description: "Create bond.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 3 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateBond(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "create-bridge",
					UsageText:   "create-bridge [BRIDGE name] dev [LINK MASTER]",
					Description: "Create bridge.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 3 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateBridge(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "create-macvlan",
					UsageText:   "create-macvlan [MACVLAN name] dev [LINK MASTER] mode [STRING]",
					Description: "Create macvlan.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 5 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateMacVLan(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "create-ipvlan",
					UsageText:   "create-ipvlan [IPVLAN name] dev [LINK MASTER] mode [STRING] flags [STRING]",
					Description: "Create ipvlan.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 3 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateIpVLan(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "create-wg",
					UsageText:   "create-wg [WIREGUARD name] dev [LINK MASTER] skey [STRING] pkey [STRING] port [string] ips [string] endpoint [STRING]",
					Description: "Create wg(wireguard).",

					Action: func(c *cli.Context) error {
						if c.NArg() < 9 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkCreateWireGuard(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "remove-netdev",
					UsageText:   "remove-netdev [NETDEV name] dev [LINK MASTER] kind [KIND {vlan|bridge|bond|vxlan|macvlan|macvtap|ipvlan|ipvtap|vrf|veth|ipip|sit|vti|gre|wg]",
					Description: "Removes .netdev and .network files.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 5 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkRemoveNetDev(c.Args(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Network device configuration",
			Subcommands: []*cli.Command{
				{
					Name:        "set-mac",
					UsageText:   "set-mac dev [LINK] macpolicy [string] macaddr [string]",
					Description: "Sets the device's MAC configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkMAC(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-name",
					UsageText:   "set-name dev [LINK] namepolicy [string] name [string]",
					Description: "Sets the device's Name configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkName(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-alt-name",
					UsageText:   "set-alt-name dev [LINK] altnamespolicy [string] altname [string]",
					Description: "Sets the device's AlternativeName configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkAltName(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-csum-offload",
					UsageText:   "set-csum-offload dev [LINK] rxco [string] txco [string]",
					Description: "Sets the device's ChecksumOffload configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkChecksumOffload(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-tcp-offload",
					UsageText:   "set-tcp-offload dev [LINK] tcpso [string] tcp6so [string]",
					Description: "Sets the device's TCPSegmentationOffload configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkTCPOffload(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-generic-offload",
					UsageText:   "set-generic-offload dev [LINK] gso [string] gro [string] grohw [string] gsomaxbytes [int] gsomaxseg [int]",
					Description: "Sets the device's GenericOffload configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkGenericOffload(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-vlan-tags",
					UsageText:   "set-vlan-tags dev [LINK] rxvlanctaghwacl [string] txvlanctaghwacl [string] rxvlanctagfilter [string] txvlanstaghwacl [string]",
					Description: "Sets the device's VLANTags configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkVLANTags(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-channel",
					UsageText:   "set-channel dev [LINK] rxch [int] txch [int] och [int] coch [int]",
					Description: "Sets the device's Channel configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkChannel(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-buffer",
					UsageText:   "set-buffer dev [LINK] rxbufsz [int] rxmbufsz [int] rxjbufsz [int] txbufsz [int]",
					Description: "Sets the device's Buffer configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkBuffer(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-flow-ctrl",
					UsageText:   "set-flow-ctrl dev [LINK] rxfctrl [string] txfctrl [string] anfctrl [string]",
					Description: "Sets the device's Buffer configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkFlowControl(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-adpt-coalesce",
					UsageText:   "set-adpt-coalesce dev [LINK] uarxc [string] uatxc [string]",
					Description: "Sets the device's AdpativeCoalesce configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkAdaptiveCoalesce(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-rx-coalesce",
					UsageText:   "set-rx-coalesce dev [LINK] rxcs [int] rxcsirq [int] rxcslow [int] rxcshigh [int]",
					Description: "Sets the device's ReceiveCoalesce configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkRxCoalesce(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-tx-coalesce",
					UsageText:   "set-tx-coalesce dev [LINK] txcs [int] txcsirq [int] txcslow [int] txcshigh [int]",
					Description: "Sets the device's TransmitCoalesce configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkTxCoalesce(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-rx-coald-frames",
					UsageText:   "set-rx-coald-frames dev [LINK] rxmcf [int] rxmcfirq [int] rxmcflow [int] rxmcfhigh [int]",
					Description: "Sets the device's ReceiveCoalescedFrames configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkRxCoalescedFrames(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-tx-coald-frames",
					UsageText:   "set-tx-coald-frames dev [LINK] txmcf [int] txmcfirq [int] txmcflow [int] txmcfhigh [int]",
					Description: "Sets the device's TransmitCoalescedFrames configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkTxCoalescedFrames(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-coalesce-pkt",
					UsageText:   "set-coalesce-pkt dev [LINK] cprlow [int] cprhigh [int] cprsis [int]",
					Description: "Sets the device's CoalescePacketRate configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkCoalescePacketRate(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-link",
					UsageText:   "set-link dev [LINK] alias,desc... [string]",
					Description: "Sets the device's alias,description,port,duplex,wakeonlan.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLink(c.Args(), c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "set-queue",
					UsageText:   "set-queue dev [LINK] txq [TransmitQueues {1...4096}] rxq [TransmitQueues {1...4096}] txqlen [TransmitQueuesLength {0...4294967294}]",
					Description: "Sets the device's queue configuration.",

					Action: func(c *cli.Context) error {
						if c.NArg() < 4 {
							fmt.Printf("Too few arguments.\n")
							return nil
						}

						networkConfigureLinkQueue(c.Args(), c.String("url"), token)
						return nil
					},
				},
			},
		},
		{
			Name:    "pkg",
			Aliases: []string{"p", "tdnf"},
			Usage:   "Package Management",
			Flags:   tdnfCreateFlags(),
			Subcommands: []*cli.Command{
				{
					Name:        "clean",
					Aliases:     []string{"c"},
					Description: "Clean Package Metadata",

					Action: func(c *cli.Context) error {
						options := tdnfParseFlags(c)
						tdnfClean(&options, c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "list",
					Aliases:     []string{"l"},
					Description: "List Packages",

					Action: func(c *cli.Context) error {
						options := tdnfParseFlags(c)
						if c.NArg() >= 1 {
							tdnfList(&options, c.Args().First(), c.String("url"), token)
						} else {
							tdnfList(&options, "", c.String("url"), token)
						}
						return nil
					},
				},
				{
					Name:        "makecache",
					Aliases:     []string{"mc"},
					Description: "Download Package Metadata",

					Action: func(c *cli.Context) error {
						options := tdnfParseFlags(c)
						tdnfMakeCache(&options, c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "repolist",
					Aliases:     []string{"rl"},
					Description: "List Repositories",

					Action: func(c *cli.Context) error {
						options := tdnfParseFlags(c)
						tdnfRepoList(&options, c.String("url"), token)
						return nil
					},
				},
				{
					Name:        "info",
					Aliases:     []string{"i"},
					Description: "Package Info",

					Action: func(c *cli.Context) error {
						options := tdnfParseFlags(c)
						if c.NArg() >= 1 {
							tdnfInfoList(&options, c.Args().First(), c.String("url"), token)
						} else {
							tdnfInfoList(&options, "", c.String("url"), token)
						}
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
