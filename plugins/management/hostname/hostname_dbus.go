// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"

	"github.com/pm-web/pkg/web"
	"github.com/pm-web/pkg/bus"
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

func (c *SDConnection) ExecuteHostNameMethod(ctx context.Context, method string, value string) error {
	if err := c.object.CallWithContext(ctx, dbusInterface+"."+method, 0, value, false).Err; err != nil {
		return err
	}

	return nil
}

func (c *SDConnection) AcquireHostNameProperty(ctx context.Context, property string) (map[string]string, error) {
	var props string

	err := c.object.CallWithContext(ctx, dbusInterface+"."+"Describe", 0).Store(&props)
	if err != nil {
		return nil, err
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
