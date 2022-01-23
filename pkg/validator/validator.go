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
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return true
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return false
	}

	return false
}
func IsArrayEmpty(str []string) bool {
	return len(str) == 0
}

func IsArrayNotEmpty(str []string) bool {
	return !IsArrayEmpty(str)
}

func IsEmptyString(str string) bool {
	return govalidator.IsNull(str)
}

func IsNotEmptyString(str string) bool {
	return govalidator.IsNotNull(str)
}

func IsIP(str string) bool {
	_, _, err := net.ParseCIDR(str)
	if err != nil {
		ip := net.ParseIP(str)
		return ip != nil
	}

	return err == nil
}

func IsNotMAC(mac string) bool {
	return !govalidator.IsMAC(mac)
}

func IsMtu(mtu string) bool {
	_, err := strconv.ParseUint(mtu, 10, 32)
	return err == nil
}

func IsVLanId(id string) bool {
	_, err := strconv.ParseUint(id, 10, 32)
	return err == nil
}


