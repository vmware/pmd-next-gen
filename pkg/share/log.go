// SPDX-License-Identifier: Apache-2.0

package share

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultLogDir  = "/var/log/api-router"
	defaultLogFile = "api-router.log"
)

// InitLog inits the logger
func InitLog() error {
	log := logrus.New()
	log.Level = logrus.InfoLevel

	viper.AutomaticEnv()
	lvl := viper.GetString("API_ROUTERD_LOG_LEVEL")
	if lvl != "" {
		l, err := logrus.ParseLevel(lvl)
		if err != nil {
			log.WithField("level", lvl).Warn("Invalid log level, fallback to 'info'")
		} else {
			log.SetLevel(l)
		}
	}

	switch viper.GetString("API_ROUTERD_LOG_FORMAT") {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
	case "text":
		log.SetFormatter(&logrus.TextFormatter{})
	}

	logDir := viper.GetString("API_ROUTERD_LOG_DIR")
	if logDir == "" {
		logDir = defaultLogDir
	}

	err := CreateDirectory(logDir, 0644)
	if err != nil {
		log.Errorf("Failed to create log directory. path: %s, err: %s", logDir, err)
		return err
	}

	logFile := path.Join(logDir, defaultLogFile)
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("Failed to create log file. path: %s, err: %s", logFile, err)
		return err
	}

	log.SetOutput(f)
	log.SetReportCaller(true)
	log.Info("Starting API Router")

	return nil
}
