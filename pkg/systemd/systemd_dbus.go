// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"fmt"

	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/share"
)

const (
	dbusInterface = "org.freedesktop.systemd1"
	dbusPath      = "/org/freedesktop/systemd1"
)

// getProperty Retrive property from systemd
func getProperty(property string) (dbus.Variant, error) {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Errorf("Failed to get dbus connection: %v", err)
		return dbus.Variant{}, err
	}
	defer conn.Close()

	c := conn.Object(dbusInterface, dbusPath)
	p, perr := c.GetProperty(dbusInterface + ".Manager." + property)
	if perr != nil {
		log.Errorf("Failed to get property '%s' from systemd: %v ", property, perr)
		return dbus.Variant{}, fmt.Errorf("Failed to get dbus property: %v", perr)
	}

	if p.Value() == nil {
		return dbus.Variant{}, fmt.Errorf("Unexpected value received: %s", property)
	}

	return p, nil
}
