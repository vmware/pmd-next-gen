package system

import (
	"errors"
	"os/exec"
	"time"

	sdbus "github.com/coreos/go-systemd/v22/dbus"
)

func UnixMicro(usec int64) time.Time {
	return time.Unix(usec/1e6, (usec%1e6)*1e3)
}

func DBusTimeToUsec(prop *sdbus.Property) (time.Time, error) {
	var usec int64

	if err := prop.Value.Store(&usec); err != nil {
		return UnixMicro(0), err
	}

	if usec == 0 {
		return UnixMicro(0), errors.New("0")
	}

	return UnixMicro(usec), nil
}

func ExecAndCapture(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	if out, err := c.CombinedOutput(); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}
