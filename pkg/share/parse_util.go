// SPDX-License-Identifier: Apache-2.0

package share

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParseBool extends ParseBool
func ParseBool(str string) (bool, error) {
	b, err := strconv.ParseBool(str)
	if err == nil {
		return b, err
	}

	if strings.EqualFold(str, "yes") || strings.EqualFold(str, "y") || strings.EqualFold(str, "on") {
		return true, nil
	} else if strings.EqualFold(str, "no") || strings.EqualFold(str, "n") || strings.EqualFold(str, "off") {
		return false, nil
	}

	return false, fmt.Errorf("ParseBool")
}

// ParseIP parse a IP
func ParseIP(ip string) (net.IP, error) {
	if len(ip) == 0 {
		return nil, fmt.Errorf("ParseIP")
	}

	a := net.ParseIP(ip)

	if a.To4() == nil || a.To16() == nil {
		return nil, fmt.Errorf("ParseIP")
	}

	return a, nil
}

// ParsePort parse a port
func ParsePort(port string) (uint16, error) {
	if len(port) == 0 {
		return 0, nil
	}

	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(p), nil
}
