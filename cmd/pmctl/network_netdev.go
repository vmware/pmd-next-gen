// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/urfave/cli/v2"
)

func networkCreateVLan(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "vlan",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		case "id":
			if validator.IsVLanId(argStrings[i+1]) {
				n.VLanSection.Id = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse VLan Id: %s\n", argStrings[i+1])
				return
			}
		}

		i++
	}

	if validator.IsEmpty(n.Link) || validator.IsEmpty(n.VLanSection.Id) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create VLan. Missing VLAN name, dev or id\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create VLan: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create VLan: %v\n", m.Errors)
	}
}

func networkCreateBond(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "bond",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		case "mode":
			n.BondSection.Mode = argStrings[i+1]
		case "thp":
			n.BondSection.TransmitHashPolicy = argStrings[i+1]
		case "ltr":
			n.BondSection.LACPTransmitRate = argStrings[i+1]
		case "mms":
			n.BondSection.MIIMonitorSec = argStrings[i+1]
		}

		i++
	}

	if validator.IsEmpty(n.Link) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create Bond. Missing BOND name, dev or mode, ltr and mms\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create Bond: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create Bond: %v\n", m.Errors)
	}
}

func networkCreateBridge(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "bridge",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		}

		i++
	}

	if validator.IsEmpty(n.Link) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create Bridge. Missing BRIDGE name or dev\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create Bridge: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create Bridge: %v\n", m.Errors)
	}
}

func networkCreateMacVLan(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "macvlan",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Link = argStrings[i+1]
		case "mode":
			n.MacVLanSection.Mode = argStrings[i+1]
		}

		i++
	}

	if validator.IsEmpty(n.Link) || validator.IsEmpty(n.Name) || validator.IsEmpty(n.MacVLanSection.Mode) {
		fmt.Printf("Failed to create Macvlan. Missing MACVLAN name, dev or mode\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create MacVLan: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create MacVLan: %v\n", m.Errors)
	}
}
