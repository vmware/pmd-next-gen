// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/distro-management-api/pkg/share"
)

const (
	Version  = "0.1"
	ConfPath = "/etc/distro-management-api"
	ConfFile = "distro-management-api"
	TLSCert  = "cert/server.crt"
	TLSKey   = "cert/server.key"

	DefaultLogLevel   = "info"
	UseAuthentication = "true"

	DefaultIP        = "127.0.0.1"
	DefaultPort      = "5208"
	ListenUnixSocket = "true"

	UnixDomainSocketPath = "/run/distro-management-api/distro-management-api.sock"
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
	Listen       string
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

	if c.Network.Listen != "" {
		share.ParseIpPort(c.Network.Listen)
		if _, _, err := share.ParseIpPort(c.Network.Listen); err != nil {
			logrus.Errorf("Failed to parse Listen=%s", c.Network.Listen)
			return nil, err
		}
	}


	return &c, nil
}
