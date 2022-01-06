## photon-mgmtd


`photon-mgmtd` is a high performance open-source, simple, and pluggable REST API gateway designed with stateless architecture. It is written in Go, and built with performance in mind. It features real time health monitoring, configuration and performance for systems (containers), networking and applications.

- Proactive Monitoring and Analytics
  `photon-mgmtd` saves network administrators time and frustration because it makes it easy to gather statistics and perform analyses.
- Platform independent REST APIs can be accessed via any application (curl, chrome, PostMan ...) from any OS (Linux, IOS, Android, Windows ...)
- Minimal data transfer using JSON.
- Plugin based Architechture. See how to write plugin section for more information.

### Features!

|Feature| Details |
| ------ | ------ |
|systemd  | information, services (start, stop, restart, status), service properties for example CPUShares
see information from ```/proc``` fs| netstat, netdev, memory and much more


### Dependencies

photon-mgmtd uses a following open source projects to work properly:

* [logrus](https://github.com/sirupsen/logrus)
* [gorilla mux](https://github.com/gorilla/mux)
* [netlink](https://github.com/vishvananda/netlink)
* [gopsutil](https://github.com/shirou/gopsutil)
* [coreos go-systemd](https://github.com/coreos/go-systemd)
* [dbus](https://github.com/godbus/dbus)
* [ethtool](https://github.com/safchain/ethtool)
* [viper](https://github.com/spf13/viper)
* [go-ini](https://github.com/go-ini/ini)


### Installation

First configure your ```$GOPATH```. If you have already done this skip this step.

```sh
# keep in ~/.bashrc
```

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export OS_OUTPUT_GOPATH=1
```

Clone inside src dir of ```$GOPATH```. In my case

```sh
$ pwd
/home/sus/go/src
```

#### Building from source
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

Configuration file `photon-mgmt.toml` located in `/etc/photon-mgmt/` directory to manage the configuration.

The `[System]` section takes following Keys:

`LogLevel=`

Specifies the log level. Takes one of `Trace`, `Debug`, `Info`, `Warning`, `Error`, `Fatal` and `Panic`. Defaults to `info`. See [sirupsen](https://github.com/sirupsen/logrus#level-logging)

`UseAuthentication=`

A boolean. Specifies whether the users should be authenticated. Defaults to `true`.


The `[Network]` section takes following Keys:

`Listen=`

Specifies the IP address and port which the REST API server will listen to. When enabled, defaults to `127.0.0.1:5208`.

`ListenUnixSocket=`

A boolean. Specifies whether the server would listen on a unix domain socket `/run/photon-mgmt/photon-mgmt.sock`. Defaults to `true`.

Note that when both `ListenUnixSocket=` and `Listen=` are enabled, server listens on the unix domain socket by default.
 ```bash
❯ sudo cat /etc/photon-mgmt/photon-mgmt.toml
[System]
LogLevel="debug"
UseAuthentication="false"

[Network]
ListenUnixSocket="true"
```

```bash
❯ sudo systemctl status photon-mgmtd.service
● photon-mgmtd.service - A REST API based configuration management microservice gateway
     Loaded: loaded (/usr/lib/systemd/system/photon-mgmtd.service; disabled; vendor preset: disabled)
     Active: active (running) since Thu 2022-01-06 16:32:19 IST; 4s ago
   Main PID: 230041 (photon-mgmtd)
      Tasks: 6 (limit: 15473)
     Memory: 2.9M
        CPU: 7ms
     CGroup: /system.slice/photon-mgmtd.service
             └─230041 /usr/bin/photon-mgmtd

Jan 06 16:32:19 Zeus systemd[1]: photon-mgmtd.service: Passing 0 fds to service
Jan 06 16:32:19 Zeus systemd[1]: photon-mgmtd.service: About to execute /usr/bin/photon-mgmtd
Jan 06 16:32:19 Zeus systemd[1]: photon-mgmtd.service: Forked /usr/bin/photon-mgmtd as 230041
Jan 06 16:32:19 Zeus systemd[1]: photon-mgmtd.service: Changed failed -> running
Jan 06 16:32:19 Zeus systemd[1]: photon-mgmtd.service: Job 56328 photon-mgmtd.service/start finished, result=done
Jan 06 16:32:19 Zeus systemd[1]: Started photon-mgmtd.service - A REST API based configuration management microservice gateway.
Jan 06 16:32:19 Zeus systemd[230041]: photon-mgmtd.service: Executing: /usr/bin/photon-mgmtd
Jan 06 16:32:19 Zeus photon-mgmtd[230041]: time="2022-01-06T16:32:19+05:30" level=info msg="photon-mgmtd: v0.1 (built go1.18beta1)"
Jan 06 16:32:19 Zeus photon-mgmtd[230041]: time="2022-01-06T16:32:19+05:30" level=info msg="Starting photon-mgmtd... Listening on unix domain socket='/run/photon-mgmt/photon-mgmt>
```

#### pmctl
----

`pmctl` is a CLI tool allows to view and configure system/network/service status.

```bash
❯ pmctl service status nginx.service
                  Name: nginx.service
           Description: The nginx HTTP and reverse proxy server
               MainPid: 45732
             LoadState: loaded
           ActiveState: active
              SubState: running
         UnitFileState: disabled
  StateChangeTimeStamp: Sun Oct 31 12:02:02 IST 2021
  ActiveEnterTimestamp: Sun Oct 31 12:02:02 IST 2021
 InactiveExitTimestamp: Sun Oct 31 12:02:02 IST 2021
   ActiveExitTimestamp: 0
 InactiveExitTimestamp: Sun Oct 31 12:02:02 IST 2021
                Active: active (running) since Sun Oct 31 12:02:02 IST 2021

```


```bash
❯ pmctl status  system
              System Name: Zeus
                   Kernel: Linux (5.14.0-0.rc7.54.fc36.x86_64) #1 SMP Mon Aug 23 13:55:32 UTC 2021
                  Chassis: vm
           Hardware Model: VMware Virtual Platform
          Hardware Vendor: VMware, Inc.
             Product UUID: 979e4d56b63718b18534e112e64cb18
         Operating System: VMware Photon OS/Linux
Operating System Home URL: https://vmware.github.io/photon/
          Systemd Version: v247.10-3.ph4
             Architecture: x86-64
           Virtualization: vmware
            Network State: routable (carrier)
     Network Online State: online
                      DNS: 172.16.130.2
                  Address: 172.16.130.132/24 on link ens33
                           172.16.130.131/24 on link ens33
                           fe80::3279:c56d:55f9:aed7/64 on link ens33
                           172.16.130.138/24 on link ens37
                  Gateway: 172.16.130.2 on link ens37
                           172.16.130.2 on link ens33
                   Uptime: Running Since (2 days, 3 hours, 8 minutes) Booted (Wed Dec 22 15:57:24 IST 2021) Users (9) Proc (284)
                   Memory: Total (13564788736) Used (13564788736) Free (589791232) Available (9723891712)
```


```
❯ pmctl status network -i ens33
             Name: ens33
Alternative Names: enp2s1
            Index: 2
        Link File: /usr/lib/systemd/network/99-default.link
     Network File: /etc/systemd/network/10-ens33.network
             Type: ether
            State: routable (configured)
           Driver: e1000
           Vendor: Intel Corporation
            Model: 82545EM Gigabit Ethernet Controller (Copper) (PRO/1000 MT Single Port Adapter)
             Path: pci-0000:02:01.0
    Carrier State: carrier
     Online State: online
IPv4Address State: routable
IPv6Address State: degraded
       HW Address: 00:0c:29:5f:d1:39
              MTU: 1500
        OperState: up
            Flags: up|broadcast|multicast
        Addresses: 172.16.130.132/24 172.16.130.131/24 fe80::3279:c56d:55f9:aed7/64
          Gateway: 172.16.130.2
              DNS: 172.16.130.2
```

### How to configure users ?

##### Unix domain socket

Any users added to the group photon-mgmt, they are allowed to access the unix socket.
```sh
# usermod -a -G photon-mgmt exampleusername
```

##### Web users via pmctl

Export the token key to the enviroment as below
```
❯ export PHOTON_MGMT_AUTH_TOKEN=secret
```

### How to configure TLS ?

Generate private key (.key)

```sh
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.......................+++++
.+++++
e is 65537 (0x010001)

openssl genrsa -out server.key 2048
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```sh
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:

```

Place ```server.crt``` and ```server.key``` in the dir ```/etc/photon-mgmt/tls```

```bash
[root@Zeus tls]# ls
server.crt  server.key
[root@Zeus tls]# pwd
/etc/photon-mgmt/cert

```

Use case: https

```sh
$ curl --header "X-Session-Token: secret" --request GET https://localhost:5208/api/v1/network/ethtool/vmnet8/get-link-features -k --tlsv1.2

```

## How to write your own plugin ?

photon-mgmtd is designed with robust plugin based architecture in mind. You can always add and remove modules to it with minimal effort
You can implement and incorporate application features very quickly. Because plug-ins are separate modules with well-defined interfaces,
you can quickly isolate and solve problems. You can create custom versions of an application with minimal source code modifications.

* Choose namespace under `plugins` directory (systemd, system, proc) where you want to put your module.
* Write sub router see for example ```plugins/systemd/```
* Write your module ```module.go``` and  ```module_router.go```
* Write ```RegisterRouterModule```
* Register ```RegisterRouterModule``` with parent router for example for ```login``` registered with
  ```RegisterRouterSystem``` under ```system``` namespace as ```login.RegisterRouterLogin```
* See examples directory how to write on your own plugin.

### Todos

 - Write Tests
 - networkd

License
----

Apache 2.0
