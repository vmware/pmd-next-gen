// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	sd "github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/web"
)

// Unit JSON message
type Unit struct {
	Action   string `json:"action"`
	Unit     string `json:"unit"`
	UnitType string `json:"unit_type"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

type Property struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

type UnitStatus struct {
	Status      string `json:"property"`
	Unit        string `json:"unit"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	LoadState   string `json:"LoadState"`
	ActiveState string `json:"pActiveState"`
	SubState    string `json:"SubState"`
	Followed    string `json:"Followed"`
	Path        string `json:"Path"`
	JobId       uint32 `json:"JobId"`
	JobType     string `json:"JobType"`
	JobPath     string `json:"JobPath"`
}

func ManagerFetchSystemProperty(w http.ResponseWriter, property string) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	v, err := conn.GetManagerProperty(property)
	if err != nil {
		log.Errorf("Failed fetch systemd manager property='%s': %v", property, err)
		return err
	}

	s, err := strconv.Unquote(string(v))
	if err != nil {
		log.Errorf("Failed to unquote systemd manager property='%s`: %v", property, err)
		return err
	}

	p := Property{
		Property: property,
		Value:    s,
	}
	return web.JSONResponse(p, w)
}

func ListUnits(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		log.Errorf("Failed list systemd units: %v", err)
		return err
	}

	return web.JSONResponse(units, w)
}

func (u *Unit) UnitActions() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %v", err)
		return err
	}
	defer conn.Close()

	c := make(chan string)
	switch u.Action {
	case "start":
		jid, err := conn.StartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to start systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully started systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "stop":
		jid, err := conn.StopUnitContext(context.Background(), u.Unit, "fail", c)
		if err != nil {
			log.Errorf("Failed to stop systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully stopped systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "restart":
		jid, err := conn.RestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to restart systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully restared systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "try-restart":
		jid, err := conn.TryRestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to try restart systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully try-restart systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "reload-or-restart":
		jid, err := conn.ReloadOrRestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload or restart systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully reload-or-restart systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "reload":
		jid, err := conn.ReloadUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload systemd unit='%s': %v", u.Unit, err)
			return err
		}

		log.Debugf("Successfully reload systemd unit='%s' job_id='%d'", u.Unit, jid)

	case "enable":
		install, changes, err := conn.EnableUnitFilesContext(context.Background(), []string{u.Unit}, false, true)
		if err != nil {
			log.Errorf("Failed to enable systemd unit='%s': %v", u.Value, err)
			return err
		}

		log.Debugf("Successfully enabled systemd unit='%s' install='%t' changes='%s'", u.Unit, install, changes)

	case "disable":
		changes, err := conn.DisableUnitFilesContext(context.Background(), []string{u.Unit}, false)
		if err != nil {
			log.Errorf("Failed to disable systemd unit='%s': %v", u.Value, err)
			return err
		}

		log.Debugf("Successfully disabled systemd unit='%s' changes='%s'", u.Unit, changes)

	case "mask":
		changes, err := conn.MaskUnitFilesContext(context.Background(), []string{u.Unit}, false, true)
		if err != nil {
			log.Errorf("Failed to mask systemd unit='%s': %v", u.Value, err)
			return err
		}

		log.Debugf("Successfully masked systemd unit='%s' changes='%s'", u.Unit, changes)

	case "unmask":
		changes, err := conn.UnmaskUnitFilesContext(context.Background(), []string{u.Unit}, false)
		if err != nil {
			log.Errorf("Failed to unmask systemd unit='%s': %v", u.Value, err)
			return err
		}

		log.Debugf("Successfully unmasked systemd unit='%s' changes='%s'", u.Unit, changes)

	case "kill":
		signal, err := strconv.ParseInt(u.Value, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse signal number='%s': %s", u.Value, err)
			return err
		}

		conn.KillUnitContext(context.Background(), u.Unit, int32(signal))

	default:
		log.Errorf("Unknown action='%s' for systemd unit='%s'", u.Action, u.Unit)
		return errors.New("unknown action")
	}

	return nil
}

func (u *Unit) GetUnitStatus(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus:: %v", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsByNamesContext(context.Background(), []string{u.Unit})
	if err != nil {
		log.Errorf("Failed fetch systemd unit='%s' status: %v", u.Unit, err)
		return err
	}

	s := UnitStatus{
		Unit:        u.Unit,
		Status:      units[0].ActiveState,
		LoadState:   units[0].LoadState,
		Name:        units[0].Name,
		Description: units[0].Description,
		ActiveState: units[0].ActiveState,
		SubState:    units[0].SubState,
		Followed:    units[0].Followed,
		Path:        string(units[0].Path),
		JobType:     units[0].JobType,
		JobPath:     string(units[0].JobPath),
	}

	return json.NewEncoder(w).Encode(s)
}

func (u *Unit) GetUnitProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %v", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetUnitPropertiesContext(context.Background(), u.Unit)
	if err != nil {
		log.Errorf("Failed to fetch systemd unit='%s' properties: %v", u.Unit, err)
		return err
	}

	return web.JSONResponse(p, w)
}

func (u *Unit) GetUnitTypeProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus:: %v", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetUnitTypePropertiesContext(context.Background(), u.Unit, u.UnitType)
	if err != nil {
		log.Errorf("Failed to fetch unit type properties unit='%s': %v", u.Unit, err)
		return err
	}

	return web.JSONResponse(p, w)
}
