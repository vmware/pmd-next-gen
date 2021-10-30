// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package system

import (
	"os/user"
	"strconv"
	"syscall"
)

func GetUserCredentials(usr string) (*syscall.Credential, error) {
	var u *user.User
	var err error

	if usr != "" {
		u, err = user.Lookup(usr)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = user.Current()
		if err != nil {
			return nil, err
		}
	}

	i, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return nil, err
	}
	uid := uint32(i)

	i, err = strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return nil, err
	}
	gid := uint32(i)

	return &syscall.Credential{Uid: uid, Gid: gid}, nil
}

func SwitchUser(c *syscall.Credential) (err error) {
	if _, _, err := syscall.RawSyscall(syscall.SYS_SETGID, uintptr(c.Gid), 0, 0); err != 0 {
		return err
	}

	if _, _, err := syscall.RawSyscall(syscall.SYS_SETUID, uintptr(c.Uid), 0, 0); err != 0 {
		return err
	}

	return nil
}
