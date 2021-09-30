// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func routerGetSystemdState(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := State(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetSystemdVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := Version(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetSystemdFeatures(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := Features(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetSystemdVirtualization(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := Virtualization(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetSystemdNFailedUnits(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := NFailedUnits(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetSystemdNNames(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := NNames(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

	}
}

func routerGetSystemdArchitecture(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := Architecture(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerConfigureSystemdConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetSystemConf(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		if err := UpdateSystemConf(rw, r); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerConfigureUnit(rw http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case "POST":
		unit := new(Unit)

		err = json.NewDecoder(r.Body).Decode(&unit)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		switch unit.Action {
		case "start":
			err = unit.StartUnit()

		case "stop":
			err = unit.StopUnit()

		case "restart":
			err = unit.RestartUnit()

		case "reload":
			err = unit.ReloadUnit()

		case "kill":
			err = unit.KillUnit()

		}
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func routerGetAllSystemdUnits(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := ListUnits(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetUnitStatus(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]

	u := Unit{
		Unit: unit,
	}

	switch r.Method {
	case "GET":
		if err := u.GetUnitStatus(rw); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func routerGetUnitProperty(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]
	property := vars["property"]

	u := Unit{
		Unit:     unit,
		Property: property,
	}

	switch r.Method {
	case "GET":
		u.GetUnitProperty(rw)
	}
}

func routerGetUnitTypeProperty(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unit := vars["unit"]
	unitType := vars["unittype"]
	property := vars["property"]

	u := Unit{
		Unit:     unit,
		UnitType: unitType,
		Property: property,
	}

	switch r.Method {
	case "GET":
		u.GetUnitTypeProperty(rw)
	}
}

// RegisterRouterSystemd register with mux
func RegisterRouterSystemd(router *mux.Router) {
	n := router.PathPrefix("/service").Subrouter()

	// property
	n.HandleFunc("/systemd/state", routerGetSystemdState)
	n.HandleFunc("/systemd/version", routerGetSystemdVersion)
	n.HandleFunc("/systemd/features", routerGetSystemdFeatures)
	n.HandleFunc("/systemd/virtualization", routerGetSystemdVirtualization)
	n.HandleFunc("/systemd/architecture", routerGetSystemdArchitecture)
	n.HandleFunc("/systemd/units", routerGetAllSystemdUnits)
	n.HandleFunc("/systemd/nnames", routerGetSystemdNNames)
	n.HandleFunc("/systemd/nfailedunits", routerGetSystemdNFailedUnits)

	// unit
	n.HandleFunc("/systemd", routerConfigureUnit)
	n.HandleFunc("/systemd/{unit}/status", routerGetUnitStatus)
	n.HandleFunc("/systemd/{unit}/property", routerGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/property/{property}", routerGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/property/{unittype}", routerGetUnitTypeProperty)

	// conf
	n.HandleFunc("/systemd/conf", routerConfigureSystemdConf)
	n.HandleFunc("/systemd/conf/update", routerConfigureSystemdConf)
}
