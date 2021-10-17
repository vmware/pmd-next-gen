// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/conf"
	"github.com/pmd/pkg/router"
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
