// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package validator

import (
	"net"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
)

func IsBool(str string) bool {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On", "0", "f", "F",
		"false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return true
	}

	return false
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

	if a.To4() == nil || a.To16() == nil {
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

func IsLinkQueue(id string) bool {
	l, err := strconv.ParseUint(id, 10, 32)
	if err != nil || l > 4096 {
		return false
	}

	return true
}
