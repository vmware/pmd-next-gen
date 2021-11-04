// SPDX-License-Identifier: Apache-2.0

package share

import (
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

	return false, errors.New("failed to parse")
}

func ParseIp(ip string) (net.IP, error) {
	a := net.ParseIP(ip)

	if a.To4() == nil || a.To16() == nil {
		return nil, errors.New("invalid IP")
	}

	return a, nil
}

func ParsePort(port string) (uint16, error) {
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, errors.Wrap(err, "invalid port")
	}

	return uint16(p), nil
}

func ParseIpPort(s string) (string, string, error) {
	ip, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", "", err
	}

	if _, err := ParseIp(ip); err != nil {
		return "", "", errors.New("invalid IP")
	}

	if _, err := ParsePort(port); err != nil {
		return "", "", errors.New("invalid port")
	}

	return ip, port, nil
}
