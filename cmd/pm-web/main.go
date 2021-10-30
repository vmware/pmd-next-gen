// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/conf"
	"github.com/pm-web/pkg/router"
	"github.com/pm-web/pkg/system"
)

func main() {
	c, err := conf.Parse()
	if err != nil {
		log.Errorf("Failed to parse conf file %s: %s", conf.ConfFile, err)
	}

	log.Infof("pm-webd: v%s (built %s)", conf.Version, runtime.Version())

	cred, err := system.GetUserCredentials("")
	if err != nil {
		log.Warningf("Failed to get current user credentials: %+v", err)
	} else {
		if cred.Uid == 0 {
			u, err := system.GetUserCredentials("pm-web")
			if err != nil {
				log.Warningf("Failed to get user 'pm-web' credentials: %+v", err)
			} else {
				if err := system.CreateStateDirs("/run/pmwebd", int(u.Uid), int(u.Gid)); err != nil {
					log.Println(err)
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

	if err := router.StartHttpServer(c); err != nil {
		log.Fatalf("Failed to start pm-webd: %v", err)
		os.Exit(1)
	}
}
