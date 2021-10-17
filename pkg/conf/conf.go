// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pmd/pkg/share"
)

const (
	Version  = "0.1"
	ConfPath = "/etc/pm-web"
	ConfFile = "pmweb"
	TLSCert  = "cert/server.crt"
	TLSKey   = "cert/server.key"

	DefaultLogLevel  = "info"
	DefaultLogFormat = "text"
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
	LogLevel  string `mapstructure:"LogLevel"`
	LogFormat string `mapstructure:"LogFormat"`
	UseAuthentication bool `mapstructure:"UseAuthentication"`
}
type Network struct {
	IPAddress        string
	Port             string
	ListenUnixSocket bool
}

func SetLogLevel(level string) error {
	if level == "" {
		return errors.New("unsupported")
	}

	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warn("Failed to parse log level, falling back to 'info'")
		return errors.New("unsupported")
	} else {
		logrus.SetLevel(l)
	}

	return nil
}

func SetLogFormat(format string) error {
	if format == "" {
		return errors.New("unsupported")
	}

	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
		})

	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		})

	default:
		logrus.Warn("Failed to parse log format, falling back to 'text'")
		return errors.New("unsupported")
	}

	return nil
}

func Parse() (*Config, error) {
	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	viper.SetDefault("System.LogFormat", DefaultLogLevel)
	viper.SetDefault("System.LogLevel", DefaultLogFormat)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to parse config file. Using defaults: %v", err)
	}

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		logrus.Errorf("Failed to decode config into struct, %v", err)
	}

	if err := SetLogLevel(viper.GetString("PM_WEBD_LOG_LEVEL")); err != nil {
		if err := SetLogLevel(c.System.LogLevel); err != nil {
			c.System.LogLevel = DefaultLogLevel
		}
	}

	logrus.Debugf("Log level set to '%+v'", logrus.GetLevel().String())

	if err := SetLogFormat(viper.GetString("PM_WEBD_LOG_FORMAT")); err != nil {
		if err = SetLogFormat(c.System.LogFormat); err != nil {
			c.System.LogLevel = DefaultLogFormat
		}
	}

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
