// SPDX-License-Identifier: Apache-2.0

package proc

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/system"
	"github.com/pmd/pkg/web"
)

const (
	procMiscPath    = "/proc/misc"
	procNetArpPath  = "/proc/net/arp"
	procModulesPath = "/proc/modules"
)

type NetARP struct {
	IPAddress string `json:"ip_address"`
	HWType    string `json:"hw_type"`
	Flags     string `json:"flags"`
	HWAddress string `json:"hw_address"`
	Mask      string `json:"mask"`
	Device    string `json:"device"`
}

type Modules struct {
	Module     string `json:"module"`
	MemorySize string `json:"memory_size"`
	Instances  string `json:"instances"`
	Dependent  string `json:"dependent"`
	State      string `json:"state"`
}

func GetVersion(w http.ResponseWriter) error {
	infoStat, err := host.Info()
	if err != nil {
		return err
	}

	return web.JSONResponse(infoStat, w)
}

func GetPlatformInformation(w http.ResponseWriter) error {
	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return err
	}

	p := struct {
		Platform string
		Family   string
		Version  string
	}{
		platform,
		family,
		version,
	}

	return web.JSONResponse(p, w)
}

func GetVirtualization(w http.ResponseWriter) error {
	system, role, err := host.Virtualization()
	if err != nil {
		return err
	}

	v := struct {
		System string
		Role   string
	}{
		system,
		role,
	}

	return web.JSONResponse(v, w)
}

func GetUserStat(w http.ResponseWriter) error {
	userstat, err := host.Users()
	if err != nil {
		return err
	}

	return web.JSONResponse(userstat, w)
}

func GetTemperatureStat(w http.ResponseWriter) error {
	tempstat, err := host.SensorsTemperatures()
	if err != nil {
		return err
	}

	return web.JSONResponse(tempstat, w)
}

// GetNetStat read netstat from proc tcp/udp/sctp
func GetNetStat(w http.ResponseWriter, protocol string) error {
	conn, err := net.Connections(protocol)
	if err != nil {
		return err
	}

	return web.JSONResponse(conn, w)
}

func GetNetStatPid(w http.ResponseWriter, protocol string, process string) error {
	pid, err := strconv.ParseInt(process, 10, 32)
	if err != nil || protocol == "" || pid == 0 {
		return errors.New("can't parse request")
	}

	conn, err := net.ConnectionsPid(protocol, int32(pid))
	if err != nil {
		return err
	}

	return web.JSONResponse(conn, w)
}

func GetProtoCountersStat(w http.ResponseWriter) error {
	protocols := []string{"ip", "icmp", "icmpmsg", "tcp", "udp", "udplite"}

	proto, err := net.ProtoCounters(protocols)
	if err != nil {
		return err
	}

	return web.JSONResponse(proto, w)
}

func GetNetDev(w http.ResponseWriter) error {
	netdev, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	return web.JSONResponse(netdev, w)
}

func GetInterfaceStat(w http.ResponseWriter) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	return web.JSONResponse(interfaces, w)
}

func GetSwapMemoryStat(w http.ResponseWriter) error {
	swap, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	return web.JSONResponse(swap, w)
}

func GetVirtualMemoryStat(w http.ResponseWriter) error {
	virt, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	return web.JSONResponse(virt, w)
}

func GetCPUInfo(w http.ResponseWriter) error {
	cpus, err := cpu.Info()
	if err != nil {
		return err
	}

	return web.JSONResponse(cpus, w)
}

func GetCPUTimeStat(w http.ResponseWriter) error {
	cpus, err := cpu.Times(true)
	if err != nil {
		return err
	}

	return web.JSONResponse(cpus, w)
}

func GetAvgStat(w http.ResponseWriter) error {
	avgstat, err := load.Avg()
	if err != nil {
		return err
	}

	return web.JSONResponse(avgstat, w)
}

func GetPartitions(w http.ResponseWriter) error {
	p, err := disk.Partitions(true)
	if err != nil {
		return err
	}

	return web.JSONResponse(p, w)
}

// GetIOCounters read IO counters
func GetIOCounters(w http.ResponseWriter) error {
	i, err := disk.IOCounters()
	if err != nil {
		return err
	}

	return web.JSONResponse(i, w)
}

// GetDiskUsage read disk usage
func GetDiskUsage(w http.ResponseWriter) error {
	u, err := disk.Usage("/")
	if err != nil {
		return err
	}

	return web.JSONResponse(u, w)
}

// GetMisc read /proc/misc
func GetMisc(w http.ResponseWriter) error {
	lines, err := system.ReadFullFile(procMiscPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procMiscPath)
		return err
	}

	miscMap := make(map[int]string)
	for _, line := range lines {
		fields := strings.Fields(line)

		deviceNum, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		miscMap[deviceNum] = fields[1]
	}

	return web.JSONResponse(miscMap, w)
}

// GetNetArp get ARP info
func GetNetArp(w http.ResponseWriter) error {
	lines, err := system.ReadFullFile(procNetArpPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procNetArpPath)
		return err
	}

	netarp := make([]NetARP, len(lines)-1)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)

		arp := NetARP{}
		arp.IPAddress = fields[0]
		arp.HWType = fields[1]
		arp.Flags = fields[2]
		arp.HWAddress = fields[3]
		arp.Mask = fields[4]
		arp.Device = fields[5]
		netarp[i-1] = arp
	}

	return web.JSONResponse(netarp, w)
}

// GetModules Get all installed modules
func GetModules(w http.ResponseWriter) error {
	lines, err := system.ReadFullFile(procModulesPath)
	if err != nil {
		log.Fatalf("Failed to read: %s", procModulesPath)
		return err
	}

	modules := make([]Modules, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)

		module := Modules{}

		for j, field := range fields {
			switch j {
			case 0:
				module.Module = field

			case 1:
				module.MemorySize = field

			case 2:
				module.Instances = field

			case 3:
				module.Dependent = field

			case 4:
				module.State = field
			}
		}

		modules[i] = module
	}

	return web.JSONResponse(modules, w)
}

// GetProcessInfo get process information from proc
func GetProcessInfo(w http.ResponseWriter, proc string, property string) error {
	pid, err := strconv.ParseInt(proc, 10, 32)
	if err != nil {
		return err
	}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}

	switch property {
	case "pid-connections":
		conn, err := p.Connections()
		if err != nil {
			return err
		}

		return web.JSONResponse(conn, w)

	case "pid-rlimit":
		rlimit, err := p.Rlimit()
		if err != nil {
			return err
		}

		return web.JSONResponse(rlimit, w)

	case "pid-rlimit-usage":
		rlimit, err := p.RlimitUsage(true)
		if err != nil {
			return err
		}

		return web.JSONResponse(rlimit, w)

	case "pid-status":
		s, err := p.Status()
		if err != nil {
			return err
		}

		return web.JSONResponse(s, w)

	case "pid-username":
		u, err := p.Username()
		if err != nil {
			return err
		}

		return web.JSONResponse(u, w)

	case "pid-open-files":
		f, err := p.OpenFiles()
		if err != nil {
			return err
		}

		return web.JSONResponse(f, w)

	case "pid-fds":
		f, err := p.NumFDs()
		if err != nil {
			return err
		}

		return web.JSONResponse(f, w)

	case "pid-name":
		n, err := p.Name()
		if err != nil {
			return err
		}

		return web.JSONResponse(n, w)

	case "pid-memory-percent":
		m, err := p.MemoryPercent()
		if err != nil {
			return err
		}

		return web.JSONResponse(m, w)

	case "pid-memory-maps":
		m, err := p.MemoryMaps(true)
		if err != nil {
			return err
		}

		return web.JSONResponse(m, w)

	case "pid-memory-info":
		m, err := p.MemoryInfo()
		if err != nil {
			return err
		}

		return web.JSONResponse(m, w)

	case "pid-io-counters":
		m, err := p.IOCounters()
		if err != nil {
			return err
		}

		return web.JSONResponse(m, w)
	}

	return nil
}
