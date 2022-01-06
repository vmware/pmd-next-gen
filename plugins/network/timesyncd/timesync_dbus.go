// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timesyncd

import (
	"context"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"

	"github.com/pmd-nextgen/pkg/bus"
	"github.com/pmd-nextgen/pkg/share"
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

func (c *SDConnection) DBusAcquireCurrentNTPServerFromTimeSync(ctx context.Context) (*NTPServer, error) {
	var wg sync.WaitGroup
	var err error
	wg.Add(2)

	var serverName dbus.Variant
	go func() {
		defer wg.Done()
		serverName, err = c.object.GetProperty(dbusManagerinterface + ".ServerName")
		if err != nil {
			log.Errorf("Failed to acquire 'ServerName': %v", err)
		}
	}()

	var serverAddress dbus.Variant
	go func() {
		defer wg.Done()
		serverAddress, err = c.object.GetProperty(dbusManagerinterface + ".ServerAddress")
		if err != nil {
			log.Errorf("Failed to acquire 'ServerAddress': %v", err)
		}
	}()

	wg.Wait()

	return &NTPServer{
		ServerName:    serverName.Value().(string),
		Family:        serverAddress.Value().([]interface{})[0].(int32),
		ServerAddress: share.BuildIPFromBytes(serverAddress.Value().([]interface{})[1].([]uint8)),
	}, nil
}

func (c *SDConnection) DBusAcquireSystemNTPServerFromTimeSync(ctx context.Context) (*NTPServer, error) {
	s, err := c.object.GetProperty(dbusManagerinterface + ".SystemNTPServers")
	if err != nil {
		return nil, err
	}

	return &NTPServer{
		SystemNTPServers: s.Value().([]string),
	}, nil
}

func (c *SDConnection) DBusAcquireLinkNTPServerFromTimeSync(ctx context.Context) (*NTPServer, error) {
	s, err := c.object.GetProperty(dbusManagerinterface + ".LinkNTPServers")
	if err != nil {
		return nil, err
	}

	return &NTPServer{
		LinkNTPServers: s.Value().([]string),
	}, nil
}
