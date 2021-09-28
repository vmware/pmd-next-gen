// SPDX-License-Identifier: Apache-2.0

package share

import (
	"os"
	"path"
	"strings"

	"golang.org/x/sys/unix"
)

// IsValidIfName tests whether it's a valid ifname
func IsValidIfName(ifname string) bool {
	s := strings.TrimSpace(ifname)
	if len(s) == 0 || len(s) > unix.IFNAMSIZ {
		return false
	}

	return true
}

// LinkExists tests whether link exists
func LinkExists(ifname string) bool {
	_, err := os.Stat(path.Join("/sys/class/net", ifname))
	if os.IsNotExist(err) {
		return false
	}

	return true
}
