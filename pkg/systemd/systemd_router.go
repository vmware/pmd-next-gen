// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd/pkg/web"
)

func routerGetSystemdManagerProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	switch r.Method {
	case "GET":
		if err := ManagerFetchSystemProperty(w, v["property"]); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerConfigureSystemdConf(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetSystemConf(w); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	case "POST":
		if err := UpdateSystemConf(w, r); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerConfigureUnit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		u := new(Unit)

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := u.UnitActions(); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}

	web.JSONResponse("", w)
}

func routerGetAllSystemdUnits(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := ListUnits(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetUnitStatus(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit: v["unit"],
	}

	switch r.Method {
	case "GET":
		if err := u.GetUnitStatus(w); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}
}

func routerGetUnitProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit:     v["unit"],
		Property: v["property"],
	}

	switch r.Method {
	case "GET":
		if err := u.GetUnitProperty(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetUnitTypeProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit:     v["unit"],
		UnitType: v["unittype"],
		Property: v["property"],
	}

	switch r.Method {
	case "GET":
		u.GetUnitTypeProperty(w)
	}
}

func RegisterRouterSystemd(router *mux.Router) {
	n := router.PathPrefix("/service").Subrouter()

	n.HandleFunc("/systemd/manager/property/{property}", routerGetSystemdManagerProperty)

	n.HandleFunc("/systemd/units", routerGetAllSystemdUnits)
	n.HandleFunc("/systemd", routerConfigureUnit)
	n.HandleFunc("/systemd/{unit}/status", routerGetUnitStatus)
	n.HandleFunc("/systemd/{unit}/property", routerGetUnitProperty)
	n.HandleFunc("/systemd/{unit}/property/{unittype}", routerGetUnitTypeProperty)

	n.HandleFunc("/systemd/conf", routerConfigureSystemdConf)
	n.HandleFunc("/systemd/conf/update", routerConfigureSystemdConf)
}
