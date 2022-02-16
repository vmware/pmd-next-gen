// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pmd-nextgen/pkg/configfile"
	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/vishvananda/netlink"
)

func setupDummy(t *testing.T, link netlink.Link) {
	_, err := netlink.LinkList()
	if err != nil {
		t.Fatal(err)
	}

	if err := netlink.LinkAdd(link); err != nil && err.Error() != "file exists" {
		t.Fatal(err)
	}

	_, err = netlink.LinkByName(link.Attrs().Name)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 1)
}

func removeDummy(t *testing.T, link netlink.Link) {
	l, err := netlink.LinkByName(link.Attrs().Name)
	if err != nil {
		t.Fatal(err)
	}

	netlink.LinkDel(l)
}

func TestNetworkDHCP(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		NetworkSection: networkd.NetworkSection{
			DHCP: "ipv4",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure DHCP: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure DHCP: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure DHCP: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure DHCP: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Network", "DHCP") != "ipv4" {
		t.Fatalf("Failed to set DHCP")
	}
}
