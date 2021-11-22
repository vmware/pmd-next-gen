// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"context"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/bus"
)

const (
	dbusInterface = "org.freedesktop.timesync1"
	dbusPath      = "/org/freedesktop/timesync1"

	dbusManagerinterface = "org.freedesktop.timesync1.Manager"
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

func DBusAcquireCurrentNTPServerFromTImeSync(ctx context.Context) (*NTPServer, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	var serverName dbus.Variant
	go func() {
		defer wg.Done()
		serverName, err = c.object.GetProperty(dbusManagerinterface + ".ServerName")

	}()

	var serverAddress dbus.Variant
	go func() {
		defer wg.Done()
		serverAddress, err = c.object.GetProperty(dbusManagerinterface + ".ServerAddress")
	}()

	wg.Wait()

	return &NTPServer{
		ServerName:    serverName.Value().(string),
		Family:        serverAddress.Value().([]interface{})[0].(int32),
		ServerAddress: fmt.Sprintf("%v", serverAddress.Value().([]interface{})[1]),
	}, nil
}

func DBusAcquireSystemNTPServerFromTImeSync(ctx context.Context) (*NTPServer, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	s, err := c.object.GetProperty(dbusManagerinterface + ".SystemNTPServers")
	if err != nil {
		return nil, err
	}

	return &NTPServer{
		SystemNTPServers:   s.Value().([]string),
	}, nil
}
