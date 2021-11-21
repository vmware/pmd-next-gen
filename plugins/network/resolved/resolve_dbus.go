// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pm-web/pkg/bus"
)

const (
	dbusInterface = "org.freedesktop.resolve1"
	dbusPath      = "/org/freedesktop/resolve1"

	dbusManagerinterface = "org.freedesktop.resolve1.Manager"
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

func DBusNetworkLinkProperty(ctx context.Context) ([]DNS, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	var linkPath dbus.ObjectPath

	c.object.CallWithContext(ctx, dbusManagerinterface+".GetLink", 0, 2).Store(&linkPath)

	linkO := c.conn.Object("org.freedesktop.resolve1", linkPath)
	variant, err := linkO.GetProperty("org.freedesktop.resolve1.Link.DNS")
	if err != nil {
		return nil, fmt.Errorf("error fetching DNS property from DBus: %v", err)
	}

	var dns []DNS
	for _, value := range variant.Value().([][]interface{}) {
		d := DNS{
			DNS: fmt.Sprintf("%v", value[1]),
		}

		d.Family, _ = value[0].(int32)
		dns = append(dns, d)
	}

	return dns, nil
}

func DBusResolveManagerDNS(ctx context.Context) ([]DNS, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	variant, err := c.object.GetProperty(dbusManagerinterface + ".DNS")
	if err != nil {
		return nil, fmt.Errorf("error fetching DNS property from DBus: %v", err)
	}

	var dns []DNS
	for _, value := range variant.Value().([][]interface{}) {
		d := DNS{
			Family: value[1].(int32),
			DNS:    fmt.Sprintf("%v", value[2]),
		}

		index := value[0].(int32)
		if index != 0 {
			link, err := netlink.LinkByIndex(int(index))
			if err != nil {
				return nil, err
			}
			d.Link = link.Attrs().Name
		}

		dns = append(dns, d)
	}

	return dns, nil
}
