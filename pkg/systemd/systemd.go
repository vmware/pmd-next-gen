// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"context"
	"encoding/json"
	"fmt"
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

// State sytemd state
func State(w http.ResponseWriter) error {
	v, err := getProperty("SystemState")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "SystemState",
		Value:    v.Value().(string),
	}

	return web.JSONResponse(prop, w)
}

// Version systemd version
func Version(w http.ResponseWriter) error {
	v, err := getProperty("Version")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Version",
		Value:    v.Value().(string),
	}

	return web.JSONResponse(prop, w)
}

// Virtualization systemd virt
func Virtualization(w http.ResponseWriter) error {
	v, err := getProperty("Virtualization")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Virtualization",
		Value:    v.Value().(string),
	}

	return web.JSONResponse(prop, w)
}

// Architecture arch of the system
func Architecture(w http.ResponseWriter) error {
	v, err := getProperty("Architecture")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Architecture",
		Value:    v.Value().(string),
	}

	return web.JSONResponse(prop, w)
}

// Features systemd features
func Features(w http.ResponseWriter) error {
	v, err := getProperty("Features")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "Features",
		Value:    v.Value().(string),
	}

	return web.JSONResponse(prop, w)
}

// NFailedUnits how many uniuts failed
func NFailedUnits(w http.ResponseWriter) error {
	v, err := getProperty("NFailedUnits")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "NFailedUnits",
		Value:    fmt.Sprint(v.Value().(uint32)),
	}

	return web.JSONResponse(prop, w)
}

// NNames number of names
func NNames(w http.ResponseWriter) error {
	v, err := getProperty("NNames")
	if err != nil {
		return err
	}

	prop := Property{
		Property: "NNames",
		Value:    fmt.Sprint(v.Value().(uint32)),
	}

	return web.JSONResponse(prop, w)
}

// ListUnits list all units
func ListUnits(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsContext(context.Background())
	if err != nil {
		log.Errorf("Failed ListUnits: %v", err)
		return err
	}

	return web.JSONResponse(units, w)
}

// StartUnit start a unit
func (u *Unit) StartUnit() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StartUnitContext(context.Background(), u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to start unit %s: %v", u.Unit, err)
		return err
	}

	return nil
}

// StopUnit stop a unit
func (u *Unit) StopUnit() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StopUnitContext(context.Background(), u.Unit, "fail", reschan)
	if err != nil {
		log.Errorf("Failed to stop unit %s: %v", u.Unit, err)
		return err
	}

	return nil
}

// RestartUnit restart a unit
func (u *Unit) RestartUnit() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.RestartUnitContext(context.Background(), u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to restart unit %s: %v", u.Unit, err)
		return err
	}

	return nil
}

// ReloadUnit reload daemon
func (u *Unit) ReloadUnit() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	err = conn.ReloadContext(context.Background())
	if err != nil {
		log.Errorf("Failed to reload unit %s: %v", u.Unit, err)
		return err
	}

	return nil
}

// KillUnit send a signal to a unit
func (u *Unit) KillUnit() error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	signal, err := strconv.ParseInt(u.Value, 10, 64)
	if err != nil {
		log.Errorf("Failed to parse signal number '%s': %s", u.Value, err)
		return err
	}

	conn.KillUnitContext(context.Background(), u.Unit, int32(signal))

	return nil
}

// GetUnitStatus get unit status
func (u *Unit) GetUnitStatus(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
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
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	if u.Property != "" {
		p, err := conn.GetServicePropertyContext(context.Background(), u.Unit, u.Property)
		if err != nil {
			log.Errorf("Failed to get service property: %v", err)
			return err
		}

		switch u.Property {
		case "CPUShares", "LimitNOFILE", "LimitNOFILESoft":
			cpu := strconv.FormatUint(p.Value.Value().(uint64), 10)
			prop := Property{Property: p.Name, Value: cpu}

			return web.JSONResponse(prop, w)
		}
	}

	p, err := conn.GetUnitPropertiesContext(context.Background(), u.Unit)
	if err != nil {
		log.Errorf("Failed to get service properties: %v", err)
		return err
	}

	return web.JSONResponse(p, w)
}

// GetUnitTypeProperty get unit type property
func (u *Unit) GetUnitTypeProperty(w http.ResponseWriter) error {
	conn, err := sd.NewSystemdConnectionContext(context.Background())
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
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
