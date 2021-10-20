// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd/pkg/web"
)

func routerFetchSystemdManagerProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	if err := ManagerFetchSystemProperty(r.Context(), w, v["property"]); err != nil {
		web.JSONResponseError(err, w)
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

		if err := u.UnitActions(r.Context()); err != nil {
			web.JSONResponseError(err, w)
			return
		}
	}

	web.JSONResponse("", w)
}

func routerFetchAllSystemdUnits(w http.ResponseWriter, r *http.Request) {

	if err := ListUnits(r.Context(),w); err != nil {
		web.JSONResponseError(err, w)
	}

}

func routerFetchUnitStatus(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit: v["unit"],
	}

	if err := u.FetchUnitStatus(r.Context(),w); err != nil {
		web.JSONResponseError(err, w)
	}

}

func routerFetchUnitProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit:     v["unit"],
		Property: v["property"],
	}

	if err := u.FetchUnitProperty(r.Context(),w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchUnitPropertyAll(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit: v["unit"],
	}

	if err := u.FetchAllUnitProperty(r.Context(),w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchUnitTypeProperty(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := Unit{
		Unit:     v["unit"],
		UnitType: v["unittype"],
		Property: v["property"],
	}

	u.GetUnitTypeProperty(r.Context(), w)
}

func RegisterRouterSystemd(router *mux.Router) {
	n := router.PathPrefix("/service").Subrouter()

	// systemd unit commands
	n.HandleFunc("/systemd", routerConfigureUnit).Methods("POST")

	// systemd unit status and property
	n.HandleFunc("/systemd/manager/property/{property}", routerFetchSystemdManagerProperty).Methods("GET")
	n.HandleFunc("/systemd/units", routerFetchAllSystemdUnits).Methods("GET")
	n.HandleFunc("/systemd/{unit}/status", routerFetchUnitStatus).Methods("GET")
	n.HandleFunc("/systemd/{unit}/property", routerFetchUnitProperty).Methods("GET")
	n.HandleFunc("/systemd/{unit}/propertyall", routerFetchUnitPropertyAll).Methods("GET")
	n.HandleFunc("/systemd/{unit}/property/{unittype}", routerFetchUnitTypeProperty).Methods("GET")

	// systemd configuration
	n.HandleFunc("/systemd/conf", routerConfigureSystemdConf)
	n.HandleFunc("/systemd/conf/update", routerConfigureSystemdConf)
}
