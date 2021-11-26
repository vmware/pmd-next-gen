package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/pm-web/pkg/web"
	"github.com/pm-web/plugins/management/hostname"
	"github.com/pm-web/plugins/systemd"
)

type Hostname struct {
	Success bool              `json:"success"`
	Message hostname.Describe `json:"message"`
	Errors  string            `json:"errors"`
}

type Systemd struct {
	Success bool             `json:"success"`
	Message systemd.Describe `json:"message"`
	Errors  string           `json:"errors"`
}

func acquireHostname(host string, token map[string]string) (*hostname.Describe, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/system/hostname/describe", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/system/hostname/describe", nil)
	}
	if err != nil {
		fmt.Printf("Failed to fetch hostname: %v\n", err)
		return nil, err
	}

	h := Hostname{}
	if err := json.Unmarshal(resp, &h); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !h.Success {
		fmt.Printf("%v\n", h.Errors)
		return nil, errors.New(h.Errors)
	}

	return &h.Message, nil
}

func acquireSystemd(host string, token map[string]string) (*systemd.Describe, error) {
	var resp []byte
	var err error

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/service/systemd/manager/describe", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/service/systemd/manager/describe", nil)
	}
	if err != nil {
		fmt.Printf("Failed to fetch hostname: %v\n", err)
		return nil, err
	}

	sd := Systemd{}
	if err := json.Unmarshal(resp, &sd); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return nil, err
	}

	if !sd.Success {
		fmt.Printf("%v\n", sd.Errors)
		return nil, errors.New(sd.Errors)
	}

	return &sd.Message, nil
}

func displayHostname(h *hostname.Describe) {
	fmt.Printf("              %v %v\n", color.HiBlueString("System Name:"), h.StaticHostname)
	fmt.Printf("                   %v %v (%v) %v\n", color.HiBlueString("Kernel:"), h.KernelName, h.KernelRelease, h.KernelVersion)
	fmt.Printf("                  %v %v\n", color.HiBlueString("Chassis:"), h.Chassis)
	if h.HardwareModel != "" {
		fmt.Printf("           %v %v\n", color.HiBlueString("Hardware Model:"), h.HardwareModel)
	}
	if h.HardwareVendor != "" {
		fmt.Printf("          %v %v\n", color.HiBlueString("Hardware Vendor:"), h.HardwareVendor)
	}
	if h.ProductUUID != "" {
		fmt.Printf("             %v %v\n", color.HiBlueString("Product UUID:"), h.ProductUUID)
	}
	fmt.Printf("         %v %v\n", color.HiBlueString("Operating System:"), h.OperatingSystemPrettyName)
	if h.OperatingSystemHomeURL != "" {
		fmt.Printf("%v %v\n", color.HiBlueString("Operating System Home URL:"), h.OperatingSystemHomeURL)
	}
}

func displaySystemd(sd *systemd.Describe) {
	fmt.Printf("          %v %v\n", color.HiBlueString("Systemd Version:"), sd.Version)
	fmt.Printf("          %v %v\n", color.HiBlueString("   Architecture:"), sd.Architecture)
	fmt.Printf("          %v %v\n", color.HiBlueString(" Virtualization:"), sd.Virtualization)

}
func acquireSystemStatus(host string, token map[string]string) {
	h, err := acquireHostname(host, token)
	if err != nil {
		return
	}

	sd, err := acquireSystemd(host, token)
	if err != nil {
		return
	}

	displayHostname(h)
	displaySystemd(sd)
}
