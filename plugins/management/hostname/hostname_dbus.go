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

type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

func NewConn() (*Conn, error) {
	c := new(Conn)

	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	}

	c.conn = conn
	c.object = conn.Object(dbusInterface, dbus.ObjectPath(dbusPath))

	return c, nil
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) SetHostName(property string, value string) error {
	if err := c.object.Call(dbusInterface+"."+property, 0, value, false).Err; err != nil {
		return err
	}

	return nil
}

func (c *Conn) GetHostName(property string) (string, error) {
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
