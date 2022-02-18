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
		fmt.Printf("Failed to create VLan: %v\n", err)
		return
	}

	time.Sleep(time.Second * 5)
	s, _:=system.ExecAndCapture("ip", "-d", "link", "show", "vlan99")
	fmt.Println(s)
	if !validator.LinkExists("vlan99") {
		t.Fatalf("Failed to create vlan='vlan99'")
	}
	defer removeLink(t, "vlan99")

	m, _, err := networkd.CreateOrParseNetDevFile("vlan99", "vlan")
	if err != nil {
		t.Fatalf("Failed to parse .netdev file of vlan='vlan99'")
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionUint("VLAN", "Id") != 10 {
		t.Fatalf("Invalid Vlan Id in .netdev file of vlan='vlan99'")
	}
}
