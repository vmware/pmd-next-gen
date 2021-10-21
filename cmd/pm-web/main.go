// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"runtime"

	"github.com/pm-web/pkg/conf"
	"github.com/pm-web/pkg/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := conf.Parse()
	if err != nil {
		log.Errorf("Failed to parse conf file %s: %s", conf.ConfFile, err)
	}

	log.Infof("pm-webd: v%s (built %s)", conf.Version, runtime.Version())

	if err := router.StartRouter(c); err != nil {
		log.Fatalf("Failed to start pm-webd: %v", err)
		os.Exit(1)
	}
}
