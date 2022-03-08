// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/firewall"
	"github.com/urfave/cli/v2"
)

type tableStats struct {
	Success bool                      `json:"success"`
	Message map[string]firewall.Table `json:"message"`
	Errors  string                    `json:"errors"`
}

func parseNftTable(args cli.Args) (*firewall.Nft, error) {
	argStrings := args.Slice()
	n := firewall.Nft{}

	for i, args := range argStrings {
		switch args {
		case "family":
			if !validator.IsNFTFamily(argStrings[i+1]) {
				return nil, fmt.Errorf("Failed to parse family: '%s'", argStrings[i+1])
			}
			n.Table.Family = argStrings[i+1]
		case "tablename":
			if validator.IsEmpty(argStrings[i+1]) {
				return nil, fmt.Errorf("Failed to parse table-name: '%s'", argStrings[i+1])
			}
			n.Table.Name = argStrings[i+1]
		}
	}

	return &n, nil
}

func networkAddNftTable(args cli.Args, host string, token map[string]string) {
	n, err := parseNftTable(args)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/firewall/nft/tables/add", token, n)
	if err != nil {
		fmt.Printf("Failed to add table: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to add table %v\n", m.Errors)
	}
}

func networkShowNftTable(args cli.Args, host string, token map[string]string) {
	n, err := parseNftTable(args)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/network/firewall/nft/tables/show", token, n)
	if err != nil {
		fmt.Printf("Failed to show table: %v\n", err)
		return
	}

	ts := tableStats{}
	if err := json.Unmarshal(resp, &ts); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !ts.Success {
		fmt.Printf("Failed to show table %v\n", ts.Errors)
	}

	for _, v := range ts.Message {
		fmt.Printf("              %v %v\n", color.HiBlueString("Name:"), v.Name)
		fmt.Printf("            %v %v\n\n", color.HiBlueString("Family:"), v.Family)
	}

	return
}

func networkDeleteNftTable(args cli.Args, host string, token map[string]string) {
	n, err := parseNftTable(args)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	resp, err := web.DispatchSocket(http.MethodDelete, host, "/api/v1/network/firewall/nft/tables/remove", token, n)
	if err != nil {
		fmt.Printf("Failed to delete table: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to delete table %v\n", m.Errors)
	}
}

func networkSaveNftTable(host string, token map[string]string) {
	resp, err := web.DispatchSocket(http.MethodPut, host, "/api/v1/network/firewall/nft/tables/save", token, nil)
	if err != nil {
		fmt.Printf("Failed to save table: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to save table: %v\n", m.Errors)
	}

	return
}
