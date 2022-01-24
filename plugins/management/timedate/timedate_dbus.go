// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timedate

import (
	"fmt"

	"github.com/godbus/dbus/v5"

	"github.com/pmd-nextgen/pkg/bus"
)

const (
	dbusInterface = "org.freedesktop.timedate1"
	dbusPath      = "/org/freedesktop/timedate1"
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

func (c *SDConnection) dBusConfigureTimeDate(property string, value string) error {
	var err error

	if property == "SetNTP" {
		err = c.object.Call(dbusInterface+"."+property, 0, true, false).Err
	} else {
		err = c.object.Call(dbusInterface+"."+property, 0, value, false).Err
	}

	return err
}

func (c *SDConnection) dbusAcquire(property string) (dbus.Variant, error) {
	p, err := c.object.GetProperty(dbusInterface + "." + property)
	if err != nil {
		return dbus.Variant{}, err
	}

	return p, nil
}
