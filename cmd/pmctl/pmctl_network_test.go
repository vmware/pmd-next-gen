// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
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
		t.Fatalf("Failed to add global Dns server: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to add Dns servers: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	dns := m.GetKeySectionString("Resolve", "DNS")
	for _, d := range s {
		if !share.StringContains(strings.Split(dns, " "), d) {
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
		t.Fatalf("Failed to add global Dns servers: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Dns: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	dns := m.GetKeySectionString("Resolve", "DNS")
	for _, d := range s {
		if share.StringContains(strings.Split(dns, " "), d) {
			t.Fatalf("Failed")
		}
	}
}

func TestNetworkAddGlobalDomain(t *testing.T) {
	s := []string{"test1.com", "test2.com", "test3.com", "test4.com"}
	n := resolved.GlobalDns{
		Domains: s,
	}
	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/resolved/add", nil, n)
	if err != nil {
		t.Fatalf("Failed to add global domain: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure domain: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	domains := m.GetKeySectionString("Resolve", "Domains")
	for _, d := range s {
		if !share.StringContains(strings.Split(domains, " "), d) {
			t.Fatalf("Failed")
		}
	}
}

func TestNetworkRemoveGlobalDomain(t *testing.T) {
	TestNetworkAddGlobalDomain(t)
	s := []string{"test1.com", "test2.com"}
	n := resolved.GlobalDns{
		Domains: s,
	}
	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodDelete, "", "/api/v1/network/resolved/remove", nil, n)
	if err != nil {
		t.Fatalf("Failed to add global Ddomain: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to remove Domain: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)

	m, err := configfile.Load("/etc/systemd/resolved.conf")
	if err != nil {
		t.Fatalf("Failed to load resolved.conf: %v\n", err)
	}

	domains := m.GetKeySectionString("Resolve", "Domains")
	for _, d := range s {
		if share.StringContains(strings.Split(domains, " "), d) {
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

func TestNetworkDHCP4ClientIdentifier(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		DHCPv4Section: networkd.DHCPv4Section{
			ClientIdentifier: "duid",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure DHCP4ClientIdentifier: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure DHCP4ClientIdentifier: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure DHCP4ClientIdentifier: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure DHCP4ClientIdentifier: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("DHCPv4", "ClientIdentifier") != "duid" {
		t.Fatalf("Failed to set ClientIdentifier")
	}
}

func TestNetworkDHCPIAID(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		DHCPv4Section: networkd.DHCPv4Section{
			IAID: "8765434",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure DHCPIAID: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure DHCPIAID: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure DHCPIAID: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure DHCPIAID: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("DHCPv4", "IAID") != "8765434" {
		t.Fatalf("Failed to set IAID")
	}
}

func TestNetworkLinkMode(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			ARP:               "yes",
			Multicast:         "yes",
			AllMulticast:      "no",
			Promiscuous:       "no",
			RequiredForOnline: "yes",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link Mode: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link Mode: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link Mode: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link Mode: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "ARP") != "yes" {
		t.Fatalf("Failed to set ARP")
	}
	if m.GetKeySectionString("Link", "Multicast") != "yes" {
		t.Fatalf("Failed to set Multicast")
	}
	if m.GetKeySectionString("Link", "AllMulticast") != "no" {
		t.Fatalf("Failed to set AllMulticast")
	}
	if m.GetKeySectionString("Link", "Promiscuous") != "no" {
		t.Fatalf("Failed to set Promiscuous")
	}
	if m.GetKeySectionString("Link", "RequiredForOnline") != "yes" {
		t.Fatalf("Failed to set RequiredForOnline")
	}
}

func TestNetworkLinkMTU(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			MTUBytes: "2048",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link MTU: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link MTU: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link MTU: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link MTU: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "MTUBytes") != "2048" {
		t.Fatalf("Failed to set MTUBytes")
	}
}

func TestNetworkLinkMAC(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			MACAddress: "00:a0:de:63:7a:e6",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link MAC: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link MAC: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link MAC: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link MAC: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "MACAddress") != "00:a0:de:63:7a:e6" {
		t.Fatalf("Failed to set MACAddress")
	}
}

func TestNetworkLinkGroup(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			Group: "2147483647",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link Group: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link Group: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link Group: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link Group: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "Group") != "2147483647" {
		t.Fatalf("Failed to set Group")
	}
}

func TestNetworkLinkOnlineFamily(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			RequiredFamilyForOnline: "ipv4",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link OnlineFamily: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link OnlineFamily: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link OnlineFamily: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link OnlineFamily: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "RequiredFamilyForOnline") != "ipv4" {
		t.Fatalf("Failed to set RequiredFamilyForOnline")
	}
}

func TestNetworkLinkActPolicy(t *testing.T) {
	setupDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})
	defer removeDummy(t, &netlink.Dummy{netlink.LinkAttrs{Name: "test99"}})

	system.ExecRun("systemctl", "restart", "systemd-networkd")
	time.Sleep(time.Second * 3)

	n := networkd.Network{
		Link: "test99",
		LinkSection: networkd.LinkSection{
			ActivationPolicy: "always-up",
		},
	}

	var resp []byte
	var err error
	resp, err = web.DispatchSocket(http.MethodPost, "", "/api/v1/network/networkd/network/configure", nil, n)
	if err != nil {
		t.Fatalf("Failed to configure Link ActPolicy: %v\n", err)
	}

	j := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &j); err != nil {
		t.Fatalf("Failed to decode json message: %v\n", err)
	}
	if !j.Success {
		t.Fatalf("Failed to configure Link ActPolicy: %v\n", j.Errors)
	}

	time.Sleep(time.Second * 3)
	link, err := netlink.LinkByName("test99")
	network, err := networkd.ParseLinkNetworkFile(link.Attrs().Index)
	if err != nil {
		t.Fatalf("Failed to configure Link ActPolicy: %v\n", err)
	}

	m, err := configfile.Load(network)
	if err != nil {
		t.Fatalf("Failed to configure Link ActPolicy: %v\n", err)
	}
	defer os.Remove(m.Path)

	if m.GetKeySectionString("Link", "ActivationPolicy") != "always-up" {
		t.Fatalf("Failed to set ActivationPolicy")
	}
}
