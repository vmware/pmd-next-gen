// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 VMware, Inc.

package proc

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vmware/pmd-next-gen/pkg/web"
)

type Proc struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func routerAcquireProcNetStat(w http.ResponseWriter, r *http.Request) {
	if err := AcquireNetStat(r.Context(), w, mux.Vars(r)["protocol"]); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcPidNetStat(w http.ResponseWriter, r *http.Request) {
	if err := AcquireNetStatPid(r.Context(), w, mux.Vars(r)["protocol"], mux.Vars(r)["pid"]); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcSysVM(w http.ResponseWriter, r *http.Request) {
	vm := VM{
		Property: mux.Vars(r)["property"],
	}

	if err := vm.GetVM(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerConfigureProcSysVM(w http.ResponseWriter, r *http.Request) {
	vm := VM{
		Property: mux.Vars(r)["property"],
	}

	v := Proc{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		web.JSONResponseError(err, w)
		return
	}

	vm.Value = v.Value
	if err := vm.SetVM(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcSysNet(w http.ResponseWriter, r *http.Request) {
	proc := SysNet{
		Path:     mux.Vars(r)["path"],
		Property: mux.Vars(r)["property"],
		Link:     mux.Vars(r)["link"],
	}

	if err := proc.GetSysNet(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func configureProcSysNet(w http.ResponseWriter, r *http.Request) {
	proc := SysNet{
		Path:     mux.Vars(r)["path"],
		Property: mux.Vars(r)["property"],
		Link:     mux.Vars(r)["link"],
	}

	v := Proc{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		web.JSONResponseError(err, w)
		return
	}

	proc.Value = v.Value
	if err := proc.SetSysNet(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcNetArp(w http.ResponseWriter, r *http.Request) {
	if err := AcquireNetArp(r.Context(), w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcProcess(w http.ResponseWriter, r *http.Request) {
	if err := AcquireProcessInfo(r.Context(), w, mux.Vars(r)["pid"], mux.Vars(r)["property"]); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireSystem(w http.ResponseWriter, r *http.Request) {
	var err error

	switch mux.Vars(r)["system"] {
	case "avgstat":
		err = AcquireAvgStat(r.Context(), w)
	case "cpuinfo":
		err = AcquireCPUInfo(r.Context(), w)
	case "cputimestat":
		err = AcquireCPUTimeStat(r.Context(), w)
	case "diskusage":
		err = AcquireDiskUsage(r.Context(), w)
	case "diskpartitions":
		err = AcquireDiskPartitions(r.Context(), w)
	case "iocounters":
		err = AcquireIOCounters(r.Context(), w)
	case "temperaturestat":
		err = AcquireTemperatureStat(r.Context(), w)
	case "modules":
		err = AcquireModules(r.Context(), w)
	case "misc":
		err = AcquireMisc(r.Context(), w)
	case "userstat":
		err = AcquireUserStat(r.Context(), w)
	case "hostinfo":
		err = AcquireHostInfo(r.Context(), w)
	case "virtualmemory":
		err = AcquireVirtualMemoryStat(r.Context(), w)
	case "virtualization":
		err = AcquireVirtualization(r.Context(), w)
	case "platform":
		err = AcquirePlatformInformation(r.Context(), w)
	case "interfaces":
		err = AcquireInterfaces(r.Context(), w)
	case "netdeviocounters":
		err = AcquireNetDevIOCounters(r.Context(), w)
	case "protocounterstat":
		err = AcquireProtoCountersStat(r.Context(), w)
	default:
		err = errors.New("not found")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/sys/net/{path}/{property}", routerAcquireProcSysNet).Methods("GET")
	n.HandleFunc("/sys/net/{path}/{link}/{property}", routerAcquireProcSysNet).Methods("GET")
	n.HandleFunc("/sys/net/{path}/{property}", configureProcSysNet).Methods("PUT")
	n.HandleFunc("/sys/net/{path}/{link}/{property}", configureProcSysNet).Methods("PUT")

	n.HandleFunc("/sys/vm/{property}", routerAcquireProcSysVM).Methods("GET")
	n.HandleFunc("/sys/vm/{property}", routerConfigureProcSysVM).Methods("PUT")

	n.HandleFunc("/{system}", routerAcquireSystem).Methods("GET")

	n.HandleFunc("/net/arp", routerAcquireProcNetArp).Methods("GET")
	n.HandleFunc("/netstat/{protocol}", routerAcquireProcNetStat).Methods("GET")

	n.HandleFunc("/process/{pid}/{property}", routerAcquireProcProcess).Methods("GET")
	n.HandleFunc("/protopidstat/{pid}/{protocol}", routerAcquireProcPidNetStat).Methods("GET")
}
