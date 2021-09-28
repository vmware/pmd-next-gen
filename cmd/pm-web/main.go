// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/conf"
	"github.com/pmd/pkg/router"
	"github.com/pmd/pkg/share"
)

func main() {
	share.InitLog()

	err := conf.InitConf()
	if err != nil {
		log.Errorf("Failed to init conf file %s: %s", conf.ConfFile, err)
	}

	log.Infof("pm-webd: v%s (built %s)", conf.Version, runtime.Version())
	log.Infof("Start Server at %s:%s", conf.IPFlag, conf.PortFlag)

	err = router.StartRouter(conf.IPFlag, conf.PortFlag, path.Join(conf.ConfPath, conf.TLSCert), path.Join(conf.ConfPath, conf.TLSKey))
	if err != nil {
		log.Fatalf("Failed to init pm-webd: %v", err)
		os.Exit(1)
	}
}
