## photon-mgmtd


`photon-mgmtd` is a high performance open-source, simple, and pluggable REST API gateway designed with stateless architecture. It is written in Go, and built with performance in mind. It features real time health monitoring, configuration and performance for systems (containers), networking and applications.

- Proactive Monitoring and Analytics
  easy to gather statistics and perform analyses.
- Platform independent REST APIs can be accessed via any application (curl, chrome, PostMan ...) from any OS (Linux, IOS, Android, Windows ...)
- Minimal data transfer using JSON.
- Plugin based architechture. See how to write plugin section for more information.

### Features!

- systemd   information, services (start, stop, restart, status), service properties for example CPUShares
- see information from ```/proc``` fs| netstat, netdev, memory , vms, ARP and much more
- system fetch and configure system information for example hostname
- network fetch and configure network information example (dns, iostat, interface)
- network link  configure network link parameters like (dhcp, linkLocalAddressing, multicastDNS, Address, route, domains, dns, ntp, ipv6AcceptRA, mode, - - mtubytes, mac, group, requiredFamilyForOnline, activationPolicy, routingPolicyRule, DHCPv4, DHCPv6, DHCPServer, Ipv6SendRA) etc
- login  fetch list of users and sessions also get information for a id
- network devices  create and remove virtual network devices like (Vlan, Bond, Bridge, MacVLan, IpVLan, VxLan, WireGuard) etc
- ethtool  fetch ethernet settings for a link also based on a action
- sysctl  used to fetch, set, load and automate kernel parameters
- user used to fetch, add, and remove user on the system
- group  used to fetch, add, and remove group on the system
- link  configure link parameters like (MACAddress, Name, AlternativeNames, Offload, VLANTAG, CHannels, Buffers, Queues, FlowControls, Coalesce) etc
- firewall  add, delete and show nft tables, chain and rules also is used to run any NFT commands
- package management (tdnf)  used to manage package management on the system like (list, info, download, update, remove, clean cache, list repositories,   search package) etc

#### Building and installation from source
----

```bash

❯ make build
❯ sudo make install
```

Due to security `photon-mgmtd` runs in non root user `photon-mgmt`. It drops all privileges except `CAP_NET_ADMIN` and `CAP_SYS_ADMIN`.

```bash

❯  useradd -M -s /usr/bin/nologin photon-mgmt
```

#### Configuration
----

Configuration file `mgmt.toml` located in `/etc/photon-mgmt/` directory to manage the configuration.

The `[System]` section takes following Keys:

`LogLevel=`

Specifies the log level. Takes one of `Trace`, `Debug`, `Info`, `Warning`, `Error`, `Fatal` and `Panic`. Defaults to `info`. See [sirupsen](https://github.com/sirupsen/logrus#level-logging)

`UseAuthentication=`
A boolean. Specifies whether the users should be authenticated. Defaults to `true`.

The `[Network]` section takes following Keys:

`Listen=`
Specifies the IP address and port which the REST API server will listen to. When enabled, defaults to `127.0.0.1:5208`.

`ListenUnixSocket=`
A boolean. Specifies whether the server would listen on a unix domain socket `/run/photon-mgmt/mgmt.sock`. Defaults to `true`.

Note that when both `ListenUnixSocket=` and `Listen=` are enabled, server listens on the unix domain socket by default.
 ```bash
❯ sudo cat /etc/photon-mgmt/mgmt.toml
[System]
LogLevel="debug"
UseAuthentication="false"

[Network]
ListenUnixSocket="true"
```

```bash
❯ sudo systemctl start photon-mgmtd
```

```bash
❯ sudo systemctl status photon-mgmtd
● photon-mgmtd.service - A REST API based configuration management microservice gateway
     Loaded: loaded (8;;file://zeus/usr/lib/systemd/system/photon-mgmtd.service^G/usr/lib/systemd/system/photon-mgmtd.service8;;^G; enabled; preset: enabled)
     Active: active (running) since Thu 2023-01-26 11:34:05 UTC; 2min 44s ago
   Main PID: 668 (photon-mgmtd)
      Tasks: 6 (limit: 18735)
     Memory: 22.8M
     CGroup: /system.slice/photon-mgmtd.service
             └─668 /usr/bin/photon-mgmtd

Jan 26 11:34:05 zeus systemd[1]: photon-mgmtd.service: Changed dead -> running
Jan 26 11:34:05 zeus systemd[1]: photon-mgmtd.service: Job 185 photon-mgmtd.service/start finished, result=done
Jan 26 11:34:05 zeus systemd[1]: Started A REST API based configuration management microservice gateway.
Jan 26 11:34:05 zeus systemd[668]: photon-mgmtd.service: Executing: /usr/bin/photon-mgmtd
Jan 26 11:34:05 zeus photon-mgmtd[668]: time="2023-01-26T11:34:05Z" level=info msg="photon-mgmtd: v0.1 (built go1.19.3)"
Jan 26 11:34:05 zeus photon-mgmtd[668]: time="2023-01-26T11:34:05Z" level=info msg="Starting photon-mgmtd... Listening on unix domain socket='/run/photon-mgmt/mgmt.sock' in HTTP>
Jan 26 11:36:43 zeus systemd[1]: photon-mgmtd.service: Trying to enqueue job photon-mgmtd.service/start/replace
Jan 26 11:36:43 zeus systemd[1]: photon-mgmtd.service: Installed new job photon-mgmtd.service/start as 596
Jan 26 11:36:43 zeus systemd[1]: photon-mgmtd.service: Enqueued job photon-mgmtd.service/start as 596
Jan 26 11:36:43 zeus systemd[1]: photon-mgmtd.service: Job 596 photon-mgmtd.service/start finished, result=done
```

For a comprehensive list use cases, see [usecases](https://github.com/vmware/pmd-next-gen/blob/main/USECASES.md).
