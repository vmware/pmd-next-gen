// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"context"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"

	"github.com/pm-web/pkg/bus"
	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/web"
)

const (
	dbusInterface = "org.freedesktop.hostname1"
	dbusPath      = "/org/freedesktop/hostname1"
)

type SDConnection struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

func NewSDConnection() (*SDConnection, error) {
	conn, err := bus.SystemBusPrivateConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	}

	return &SDConnection{
		conn:   conn,
		object: conn.Object(dbusInterface, dbus.ObjectPath(dbusPath)),
	}, nil
}

func (c *SDConnection) Close() {
	c.conn.Close()
}

func (c *SDConnection) DBusExecuteHostNameMethod(ctx context.Context, method string, value string) error {
	if err := c.object.CallWithContext(ctx, dbusInterface+"."+method, 0, value, true).Err; err != nil {
		return err
	}

	return nil
}

func (c *SDConnection) DBusHostNameDescribe(ctx context.Context) (map[string]string, error) {
	var props string

	err := c.object.CallWithContext(ctx, dbusInterface+"."+"Describe", 0).Store(&props)
	if err != nil {
		m, err := c.DBusHostNameDescribeFallback(ctx)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	msg, err := web.JSONUnmarshal([]byte(props))
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for k, v := range msg {
		if v != nil {
			m[k] = fmt.Sprintf("%v", v)
		} else {
			m[k] = ""
		}
	}

	return m, nil
}

func (c *SDConnection) DBusHostNameDescribeFallback(ctx context.Context) (map[string]string, error) {
	m := make(map[string]string)

	var wg sync.WaitGroup
	wg.Add(17)

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".StaticHostname")
		if err == nil {
			m["StaticHostname"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".Hostname")
		if err == nil {
			m["Hostname"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".PrettyHostname")
		if err == nil {
			m["PrettyHostname"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".IconName")
		if err == nil {
			m["IconName"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".Chassis")
		if err == nil {
			m["Chassis"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".Deployment")
		if err == nil {
			m["Deployment"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".Location")
		if err == nil {
			m["Location"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".KernelName")
		if err == nil {
			m["KernelName"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".KernelRelease")
		if err == nil {
			m["KernelRelease"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".KernelVersion")
		if err == nil {
			m["KernelVersion"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".OperatingSystemPrettyName")
		if err == nil {
			m["OperatingSystemPrettyName"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".OperatingSystemCPEName")
		if err == nil {
			m["OperatingSystemCPEName"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".HomeURL")
		if err == nil {
			m["HomeURL"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".HardwareVendor")
		if err == nil {
			m["HardwareVendor"] = s.Value().(string)
		}
	}()
	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".HardwareModel")
		if err == nil {
			m["HardwareModel"] = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()

		var uuid []uint8
		err := c.object.Call(dbusInterface+".GetProductUUID", 0, false).Store(&uuid)
		if err == nil {
			m["ProductUUID"] = share.BuildHexFromBytes(uuid)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.object.GetProperty(dbusInterface + ".HostnameSource")
		if err == nil {
			m["HostnameSource"] = s.Value().(string)
		}
	}()

	wg.Wait()

	return m, nil
}
