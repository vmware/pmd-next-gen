// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pm-web/pkg/web"
)

type Info struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func routerAcquireProcNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	if err := AcquireNetStat(w, protocol); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcPidNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]
	pid := vars["pid"]

	if err := AcquireNetStatPid(w, protocol, pid); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcSysVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vm := VM{
		Property: vars["path"],
	}

	if err := vm.GetVM(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerConfigureProcSysVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vm := VM{
		Property: vars["path"],
	}

	v := new(Info)

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
	vars := mux.Vars(r)
	proc := SysNet{
		Path: vars["path"],
		Property: vars["conf"],
		Link: vars["link"],
	}

	if err := proc.GetSysNet(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func configureProcSysNet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{
		Path: vars["path"],
		Property: vars["conf"],
		Link: vars["link"],
	}

	v := new(Info)

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
	if err := AcquireNetArp(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireProcProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	if err := AcquireProcessInfo(w, pid, property); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerAcquireSystem(w http.ResponseWriter, r *http.Request) {
	var err error
	v := mux.Vars(r)

	switch v["system"] {
	case "avgstat":
		err = AcquireAvgStat(w)
	case "cpuinfo":
		err = AcquireCPUInfo(w)
	case "cputimestat":
		err = AcquireCPUTimeStat(w)
	case "diskusage":
		err = AcquireDiskUsage(w)
	case "iocounters":
		err = AcquireIOCounters(w)
	case "partitions":
		err = AcquirePartitions(w)
	case "temperaturestat":
		err = AcquireTemperatureStat(w)
	case "modules":
		err = AcquireModules(w)
	case "misc":
		err = AcquireMisc(w)
	case "userstat":
		err = AcquireUserStat(w)
	case "version":
		err = AcquireVersion(w)
	case "virtualmemory":
		err = AcquireVirtualMemoryStat(w)
	case "virtualization":
		err = AcquireVirtualization(w)
	case "platform":
		err = AcquirePlatformInformation(w)
	case "swapmemory":
		err = AcquireSwapMemoryStat(w)
	case "interfaces":
		err = AcquireInterfaces(w)
	case "netdeviocounters":
		err = AcquireNetDevIOCounters(w)
	case "protocounterstat":
		err = AcquireProtoCountersStat(w)
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/sys/net/{path}/{link}/{conf}", routerAcquireProcSysNet).Methods("GET")
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", configureProcSysNet).Methods("PUT")

	n.HandleFunc("/sys/vm/{path}", routerAcquireProcSysVM).Methods("GET")
	n.HandleFunc("/sys/vm/{path}", routerConfigureProcSysVM).Methods("PUT")

	n.HandleFunc("/{system}", routerAcquireSystem).Methods("GET")

	n.HandleFunc("/net/arp", routerAcquireProcNetArp).Methods("GET")
	n.HandleFunc("/netstat/{protocol}", routerAcquireProcNetStat).Methods("GET")

	n.HandleFunc("/process/{pid}/{property}/", routerAcquireProcProcess).Methods("GET")
	n.HandleFunc("/protopidstat/{pid}/{protocol}", routerAcquireProcPidNetStat).Methods("GET")

}
