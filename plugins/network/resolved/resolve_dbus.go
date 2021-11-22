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

func buildDNSMessage(variant dbus.Variant, link bool) ([]DNS, error) {
	var dns []DNS
	for _, v := range variant.Value().([][]interface{}) {
		d := DNS{}
		if link {
			d.Family = v[0].(int32)
			d.DNS = fmt.Sprintf("%v", v[1])
		} else {
			d.Family = v[1].(int32)
			d.DNS = fmt.Sprintf("%v", v[2])

			index := v[0].(int32)
			if index != 0 {
				link, err := netlink.LinkByIndex(int(index))
				if err != nil {
					return nil, err
				}
				d.Link = link.Attrs().Name
			}
		}

		dns = append(dns, d)
	}

	return dns, nil
}

func buildDomainsMessage(variant dbus.Variant) ([]Domains, error) {
	var domains []Domains
	for _, v := range variant.Value().([][]interface{}) {
		d := Domains{
			Domain: fmt.Sprintf("%v", v[1]),
		}

		index := v[0].(int32)
		if index != 0 {
			link, err := netlink.LinkByIndex(int(index))
			if err != nil {
				return nil, err
			}
			d.Link = link.Attrs().Name
		}

		domains = append(domains, d)
	}

	return domains, nil
}

func DBusAcquireDNSFromResolveLink(ctx context.Context, link string) ([]DNS, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	l, err := netlink.LinkByName(link)
	if err != nil {
		return nil, err
	}

	var linkPath dbus.ObjectPath

	c.object.CallWithContext(ctx, dbusManagerinterface+".GetLink", 0, l.Attrs().Index).Store(&linkPath)
	variant, err := c.conn.Object("org.freedesktop.resolve1", linkPath).GetProperty("org.freedesktop.resolve1.Link.DNS")
	if err != nil {
		return nil, fmt.Errorf("error fetching DNS from resolve: %v", err)
	}

	return buildDNSMessage(variant, true)
}

func DBusAcquireDomainsFromResolveLink(ctx context.Context, link string) ([]Domains, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	l, err := netlink.LinkByName(link)
	if err != nil {
		return nil, err
	}

	var linkPath dbus.ObjectPath

	c.object.CallWithContext(ctx, dbusManagerinterface+".GetLink", 0, l.Attrs().Index).Store(&linkPath)
	variant, err := c.conn.Object("org.freedesktop.resolve1", linkPath).GetProperty("org.freedesktop.resolve1.Link.Domains")
	if err != nil {
		return nil, fmt.Errorf("error fetching Domains from resolve: %v", err)
	}

	return buildDomainsMessage(variant)
}

func DBusAcquireDNSFromResolveManager(ctx context.Context) ([]DNS, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	variant, err := c.object.GetProperty(dbusManagerinterface + ".DNS")
	if err != nil {
		return nil, fmt.Errorf("error fetching DNS from resolve: %v", err)
	}

	return buildDNSMessage(variant, false)
}

func DBusAcquireDomainsFromResolveManager(ctx context.Context) ([]Domains, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	variant, err := c.object.GetProperty(dbusManagerinterface + ".Domains")
	if err != nil {
		return nil, fmt.Errorf("error fetching Domains from resolve: %v", err)
	}

	return buildDomainsMessage(variant)
}
