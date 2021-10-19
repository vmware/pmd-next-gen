// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"errors"
	"net/http"
	"path"

	"github.com/pmd/pkg/system"
	"github.com/pmd/pkg/web"
)

const (
	sysNetPath     = "/proc/sys/net"
	sysNetPathCore = "core"
	sysNetPathIPv4 = "ipv4"
	sysNetPathIPv6 = "ipv6"
)

// SysNet Json request
type SysNet struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
	Link     string `json:"link"`
}

// getPath read info from proc
func (r *SysNet) getPath() (string, error) {
	var procPath string

	switch r.Path {
	case sysNetPathCore:
		procPath = path.Join(path.Join(sysNetPath, sysNetPathCore), r.Property)

	case sysNetPathIPv4:

		if r.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(sysNetPath, sysNetPathIPv4), "conf"), r.Link), r.Property)
		} else {
			procPath = path.Join(path.Join(sysNetPath, sysNetPathIPv4), r.Property)
		}

	case sysNetPathIPv6:

		if r.Link != "" {
			procPath = path.Join(path.Join(path.Join(path.Join(sysNetPath, sysNetPathIPv6), "conf"), r.Link), r.Property)
		} else {
			procPath = path.Join(path.Join(sysNetPath, sysNetPathIPv6), r.Property)
		}

	default:
		return "", errors.New("path not found")
	}

	return procPath, nil
}

// GetSysNet read proc value and send response
func (r *SysNet) GetSysNet(rw http.ResponseWriter) error {
	path, err := r.getPath()
	if err != nil {
		return err
	}

	line, err := system.ReadOneLineFile(path)
	if err != nil {
		return err
	}

	s := SysNet{
		Path:     r.Path,
		Property: r.Property,
		Value:    line,
		Link:     r.Link,
	}

	return web.JSONResponse(s, rw)
}

// SetSysNet sets a value to proc
func (r *SysNet) SetSysNet(rw http.ResponseWriter) error {
	path, err := r.getPath()
	if err != nil {
		return err
	}

	return system.WriteOneLineFile(path, r.Value)
}
