// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"fmt"

	"github.com/godbus/dbus/v5"

	"github.com/pm-web/pkg/share"
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
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	}

	c := SDConnection {
		conn : conn,
		object :conn.Object(dbusInterface, dbus.ObjectPath(dbusPath)),
	}

	return &c, nil
}

func (c *SDConnection) Close() {
	c.conn.Close()
}

func (c *SDConnection) SetHostName(method string, value string) error {
	if err := c.object.Call(dbusInterface+"."+method, 0, value, false).Err; err != nil {
		return err
	}

	return nil
}

func (c *SDConnection) GetHostName(property string) (string, error) {
	p, err := c.object.GetProperty(dbusInterface + "." + property)
	if err != nil {
		return "", err
	}

	v, b := p.Value().(string)
	if !b {
		return "", fmt.Errorf("empty value received: %s", property)
	}

	return v, nil
}
