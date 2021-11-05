// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"net/http"
	"path"

	"github.com/pm-web/pkg/system"
	"github.com/pm-web/pkg/web"
)

const (
	vmPath = "/proc/sys/vm"
)

type VM struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (r *VM) GetVM(w http.ResponseWriter) error {
	line, err := system.ReadOneLineFile(path.Join(vmPath, r.Property))
	if err != nil {
		return err
	}

	vm := VM{
		Property: r.Property,
		Value:    line,
	}

	return web.JSONResponse(vm, w)
}

func (r *VM) SetVM(w http.ResponseWriter) error {
	if err := system.WriteOneLineFile(path.Join(vmPath, r.Property), r.Value); err != nil {
		return err
	}

	line, err := system.ReadOneLineFile(path.Join(vmPath, r.Property))
	if err != nil {
		return err
	}

	vm := VM{
		Property: r.Property,
		Value:    line,
	}

	return web.JSONResponse(vm, w)
}
