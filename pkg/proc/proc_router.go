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

func routerGetProcNetDev(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetNetDev(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcVersion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetVersion(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcPlatformInformation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetPlatformInformation(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcVirtualization(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetVirtualization(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcUserStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetUserStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcTemperatureStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetTemperatureStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]

	switch r.Method {
	case "GET":
		if err := GetNetStat(w, protocol); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcPidNetStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	protocol := vars["protocol"]
	pid := vars["pid"]

	switch r.Method {
	case "GET":
		if err := GetNetStatPid(w, protocol, pid); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcInterfaceStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetInterfaceStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcProtoCountersStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetProtoCountersStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcGetSwapMemoryStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetSwapMemoryStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcVirtualMemoryStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetVirtualMemoryStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcCPUInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetCPUInfo(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcCPUTimeStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetCPUTimeStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcAvgStat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetAvgStat(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func configureProcSysVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vm := VM{
		Property: vars["path"],
	}

	switch r.Method {
	case "GET":
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
	case "GET":
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

func routerGetProcMisc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetMisc(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcNetArp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetNetArp(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcModules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetModules(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetProcProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid := vars["pid"]
	property := vars["property"]

	switch r.Method {
	case "GET":

		if err := GetProcessInfo(w, pid, property); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetPartitions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := GetPartitions(w)
		if err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetIOCounters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetIOCounters(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func routerGetDiskUsage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := GetDiskUsage(w); err != nil {
			web.JSONResponseError(err, w)
		}
	}
}

func RegisterRouterProc(router *mux.Router) {
	n := router.PathPrefix("/proc").Subrouter().StrictSlash(false)

	n.HandleFunc("/avgstat", routerGetProcAvgStat)
	n.HandleFunc("/cpuinfo", routerGetProcCPUInfo)
	n.HandleFunc("/cputimestat", routerGetProcCPUTimeStat)
	n.HandleFunc("/diskusage", routerGetDiskUsage)
	n.HandleFunc("/interface-stat", routerGetProcInterfaceStat)
	n.HandleFunc("/iocounters", routerGetIOCounters)
	n.HandleFunc("/misc", routerGetProcMisc)
	n.HandleFunc("/modules", routerGetProcModules)
	n.HandleFunc("/net/arp", routerGetProcNetArp)
	n.HandleFunc("/netdev", routerGetProcNetDev)
	n.HandleFunc("/netstat/{protocol}", routerGetProcNetStat)
	n.HandleFunc("/partitions", routerGetPartitions)
	n.HandleFunc("/platform", routerGetProcPlatformInformation)
	n.HandleFunc("/process/{pid}/{property}/", routerGetProcProcess)
	n.HandleFunc("/proto-counter-stat", routerGetProcProtoCountersStat)
	n.HandleFunc("/proto-pid-stat/{pid}/{protocol}", routerGetProcPidNetStat)
	n.HandleFunc("/swap-memory", routerGetProcGetSwapMemoryStat)
	n.HandleFunc("/sys/net/{path}/{link}/{conf}", configureProcSysNet)
	n.HandleFunc("/sys/vm/{path}", configureProcSysVM)
	n.HandleFunc("/temperaturestat", routerGetProcTemperatureStat)
	n.HandleFunc("/userstat", routerGetProcUserStat)
	n.HandleFunc("/version", routerGetProcVersion)
	n.HandleFunc("/virtual-memory", routerGetProcVirtualMemoryStat)
	n.HandleFunc("/virtualization", routerGetProcVirtualization)
}
