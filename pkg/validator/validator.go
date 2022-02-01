// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package validator

import (
	"net"
	"strconv"

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

func IsScope(id string) bool {
	scope, err := strconv.ParseUint(id, 10, 32)
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
	if mode == "balance-rr" || mode == "active-backup" || mode == "balance-xor" ||
		mode == "broadcast" || mode == "802.3ad" || mode == "balance-tlb" || mode == "balance-alb" {
		return true
	} else {
		return false
	}
}

func IsBondTransmitHashPolicy(mode, thp string) bool {
	if (thp == "layer2" || thp == "layer3+4" || thp == "layer2+3" || thp == "encap2+3" || thp == "encap3+4") &&
		(mode == "balance-xor" || mode == "802.3ad" || mode == "balance-tlb") {
		return true
	} else {
		return false
	}
}

func IsBondLACPTransmitRate(ltr string) bool {
	if ltr == "slow" || ltr == "fast" {
		return true
	} else {
		return false
	}
}
