## distro-management-apid


distro-management-apid is a high performance open-source, simple, and pluggable REST API gateway designed with stateless architecture.It is written in Go, and built with performance in mind. It features real time health monitoring, configuration and performance for systems (containers), networking and applications.

- Proactive Monitoring and Analytics
  distro-management-apid saves network administrators time and frustration because it makes it easy to gather statistics and perform analyses.
- Platform independent REST APIs can be accessed via any application (curl, chrome, PostMan ...) from any OS (Linux, IOS, Android, Windows ...)
- Minimal data transfer using JSON.
- Plugin based Architechture. See how to write plugin section for more information.

### Features!

|Feature| Details |
| ------ | ------ |
|systemd  | information, services (start, stop, restart, status), service properties for example CPUShares
see information from ```/proc``` fs| netstat, netdev, memory and much more


### Dependencies

distro-management-apid uses a following open source projects to work properly:

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

Due to security `distro-management-apid` runs in non root user `distro-management-api`. It drops all privileges except `CAP_NET_ADMIN` and `CAP_SYS_ADMIN`.

```bash

❯  useradd -M -s /usr/bin/nologin distro-management-api
```

#### Configuration
----

Configuration file `distro-management-api.toml` located in `/etc/distro-management-api/` directory to manage the configuration.

The `[System]` section takes following Keys:

`LogLevel=`

Specifies the log level. Takes one of `Trace`, `Debug`, `Info`, `Warning`, `Error`, `Fatal` and `Panic`. Defaults to `info`. See [sirupsen](https://github.com/sirupsen/logrus#level-logging)

`UseAuthentication=`

A boolean. Specifies whether the users should be authenticated. Defaults to `true`.


The `[Network]` section takes following Keys:

`Listen=`

Specifies the IP address and port which the REST API server will listen to. When enabled, defaults to `127.0.0.1:5208`.

`ListenUnixSocket=`

A boolean. Specifies whether the server would listen on a unix domain socket `/run/distro-management-api/distro-management-api.sock`. Defaults to `true`.

Note that when both `ListenUnixSocket=` and `Listen=` are enabled, server listens on the unix domain socket by default.
 ```bash
❯ sudo cat /etc/distro-management-api/pmweb.toml                                     
[System]
LogLevel="debug"
UseAuthentication="false"

[Network]
ListenUnixSocket="true"
```

```bash
❯ sudo systemctl status distro-management-apid.service
● distro-management-apid.service - A REST API Microservice Gateway
     Loaded: loaded (/usr/lib/systemd/system/distro-management-apid.service; disabled; vendor preset: disabled)
     Active: active (running) since Mon 2021-12-13 14:07:46 IST; 22s ago
       Docs: man:distro-management-apid.conf(5)
   Main PID: 27981 (distro-management-apid)
      Tasks: 6 (limit: 15473)
     Memory: 2.1M
        CPU: 8ms
     CGroup: /system.slice/distro-management-apid.service
             └─27981 /usr/bin/distro-management-apid

Dec 13 14:07:46 Zeus systemd[1]: distro-management-apid.service: Forked /usr/bin/distro-management-apid as 27981
Dec 13 14:07:46 Zeus distro-management-apid[27981]: time="2021-12-13T14:07:46+05:30" level=info msg="distro-management-apid: v0.1 (built go1.17.4)"
Dec 13 14:07:46 Zeus systemd[1]: distro-management-apid.service: Changed dead -> running
Dec 13 14:07:46 Zeus systemd[1]: distro-management-apid.service: Job 23713 distro-management-apid.service/start finished, result=done
Dec 13 14:07:46 Zeus systemd[1]: Started A REST API Microservice Gateway.
Dec 13 14:07:46 Zeus systemd[27981]: distro-management-apid.service: Executing: /usr/bin/distro-management-apid
Dec 13 14:07:46 Zeus distro-management-apid[27981]: time="2021-12-13T14:07:46+05:30" level=info msg="Starting distro-management-apid server at unix domain socket='/run/pmwebd/pmwebd.sock' in HTTP mode"
                                  
```

#### pmctl
----

`pmctl` is a CLI tool allows to view system/service status.

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

### How to configure users ?

##### Unix domain socket

Any users added to the group distro-management-api, they are allowed to access the unix socket.
```sh
# usermod -a -G distro-management-api exampleusername
```

##### Web users via pmctl

Export the token key to the enviroment as below
```
❯ export PM_WEB_AUTH_TOKEN=secret
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

Place ```server.crt``` and ```server.key``` in the dir ```/etc/distro-management-api/tls```

```bash
[root@Zeus tls]# ls
server.crt  server.key
[root@Zeus tls]# pwd
/etc/distro-management-api/cert

```

Use case: https

```sh
$ curl --header "X-Session-Token: secret" --request GET https://localhost:5208/api/v1/network/ethtool/vmnet8/get-link-features -k --tlsv1.2

```

## How to write your own plugin ?

distro-management-apid is designed with robust plugin based architecture in mind. You can always add and remove modules to it with minimal effort
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
