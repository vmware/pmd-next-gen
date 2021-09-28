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
		err := State(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdVersion(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Version(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdFeatures(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Features(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdVirtualization(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Virtualization(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdNFailedUnits(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := NFailedUnits(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdNNames(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := NNames(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerGetSystemdArchitecture(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Architecture(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
	}
}

func routerConfigureSystemdConf(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetSystemConf(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break

	case "POST":
		err := UpdateSystemConf(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		break
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
			break
		case "stop":
			err = unit.StopUnit()
			break
		case "restart":
			err = unit.RestartUnit()
			break
		case "reload":
			err = unit.ReloadUnit()
			break
		case "kill":
			err = unit.KillUnit()
			break
		}
		break
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func routerGetAllSystemdUnits(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := ListUnits(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break
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
		err := u.GetUnitStatus(rw)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		break
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
		break
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
		break
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
	n.HandleFunc("/systemd/{unit}/get", routerGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/get/{property}", routerGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/gettype/{unittype}", routerGetUnitTypeProperty)

	// conf
	n.HandleFunc("/systemd/conf", routerConfigureSystemdConf)
	n.HandleFunc("/systemd/conf/update", routerConfigureSystemdConf)
}
