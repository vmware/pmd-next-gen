// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/bus"
)

const (
	dbusInterface = "org.freedesktop.network1"
	dbusPath      = "/org/freedesktop/network1"

	dbusManagerinterface = "org.freedesktop.network1.Manager"
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

func DBusNetworkReconfigureLink(ctx context.Context, index int) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	if err := c.object.CallWithContext(ctx, dbusManagerinterface+"."+"ReconfigureLink", 0, index).Err; err != nil {
		return err
	}

	return nil
}

func DBusNetworkReload(ctx context.Context) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	if err := c.object.CallWithContext(ctx, dbusManagerinterface+"."+"Reload", 0).Err; err != nil {
		return err
	}

	return nil
}

func DBusNetworkLinkProperty(ctx context.Context) (map[string]interface{}, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	var props string
	err = c.object.CallWithContext(ctx, dbusManagerinterface+"."+"Describe", 0).Store(&props)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(props), &m); err != nil {
		return nil, err
	}

	return m, nil
}
