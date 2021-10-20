// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pmd/pkg/web"
)

type Info struct {
	Path     string `json:"path"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func routerFetchProcNetDev(w http.ResponseWriter, r *http.Request) {
	if err := FetchNetDev(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcVersion(w http.ResponseWriter, r *http.Request) {
	if err := FetchVersion(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcPlatformInformation(w http.ResponseWriter, r *http.Request) {
	if err := FetchPlatformInformation(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcVirtualization(w http.ResponseWriter, r *http.Request) {
	if err := FetchVirtualization(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcUserStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchUserStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcTemperatureStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchTemperatureStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
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

func routerFetchProcInterfaceStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchInterfaceStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcProtoCountersStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchProtoCountersStat(w); err != nil {
		web.JSONResponseError(err, w)
	}

}

func routerFetchProcFetchSwapMemoryStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchSwapMemoryStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcVirtualMemoryStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchVirtualMemoryStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcCPUInfo(w http.ResponseWriter, r *http.Request) {
	if err := FetchCPUInfo(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcCPUTimeStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchCPUTimeStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcAvgStat(w http.ResponseWriter, r *http.Request) {
	if err := FetchAvgStat(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func configureProcSysVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vm := VM{
		Property: vars["path"],
	}

	switch r.Method {
	case "Fetch":
		if err := vm.GetVM(w); err != nil {
			web.JSONResponseError(err, w)
		}
	case "PUT":

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
}

func configureProcSysNet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := SysNet{Path: vars["path"], Property: vars["conf"], Link: vars["link"]}

	switch r.Method {
	case "Get":
		if err := proc.GetSysNet(w); err != nil {
			web.JSONResponseError(err, w)
		}

	case "PUT":
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
}

func routerFetchProcMisc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "Fetch":
		if err := FetchMisc(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerFetchProcNetArp(w http.ResponseWriter, r *http.Request) {
	if err := FetchNetArp(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchProcModules(w http.ResponseWriter, r *http.Request) {
	if err := FetchModules(w); err != nil {
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

func routerFetchPartitions(w http.ResponseWriter, r *http.Request) {
	err := FetchPartitions(w)
	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchIOCounters(w http.ResponseWriter, r *http.Request) {
	if err := FetchIOCounters(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func routerFetchDiskUsage(w http.ResponseWriter, r *http.Request) {
	if err := FetchDiskUsage(w); err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/sys/net/{path}/{link}/{conf}", configureProcSysNet)
	n.HandleFunc("/sys/vm/{path}", configureProcSysVM)

	n.HandleFunc("/avgstat", routerFetchProcAvgStat).Methods("GET")
	n.HandleFunc("/cpuinfo", routerFetchProcCPUInfo).Methods("GET")
	n.HandleFunc("/cputimestat", routerFetchProcCPUTimeStat).Methods("GET")
	n.HandleFunc("/diskusage", routerFetchDiskUsage).Methods("GET")
	n.HandleFunc("/interface-stat", routerFetchProcInterfaceStat).Methods("GET")
	n.HandleFunc("/iocounters", routerFetchIOCounters).Methods("GET")
	n.HandleFunc("/misc", routerFetchProcMisc).Methods("GET")
	n.HandleFunc("/temperaturestat", routerFetchProcTemperatureStat).Methods("GET")
	n.HandleFunc("/userstat", routerFetchProcUserStat).Methods("GET")
	n.HandleFunc("/version", routerFetchProcVersion).Methods("GET")
	n.HandleFunc("/virtual-memory", routerFetchProcVirtualMemoryStat).Methods("GET")
	n.HandleFunc("/virtualization", routerFetchProcVirtualization).Methods("GET")
	n.HandleFunc("/modules", routerFetchProcModules).Methods("GET")
	n.HandleFunc("/net/arp", routerFetchProcNetArp).Methods("GET")
	n.HandleFunc("/netdev", routerFetchProcNetDev).Methods("GET")
	n.HandleFunc("/netstat/{protocol}", routerFetchProcNetStat).Methods("GET")
	n.HandleFunc("/partitions", routerFetchPartitions).Methods("GET")
	n.HandleFunc("/platform", routerFetchProcPlatformInformation).Methods("GET")
	n.HandleFunc("/process/{pid}/{property}/", routerFetchProcProcess).Methods("GET")
	n.HandleFunc("/proto-counter-stat", routerFetchProcProtoCountersStat).Methods("GET")
	n.HandleFunc("/proto-pid-stat/{pid}/{protocol}", routerFetchProcPidNetStat).Methods("GET")
	n.HandleFunc("/swap-memory", routerFetchProcFetchSwapMemoryStat).Methods("GET")
}
