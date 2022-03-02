// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/vishvananda/netlink"
)

func configureNetDev(t *testing.T, n networkd.NetDev) error {
	var resp []byte
	var err error

	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/netdev/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure netdev: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure netdev: %v\n", j.Errors)
	}

	return nil
}

func TestNetDevCreateVLan(t *testing.T) {
	setupLink(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeLink(t, "test99")

	n := networkd.NetDev{
		Name:  "vlan99",
		Kind:  "vlan",
		Links: []string{"test99"},
		VLanSection: networkd.VLan{
			Id: 10,
		},
	}

	if err := configureNetDev(t, n); err != nil {
		t.Fatalf("Failed to create VLan: %v\n", err)
	}

	time.Sleep(time.Second * 5)

	if !validator.LinkExists("vlan99") {
		t.Fatalf("Failed to create vlan='vlan99'")
	}

	s, _ := system.ExecAndCapture("ip", "-d", "link", "show", "vlan99")
	fmt.Println(s)

	m, _, err := networkd.CreateOrParseNetDevFile("vlan99", "vlan")
	if err != nil {
		t.Fatalf("Failed to parse .netdev file of vlan='vlan99'")
	}

	if m.GetKeySectionString("NetDev", "Kind") != "vlan" {
		t.Fatalf("Vlan kind is not 'vlan' in .netdev file of vlan='vlan99'")
	}

	if m.GetKeySectionUint("VLAN", "Id") != 10 {
		t.Fatalf("Invalid Vlan Id in .netdev file of vlan='vlan99'")
	}

	m, err = networkd.CreateOrParseNetworkFile("vlan99")
	if err != nil {
		t.Fatalf("Failed to parse .network file of vlan='vlan99'")
	}

	if m.GetKeySectionString("Match", "Name") != "vlan99" {
		t.Fatalf("Invalid netdev name in .network file of vlan='vlan99'")
	}

	m, err = networkd.CreateOrParseNetworkFile("test99")
	if err != nil {
		t.Fatalf("Failed to parse .network file of test99")
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Network", "VLAN") != "vlan99" {
		t.Fatalf("Failed to parse .network file of test99")
	}

	if err := networkd.RemoveNetDev(n.Name, n.Kind); err != nil {
		t.Fatalf("Failed to remove .network file='%v'", err)
	}
}

func TestNetDevCreateBond(t *testing.T) {
	setupLink(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	setupLink(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test98"}})
	defer removeLink(t, "test99")
	defer removeLink(t, "test98")

	n := networkd.NetDev{
		Name:  "bond99",
		Kind:  "bond",
		Links: []string{"test99", "test98"},
		BondSection: networkd.Bond{
			Mode: "balance-rr",
		},
	}

	if err := configureNetDev(t, n); err != nil {
		t.Fatalf("Failed to create Bond: %v\n", err)
	}

	time.Sleep(time.Second * 5)

	s, _ := system.ExecAndCapture("ip", "-d", "link", "show", "bond99")
	fmt.Println(s)
	
	if !validator.LinkExists("bond99") {
		t.Fatalf("Failed to create bond='bond99'")
	}

	m, _, err := networkd.CreateOrParseNetDevFile("bond99", "bond")
	if err != nil {
		t.Fatalf("Failed to parse .netdev file of bond='bond99'")
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("NetDev", "Kind") != "bond" {
		t.Fatalf("Bond kind is not 'bond' in .netdev file of bond='bond99'")
	}

	if m.GetKeySectionString("Bond", "Mode") != "balance-rr" {
		t.Fatalf("Invalid bond mode in .netdev file of bond='bond99'")
	}

	m1, err := networkd.CreateOrParseNetworkFile("test99")
	if err != nil {
		t.Fatalf("Failed to parse .network file of test99")
	}
	defer os.Remove(m1.Path)

	if m1.GetKeySectionString("Network", "Bond") != "bond99" {
		t.Fatalf("Failed to parse Bond=bond99 in .network file")
	}
	
	m2, err := networkd.CreateOrParseNetworkFile("test98")
	if err != nil {
		t.Fatalf("Failed to parse .network file of test99")
	}
	defer os.Remove(m2.Path)

	if m2.GetKeySectionString("Network", "Bond") != "bond99" {
		t.Fatalf("Failed to parse Bond=bond99 in .network file")
	}

	if err := networkd.RemoveNetDev(n.Name, n.Kind); err != nil {
		t.Fatalf("Failed to remove .network file='%v'", err)
	}
}
