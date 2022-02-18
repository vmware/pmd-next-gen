// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package validator

import (
	"net"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/vishvananda/netlink"
)

func IsBool(str string) bool {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On", "0", "f", "F",
		"false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return true
	}

	return false
}

func BoolToString(str string) string {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return "yes"
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return "no"
	}

	return "n/a"
}

func IsArrayEmpty(str []string) bool {
	return len(str) == 0
}

func IsEmpty(str string) bool {
	return govalidator.IsNull(str)
}

func IsUintOrMax(s string) bool {
	if strings.EqualFold(s, "max") {
		return true
	}

	_, err := strconv.ParseUint(s, 10, 32)
	return err == nil
}

func IsPort(port string) bool {
	_, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return false
	}

	return true
}

func IsHost(host string) bool {
	_, err := net.LookupHost(host)
	if err != nil {
		return false
	}

	return true
}

func IsValidIP(ip string) bool {
	a := net.ParseIP(ip)
	if a.To4() == nil && a.To16() == nil {
		return false
	}

	return true
}

func IsIP(str string) bool {
	_, _, err := net.ParseCIDR(str)
	if err != nil {
		ip := net.ParseIP(str)
		return ip != nil
	}

	return err == nil
}

func IsIPs(s []string) bool {
	for _, ip := range s {
		if !IsValidIP(ip) {
			return false
		}
	}

	return true
}

func IsClientIdentifier(identifier string) bool {
	return identifier == "mac" || identifier == "duid" || identifier == "duid-only"
}

func IsNotMAC(mac string) bool {
	return !govalidator.IsMAC(mac)
}

func IsMtu(mtu string) bool {
	_, err := strconv.ParseUint(mtu, 10, 32)
	return err == nil
}

func IsIaId(iaid string) bool {
	_, err := strconv.ParseUint(iaid, 10, 32)
	return err == nil
}

func IsVLanId(id string) bool {
	_, err := strconv.ParseUint(id, 10, 32)
	return err == nil
}

func IsScope(s string) bool {
	switch s {
	case "global", "link", "host":
		return true
	}

	scope, err := strconv.ParseUint(s, 10, 32)
	if err != nil || scope >= 256 {
		return false
	}

	return true
}

func IsBoolWithIp(s string) bool {
	switch s {
	case "yes", "no", "ipv4", "ipv6":
		return true
	}

	return false
}

func IsDHCP(s string) bool {
	return IsBoolWithIp(s)
}

func IsLinkLocalAddressing(s string) bool {
	return IsBoolWithIp(s)
}

func IsMulticastDNS(s string) bool {
	return IsBool(s) || s == "resolve"
}

func IsBondMode(mode string) bool {
	return mode == "balance-rr" || mode == "active-backup" || mode == "balance-xor" ||
		mode == "broadcast" || mode == "802.3ad" || mode == "balance-tlb" || mode == "balance-alb"
}

func IsBondTransmitHashPolicy(mode, thp string) bool {
	if (thp == "layer2" || thp == "layer3+4" || thp == "layer2+3" || thp == "encap2+3" || thp == "encap3+4") &&
		(mode == "balance-xor" || mode == "802.3ad" || mode == "balance-tlb") {
		return true
	}

	return false
}

func IsBondLACPTransmitRate(ltr string) bool {
	return ltr == "slow" || ltr == "fast"
}

func IsMacVLanMode(mode string) bool {
	return mode == "private" || mode == "vepa" || mode == "bridge" || mode == "passthru" || mode == "source"
}

func IsIpVLanMode(mode string) bool {
	return mode == "l2" || mode == "l3" || mode == "l3s"
}

func IsIpVLanFlags(flags string) bool {
	return flags == "bridge" || flags == "private" || flags == "vepa"
}

func IsWireGuardListenPort(port string) bool {
	return port == "auto" || IsPort(port)
}

func IsWireGuardPeerEndpoint(endPoint string) bool {
	ip, port, err := net.SplitHostPort(endPoint)
	if err != nil {
		return false
	}
	if !IsValidIP(ip) && !IsHost(ip) {
		return false
	}
	if !IsPort(port) {
		return false
	}

	return true
}

func IsLinkMACAddressPolicy(policy string) bool {
	return policy == "persistent" || policy == "random" || policy == "none"
}

func IsLinkNamePolicy(policy string) bool {
	return policy == "kernel" || policy == "database" || policy == "onboard" ||
		policy == "slot" || policy == "path" || policy == "mac" || policy == "keep"
}

func IsLinkName(name string) bool {
	if strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "ens") || strings.HasPrefix(name, "lo") {
		return false
	}

	return true
}

func IsLinkAlternativeNamesPolicy(policy string) bool {
	return policy == "database" || policy == "onboard" || policy == "slot" ||
		policy == "path" || policy == "mac"
}

func IsLinkQueue(id string) bool {
	l, err := strconv.ParseUint(id, 10, 32)
	if err != nil || l > 4096 {
		return false
	}

	return true
}

func IsLinkQueueLength(queueLength string) bool {
	l, err := strconv.ParseUint(queueLength, 10, 32)
	if err != nil || l > 4294967294 {
		return false
	}

	return true
}

func IsLinkMtu(value string) bool {
	if strings.HasSuffix(value, "K") || strings.HasSuffix(value, "M") || strings.HasSuffix(value, "G") {
		return true
	}
	_, err := strconv.ParseUint(value, 10, 32)
	return err == nil
}

func IsLinkBitsPerSecond(value string) bool {
	if strings.HasSuffix(value, "K") || strings.HasSuffix(value, "M") || strings.HasSuffix(value, "G") {
		return true
	}
	_, err := strconv.ParseUint(value, 10, 32)
	return err == nil
}

func IsLinkDuplex(duplex string) bool {
	return duplex == "full" || duplex == "half"
}

func IsLinkWakeOnLan(value string) bool {
	return value == "off" || value == "phy" || value == "unicast" || value == "multicast" ||
		value == "broadcast" || value == "arp" || value == "magic" || value == "secureon"
}

func IsLinkPort(port string) bool {
	return port == "tp" || port == "aui" || port == "bnc" || port == "mii" || port == "fibre"
}

func IsLinkAdvertise(advertise string) bool {
	return advertise == "10baset-half" || advertise == "10baset-full" || advertise == "100baset-half" ||
		advertise == "100baset-full" || advertise == "1000baset-half" || advertise == "1000baset-full" ||
		advertise == "10000baset-full" || advertise == "2500basex-full" || advertise == "1000basekx-full" ||
		advertise == "10000basekx4-full" || advertise == "10000basekr-full" || advertise == "10000baser-fec" ||
		advertise == "20000basemld2-full" || advertise == "20000basekr2-full"
}

func IsLinkGSO(value string) bool {
	if strings.HasSuffix(value, "K") || strings.HasSuffix(value, "M") || strings.HasSuffix(value, "G") {
		return true
	}

	l, err := strconv.ParseUint(value, 10, 32)
	if err != nil || l > 65536 {
		return false
	}

	return true
}

func IsLinkGroup(value string) bool {
	l, err := strconv.ParseUint(value, 10, 32)
	if err != nil || l > 2147483647 {
		return false
	}

	return true
}

func IsLinkRequiredFamilyForOnline(family string) bool {
	return family == "ipv4" || family == "ipv6" || family == "both" || family == "any"
}

func IsLinkActivationPolicy(policy string) bool {
	return policy == "up" || policy == "always-up" || policy == "down" ||
		policy == "always-down" || policy == "manual" || policy == "bound"
}

func LinkExists(link string) bool {
	_, err := netlink.LinkByName(link)
	return err == nil
}
