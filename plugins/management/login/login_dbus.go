package login

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/bus"
)

const (
	dbusManagerinterface = "org.freedesktop.login1.Manager"
	dbusPath             = "/org/freedesktop/login1"
	dbusInterface        = "org.freedesktop.login1"
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

func DBusAcquireUsersFromLogin(ctx context.Context) ([]User, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	out := [][]interface{}{}
	if err := c.object.Call(dbusManagerinterface+".ListUsers", 0).Store(&out); err != nil {
		return nil, err
	}

	users := []User{}
	for _, v := range out {
		u := User{
			UID:  v[0].(uint32),
			Name: fmt.Sprintf("%v", v[1]),
			Path: fmt.Sprintf("%v", v[2]),
		}

		users = append(users, u)
	}

	return users, nil
}

func DBusAcquireUSessionsFromLogin(ctx context.Context) ([]Session, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	out := [][]interface{}{}
	if err := c.object.Call(dbusManagerinterface+".ListSessions", 0).Store(&out); err != nil {
		return nil, err
	}

	sessions := []Session{}
	for _, v := range out {
		s := Session{
			ID:   fmt.Sprintf("%v", v[0]),
			UID:  v[1].(uint32),
			User: fmt.Sprintf("%v", v[2]),
			Seat: fmt.Sprintf("%v", v[3]),
			Path: fmt.Sprintf("%v", v[4]),
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}
