// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/syndtr/gocapability/capability"
	"golang.org/x/sys/unix"
)

func AllowedCapabilities() ([]capability.Cap, error) {
	return capSlice([]string{
		"CAP_CHOWN",
		"CAP_SYS_ADMIN",
		"CAP_NET_ADMIN",
		"CAP_NET_BIND_SERVICE",
	})
}

var capabilityMap map[string]capability.Cap

func init() {
	capabilityMap = make(map[string]capability.Cap, capability.CAP_LAST_CAP+1)
	for _, c := range capability.List() {
		if c > capability.CAP_LAST_CAP {
			continue
		}
		capabilityMap["CAP_"+strings.ToUpper(c.String())] = c
	}
}

func capSlice(caps []string) ([]capability.Cap, error) {
	out := make([]capability.Cap, len(caps))
	for i, c := range caps {
		v, ok := capabilityMap[c]
		if !ok {
			return nil, fmt.Errorf("unknown capability %q", c)
		}
		out[i] = v
	}
	return out, nil
}

func ApplyCapability(cred *syscall.Credential) error {
	caps, err := capability.NewPid2(0)
	if err != nil {
		return err
	}

	allCapabilityTypes := capability.CAPS | capability.BOUNDS | capability.AMBS
	allowedCap, err := AllowedCapabilities()
	if err != nil {
		return err
	}

	caps.Clear(capability.CAPS | capability.BOUNDS | capability.AMBS)
	caps.Set(capability.BOUNDS, allowedCap...)
	caps.Set(capability.PERMITTED, allowedCap...)
	caps.Set(capability.INHERITABLE, allowedCap...)
	caps.Set(capability.EFFECTIVE, allowedCap...)

	caps.Clear(capability.AMBIENT)

	return caps.Apply(allCapabilityTypes)
}

func EnableKeepCapability() error {
	if err := unix.Prctl(unix.PR_SET_KEEPCAPS, 1, 0, 0, 0); err != nil {
		return err
	}

	return nil
}

func DisableKeepCapability() error {
	if err := unix.Prctl(unix.PR_SET_KEEPCAPS, 0, 0, 0, 0); err != nil {
		return err
	}

	return nil
}
