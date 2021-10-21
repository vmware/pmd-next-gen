// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"github.com/pm-web/pkg/share"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	Version  = "0.1"
	ConfPath = "/etc/pm-web"
	ConfFile = "pmweb"
	TLSCert  = "cert/server.crt"
	TLSKey   = "cert/server.key"

	DefaultLogLevel   = "info"
	UseAuthentication = "true"

	DefaultIP        = "0.0.0.0"
	DefaultPort      = "8080"
	ListenUnixSocket = "true"
)

type Config struct {
	System  System  `mapstructure:"System"`
	Network Network `mapstructure:"Network"`
}

type System struct {
	LogLevel          string `mapstructure:"LogLevel"`
	UseAuthentication bool   `mapstructure:"UseAuthentication"`
}
type Network struct {
	IPAddress        string
	Port             string
	ListenUnixSocket bool
}

func Parse() (*Config, error) {
	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	viper.SetDefault("System.LogLevel", DefaultLogLevel)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to parse config file. Using defaults: %v", err)
	}

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		logrus.Errorf("Failed to decode config into struct, %v", err)
	}

	if l, err := logrus.ParseLevel(c.System.LogLevel); err != nil {
		logrus.Warn("Failed to parse log level='%s', falling back to 'info': %v", c.System.LogLevel, err)
		c.System.LogLevel = DefaultLogLevel
	} else {
		logrus.SetLevel(l)
	}

	logrus.Debugf("Log level set to '%+v'", logrus.GetLevel().String())

	if c.Network.IPAddress != "" {
		if _, err := share.ParseIP(c.Network.IPAddress); err != nil {
			logrus.Errorf("Failed to parse IPAddress=%s, %s", c.Network.IPAddress, c.Network.Port)
		}
	}

	if c.Network.Port != "" {
		if _, err := share.ParsePort(c.Network.Port); err != nil {
			logrus.Errorf("Failed to parse conf file Port=%s", c.Network.Port)
		}
	}

	if c.Network.IPAddress != "" && c.Network.Port != "" {
		logrus.Debugf("Parsed IPAddress=%s and Port=%s", c.Network.IPAddress, c.Network.Port)
	}

	return &c, nil
}
