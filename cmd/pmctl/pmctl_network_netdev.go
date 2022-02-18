// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			n.Links = strings.Fields(argStrings[i+1])
		case "id":
			if validator.IsVLanId(argStrings[i+1]) {
				id, err := strconv.ParseUint(argStrings[i+1], 10, 32)
				if err != nil {
					fmt.Printf("Failed to parse VLan Id: %s\n", argStrings[i+1])
					return
				}
				n.VLanSection.Id = uint(id)
			}
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || n.VLanSection.Id == 0 || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create VLan. Missing VLan name, dev or id\n")
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
			n.Links = strings.Split(argStrings[i+1], ",")
		case "mode":
			if validator.IsBondMode(argStrings[i+1]) {
				n.BondSection.Mode = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse bond mode: %s\n", argStrings[i+1])
				return
			}
		case "thp":
			if validator.IsBondTransmitHashPolicy(n.BondSection.Mode, argStrings[i+1]) {
				n.BondSection.TransmitHashPolicy = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse transmit hash policy: %s\n", argStrings[i+1])
				return
			}
		case "ltr":
			if validator.IsBondLACPTransmitRate(argStrings[i+1]) {
				n.BondSection.LACPTransmitRate = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse LACP transmit rate: %s\n", argStrings[i+1])
				return
			}
		case "mms":
			n.BondSection.MIIMonitorSec = argStrings[i+1]
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Name) {
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
			n.Links = strings.Fields(argStrings[i+1])
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create bridge. Missing bridge name or dev\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create bridge: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create bridge: %v\n", m.Errors)
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
			n.Links = strings.Fields(argStrings[i+1])
		case "mode":
			if validator.IsMacVLanMode(argStrings[i+1]) {
				n.MacVLanSection.Mode = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse mode: %s\n", argStrings[i+1])
				return
			}
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Name) || validator.IsEmpty(n.MacVLanSection.Mode) {
		fmt.Printf("Failed to create MacVLan. Missing MACVLAN name, dev or mode\n")
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

func networkCreateIpVLan(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "ipvlan",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Links = strings.Fields(argStrings[i+1])
		case "mode":
			if validator.IsIpVLanMode(argStrings[i+1]) {
				n.IpVLanSection.Mode = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse mode: %s\n", argStrings[i+1])
				return
			}
		case "flags":
			if validator.IsIpVLanFlags(argStrings[i+1]) {
				n.IpVLanSection.Flags = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse flags: %s\n", argStrings[i+1])
				return
			}
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to create IpVLan. Missing IPVLAN name or dev\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create IpVLan: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create IpVLan: %v\n", m.Errors)
	}
}

func networkCreateWireGuard(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
		Kind: "wireguard",
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Links = strings.Fields(argStrings[i+1])
		case "skey":
			n.WireGuardSection.PrivateKey = argStrings[i+1]
		case "pkey":
			n.WireGuardPeerSection.PublicKey = argStrings[i+1]
		case "port":
			if validator.IsWireGuardListenPort(argStrings[i+1]) {
				n.WireGuardSection.ListenPort = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse listen port: %s\n", argStrings[i+1])
				return
			}
		case "ips":
			ips := strings.Split(argStrings[i+1], ",")
			for _, ip := range ips {
				if !validator.IsIP(ip) {
					fmt.Printf("Failed to parse allowed ips: %s\n", argStrings[i+1])
					return
				}
			}
			n.WireGuardPeerSection.AllowedIPs = ips
		case "endpoint":
			if validator.IsWireGuardPeerEndpoint(argStrings[i+1]) {
				n.WireGuardPeerSection.Endpoint = argStrings[i+1]
			} else {
				fmt.Printf("Failed to parse endpoint: %s\n", argStrings[i+1])
				return
			}
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Name) || validator.IsEmpty(n.WireGuardSection.PrivateKey) ||
		validator.IsEmpty(n.WireGuardPeerSection.PublicKey) || validator.IsEmpty(n.WireGuardPeerSection.Endpoint) {
		fmt.Printf("Failed to create WireGuard. Missing WireGuard name, skey, pkey or dev\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/netdev/configure", token, n)
	if err != nil {
		fmt.Printf("Failed to create WireGuard: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to create WireGuard: %v\n", m.Errors)
	}
}

func networkRemoveNetDev(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	n := networkd.NetDev{
		Name: argStrings[0],
	}

	for i := 1; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			n.Links = strings.Fields(argStrings[i+1])
		case "kind":
			n.Kind = argStrings[i+1]
		}

		i++
	}

	if validator.IsArrayEmpty(n.Links) || validator.IsEmpty(n.Kind) || validator.IsEmpty(n.Name) {
		fmt.Printf("Failed to remove netdev. Missing name, dev or kind\n")
		return
	}

	resp, err := web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/networkd/netdev/remove", token, n)
	if err != nil {
		fmt.Printf("Failed to remove netdev: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to remove netdev %v\n", m.Errors)
	}
}
