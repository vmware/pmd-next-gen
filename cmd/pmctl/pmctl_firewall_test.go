// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/fatih/color"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/firewall"
)

func addNFTTable() error {
	tbl := firewall.Nft{
		Table: firewall.Table{
			Name:   "test99",
			Family: "inet",
		},
	}

	resp, err := web.DispatchSocket(http.MethodPost, "", "/api/v1/network/firewall/nft/table/add", nil, tbl)
	if err != nil {
		return err
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		return err
	}

	if !m.Success {
		return fmt.Errorf("%v", m.Errors)
	}

	return nil
}

func deleteNFTTable() error {
	tbl := firewall.Nft{
		Table: firewall.Table{
			Name:   "test99",
			Family: "inet",
		},
	}

	resp, err := web.DispatchSocket(http.MethodDelete, "", "/api/v1/network/firewall/nft/table/remove", nil, tbl)
	if err != nil {
		return err
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		return err
	}

	if !m.Success {
		return fmt.Errorf("%v", m.Errors)
	}

	return nil
}

func TestAddNFTTable(t *testing.T) {
	if err := addNFTTable(); err != nil {
		t.Fatalf("Failed to add table: %v\n", err)
	}
}

func TestShowNFTTable(t *testing.T) {
	resp, err := web.DispatchSocket(http.MethodGet, "", "/api/v1/network/firewall/nft/table/show", nil, nil)
	if err != nil {
		t.Fatalf("Failed to acquire tables: %v\n", err)
	}

	ts := tableStats{}
	if err := json.Unmarshal(resp, &ts); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}

	if ts.Success {
		for _, v := range ts.Message {
			fmt.Printf("             %v %v\n", color.HiBlueString("Table:"), v.Name)
			fmt.Printf("            %v %v\n\n", color.HiBlueString("Family:"), v.Family)
		}
	} else {
		t.Fatalf(ts.Errors)
	}
}

func TestDeleteNFTTable(t *testing.T) {
	if err := deleteNFTTable(); err != nil {
		t.Fatalf("Failed to remove table: %v\n", err)
	}
}

func TestAddNFTChain(t *testing.T) {
	if err := addNFTTable(); err != nil {
		t.Fatalf("Failed to add table: %v\n", err)
	}

	c := firewall.Nft{
		Chain: firewall.Chain{
			Name:     "chaintest99",
			Table:    "test99",
			Family:   "inet",
			Hook:     "input",
			Priority: "300",
			Type:     "filter",
		},
	}

	resp, err := web.DispatchSocket(http.MethodPost, "", "/api/v1/network/firewall/nft/chain/add", nil, c)
	if err != nil {
		t.Fatalf("Failed to add chain: %v\n", err)
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}

	if !m.Success {
		t.Fatalf("Failed to add chain: %v\n", m.Errors)
	}
}

func TestShowNFTChain(t *testing.T) {
	resp, err := web.DispatchSocket(http.MethodGet, "", "/api/v1/network/firewall/nft/chain/show", nil, nil)
	if err != nil {
		t.Fatalf("Failed to acquire chains: %v\n", err)
	}

	c := chainStats{}
	if err := json.Unmarshal(resp, &c); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}

	if c.Success {
		for _, v := range c.Message {
			fmt.Printf("               %v %v\n", color.HiBlueString("Table:"), v.Table.Name)
			fmt.Printf("              %v %v\n", color.HiBlueString("Family:"), v.Table.Family)
			fmt.Printf("               %v %v\n", color.HiBlueString("Chain:"), v.Name)
			fmt.Printf("                %v %v\n", color.HiBlueString("Hook:"), v.Hooknum)
			fmt.Printf("                %v %v\n", color.HiBlueString("Type:"), v.Type)
			fmt.Printf("            %v %v\n\n", color.HiBlueString("Priority:"), v.Priority)
		}
	} else {
		t.Fatalf(c.Errors)
	}
}

func TestDeleteNFTChain(t *testing.T) {
	c := firewall.Nft{
		Chain: firewall.Chain{
			Name:   "chaintest99",
			Table:  "test99",
			Family: "inet",
		},
	}

	resp, err := web.DispatchSocket(http.MethodDelete, "", "/api/v1/network/firewall/nft/chain/remove", nil, c)
	if err != nil {
		t.Fatalf("Failed to remove chain: %v\n", err)
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}

	if !m.Success {
		t.Fatalf("Failed to remove chain: %v\n", m.Errors)
	}

	if err := deleteNFTTable(); err != nil {
		t.Fatalf("Failed to remove table: %v\n", err)
	}
}
