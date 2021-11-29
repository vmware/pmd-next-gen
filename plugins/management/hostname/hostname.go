// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/web"
)

type Hostname struct {
	Method string `json:"Method"`
	Value  string `json:"Value"`
}

type Describe struct {
	Chassis                   string `json:"Chassis"`
	DefaultHostname           string `json:"DefaultHostname"`
	Deployment                string `json:"Deployment"`
	HardwareModel             string `json:"HardwareModel"`
	HardwareVendor            string `json:"HardwareVendor"`
	Hostname                  string `json:"Hostname"`
	HostnameSource            string `json:"HostnameSource"`
	IconName                  string `json:"IconName"`
	KernelName                string `json:"KernelName"`
	KernelRelease             string `json:"KernelRelease"`
	KernelVersion             string `json:"KernelVersion"`
	Location                  string `json:"Location"`
	OperatingSystemCPEName    string `json:"OperatingSystemCPEName"`
	OperatingSystemHomeURL    string `json:"OperatingSystemHomeURL"`
	OperatingSystemPrettyName string `json:"OperatingSystemPrettyName"`
	PrettyHostname            string `json:"PrettyHostname"`
	ProductUUID               string `json:"ProductUUID"`
	StaticHostname            string `json:"StaticHostname"`
}

func (h *Hostname) SetHostname(ctx context.Context, w http.ResponseWriter) error {
	conn, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	if err := conn.DBusExecuteHostNameMethod(ctx, h.Method, h.Value); err != nil {
		return err
	}

	return web.JSONResponse("hostname set to: " + h.Value , w)
}

func HostnameDescribe(ctx context.Context, w http.ResponseWriter) error {
	conn, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	desc, err := conn.DBusHostNameDescribe(ctx)
	if err != nil {
		return err
	}

	return web.JSONResponse(desc, w)

}
