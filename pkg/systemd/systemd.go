// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"context"
	"encoding/json"
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

// Property generic property and value
type Property struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

// UnitStatus unit status
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

// ListUnits list all units
func ListUnits(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		log.Errorf("Failed list units: %v", err)
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
		_, err = conn.StartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to start unit '%s': %v", u.Unit, err)
			return err
		}
	case "stop":
		_, err = conn.StopUnitContext(context.Background(), u.Unit, "fail", c)
		if err != nil {
			log.Errorf("Failed to stop unit '%s': %v", u.Unit, err)
			return err
		}
	case "restart":
		_, err = conn.RestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to restart unit '%s': %v", u.Unit, err)
			return err
		}

	case "try-restart":
		_, err = conn.TryRestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to try restart unit '%s': %v", u.Unit, err)
			return err
		}

	case "reload-or-restart":
		_, err = conn.ReloadOrRestartUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload or restart unit '%s': %v", u.Unit, err)
			return err
		}

	case "reload":
		_, err = conn.ReloadUnitContext(context.Background(), u.Unit, "replace", c)
		if err != nil {
			log.Errorf("Failed to reload unit '%s': %v", u.Unit, err)
			return err
		}

	case "kill":
		signal, err := strconv.ParseInt(u.Value, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse signal number '%s': %s", u.Value, err)
			return err
		}

		conn.KillUnitContext(context.Background(), u.Unit, int32(signal))
	}

	return nil
}

// GetUnitStatus get unit status
func (u *Unit) GetUnitStatus(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus:: %v", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsByNamesContext(context.Background(), []string{u.Unit})
	if err != nil {
		log.Errorf("Failed get unit '%s' status: %v", u.Unit, err)
		return err
	}

	status := UnitStatus{
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

	json.NewEncoder(w).Encode(status)

	return nil
}

// GetUnitProperty get unit property
func (u *Unit) GetUnitProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus: %v", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetUnitPropertiesContext(context.Background(), u.Unit)
	if err != nil {
		log.Errorf("Failed to get unit properties: %v", err)
		return err
	}

	return web.JSONResponse(p, w)
}

// GetUnitTypeProperty get unit type property
func (u *Unit) GetUnitTypeProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to establishes connection to the system bus:: %v", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetUnitTypePropertiesContext(context.Background(), u.Unit, u.UnitType)
	if err != nil {
		log.Errorf("Failed to get unit type properties: %v", err)
		return err
	}

	return web.JSONResponse(p, w)
}
