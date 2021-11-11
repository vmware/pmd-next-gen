// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/web"
)

type Hostname struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

var hostNameMethods share.Set

func (h *Hostname) SetHostname() error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	b := hostNameMethods.Contains(h.Property)
	if !b {
		return fmt.Errorf("failed to set hostname property: '%s' not found", h.Property)
	}

	return conn.SetHostName(h.Property, h.Value)
}

func AcquireHostnameProperties(w http.ResponseWriter) error {
	conn, err := NewConn()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	hostNameProperties := map[string]string{}

	var wg sync.WaitGroup
	wg.Add(13)

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("Hostname"); err == nil {
			hostNameProperties["Hostname"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("StaticHostname"); err == nil {
			hostNameProperties["StaticHostname"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("PrettyHostname"); err == nil {
			hostNameProperties["PrettyHostname"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("IconName"); err == nil {
			hostNameProperties["IconName"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("Chassis"); err == nil {
			hostNameProperties["Chassis"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("Deployment"); err == nil {
			hostNameProperties["Deployment"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("Location"); err == nil {
			hostNameProperties["Location"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("KernelName"); err == nil {
			hostNameProperties["KernelName"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("KernelRelease"); err == nil {
			hostNameProperties["LKernelRelease"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("KernelVersion"); err == nil {
			hostNameProperties["KernelVersion"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("OperatingSystemPrettyName"); err == nil {
			hostNameProperties["OperatingSystemPrettyName"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("OperatingSystemCPEName"); err == nil {
			hostNameProperties["OperatingSystemCPEName"] = p
		}

	}()

	go func() {
		defer wg.Done()

		if p, err := conn.GetHostName("HomeURL"); err == nil {
			hostNameProperties["HomeURL"] = p
		}

	}()

	wg.Wait()

	return web.JSONResponse(hostNameProperties, w)
}

func Init() {
	hostNameMethods := share.NewSet()

	hostNameMethods.Add("SetHostname")
	hostNameMethods.Add("SetStaticHostname")
	hostNameMethods.Add("SetPrettyHostname")
	hostNameMethods.Add("SetIconName")
	hostNameMethods.Add("SetChassis")
	hostNameMethods.Add("SetDeployment")
	hostNameMethods.Add("SetLocation")
}
