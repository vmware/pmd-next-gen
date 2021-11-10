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

func routerFetchProcNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	if err := FetchNetStat(w, protocol); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcPidNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]
	pid := vars["pid"]

	if err := FetchNetStatPid(w, protocol, pid); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcSysVM(w http.ResponseWriter, r *http.Request) {
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

func routerFetchProcSysNet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{Path: vars["path"], Property: vars["conf"], Link: vars["link"]}

	if err := proc.GetSysNet(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func configureProcSysNet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{Path: vars["path"], Property: vars["conf"], Link: vars["link"]}

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

func routerFetchProcNetArp(w http.ResponseWriter, r *http.Request) {
	if err := FetchNetArp(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	if err := FetchProcessInfo(w, pid, property); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchSystem(w http.ResponseWriter, r *http.Request) {
	var err error
	v := mux.Vars(r)

	switch v["system"] {
	case "avgstat":
		err = FetchAvgStat(w)
	case "cpuinfo":
		err = FetchCPUInfo(w)
	case "cputimestat":
		err = FetchCPUTimeStat(w)
	case "diskusage":
		err = FetchDiskUsage(w)
	case "iocounters":
		err = FetchIOCounters(w)
	case "partitions":
		err = FetchPartitions(w)
	case "temperaturestat":
		err = FetchTemperatureStat(w)
	case "modules":
		err = FetchModules(w)
	case "misc":
		err = FetchMisc(w)
	case "userstat":
		err = FetchUserStat(w)
	case "version":
		err = FetchVersion(w)
	case "virtualmemory":
		err = FetchVirtualMemoryStat(w)
	case "virtualization":
		err = FetchVirtualization(w)
	case "platform":
		err = FetchPlatformInformation(w)
	case "swapmemory":
		err = FetchSwapMemoryStat(w)
	case "interfaces":
		err = FetchInterfaces(w)
	case "netdeviocounters":
		err = FetchNetDevIOCounters(w)
	case "protocounterstat":
		err = FetchProtoCountersStat(w)
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/sys/net/{path}/{link}/{conf}", routerFetchProcSysNet).Methods("GET")
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", configureProcSysNet).Methods("PUT")

	n.HandleFunc("/sys/vm/{path}", routerFetchProcSysVM).Methods("GET")
	n.HandleFunc("/sys/vm/{path}", routerConfigureProcSysVM).Methods("PUT")

	n.HandleFunc("/{system}", routerFetchSystem).Methods("GET")

	n.HandleFunc("/net/arp", routerFetchProcNetArp).Methods("GET")
	n.HandleFunc("/netstat/{protocol}", routerFetchProcNetStat).Methods("GET")

	n.HandleFunc("/process/{pid}/{property}/", routerFetchProcProcess).Methods("GET")
	n.HandleFunc("/protopidstat/{pid}/{protocol}", routerFetchProcPidNetStat).Methods("GET")

}
