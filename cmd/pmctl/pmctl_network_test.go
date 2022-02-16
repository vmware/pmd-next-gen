// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pmd-nextgen/pkg/configfile"
	"github.com/pmd-nextgen/pkg/share"
	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/pmd-nextgen/plugins/network/resolved"
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

func TestNetworkAddGlobalDns(t *testing.T) {
	s := []string{"8.8.8.8", "8.8.4.4", "8.8.8.1", "8.8.8.2"}
	n := resolved.GlobalDns{
		DnsServers: s,
	}
	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/resolved/add", nil, n)
	if err != nil {
		fmt.Printf("Failed to add global Dns server: %v\n", err)
		return
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure DNS: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	dns := m.GetKeySectionString("Resolve", "DNS")
	for _, d := range strings.Split(dns, " ") {
		if !share.StringContains(s, d) {
			t.Fatalf("Failed")
		}
	}
}

func TestNetworkRemoveGlobalDns(t *testing.T) {
	TestNetworkAddGlobalDns(t)
	s := []string{"8.8.8.8", "8.8.4.4"}
	n := resolved.GlobalDns{
		DnsServers: s,
	}
	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodDelete, "", "/api/v1/network/resolved/remove", nil, n)
	if err != nil {
		fmt.Printf("Failed to add global Dns server: %v\n", err)
		return
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure DNS: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	dns := m.GetKeySectionString("Resolve", "DNS")
	for _, d := range strings.Split(dns, " ") {
		if share.StringContains(s, d) {
			t.Fatalf("Failed")
		}
	}
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
