// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/distro-management-api/pkg/conf"
	"github.com/distro-management-api/pkg/server"
	"github.com/distro-management-api/pkg/system"
)

func main() {
	c, err := conf.Parse()
	if err != nil {
		log.Errorf("Failed to parse conf file %s: %s", conf.ConfFile, err)
	}

	log.Infof("distro-management-apid: v%s (built %s)", conf.Version, runtime.Version())

	runtime.LockOSThread()

	cred, err := system.GetUserCredentials("")
	if err != nil {
		log.Warningf("Failed to get current user credentials: %+v", err)
		os.Exit(1)
	} else {
		if cred.Uid == 0 {
			u, err := system.GetUserCredentials("distro-management-api")
			if err != nil {
				log.Errorf("Failed to get user 'distro-management-api' credentials: %+v", err)
				os.Exit(1)
			} else {
				if err := system.CreateStateDirs("/run/distro-management-api", int(u.Uid), int(u.Gid)); err != nil {
					log.Errorf("Failed to create runtime dir '/run/distro-management-api': %+v", err)
					os.Exit(1)
				}

				if err := system.EnableKeepCapability(); err != nil {
					log.Warningf("Failed to enable keep capabilities: %+v", err)
				}

				if err := system.SwitchUser(u); err != nil {
					log.Warningf("Failed to switch user: %+v", err)
				}

				if err := system.DisableKeepCapability(); err != nil {
					log.Warningf("Failed to disable keep capabilities: %+v", err)
				}

				err := system.ApplyCapability(u)
				if err != nil {
					log.Warningf("Failed to apply capabilities: +%v", err)
				}
			}
		}
	}

	if err := server.Run(c); err != nil {
		log.Fatalf("Failed to start distro-management-apid: %v", err)
		os.Exit(1)
	}
}
