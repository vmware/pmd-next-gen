// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pmd/pkg/share"
)

// App Version
const (
	Version  = "0.1"
	ConfPath = "/etc/pm-web"
	ConfFile = "pmweb"
	TLSCert  = "cert/server.crt"
	TLSKey   = "cert/server.key"

	DefaultLogLevel  = "info"
	DefaultLogFormat = "text"

	DefaultIP   = "0.0.0.0"
	DefaultPort = "8080"
)

// flag
var (
	IPFlag   string
	PortFlag string
)

//Config config file key value
type Config struct {
	System  System  `mapstructure:"System"`
	Network Network `mapstructure:"Network"`
}

type System struct {
	LogLevel  string `mapstructure:"LogLevel"`
	LogFormat string `mapstructure:"LogFormat"`
}
type Network struct {
	IPAddress string
	Port      string
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

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Failed to parse  config file, %v", err)
	}

	viper.SetDefault("System.LogFormat", DefaultLogLevel)
	viper.SetDefault("System.LogLevel", DefaultLogFormat)
	viper.SetDefault("Network.IPAddress", DefaultIP)
	viper.SetDefault("Network.Port", DefaultPort)

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

	_, err := share.ParseIP(c.Network.IPAddress)
	if err != nil {
		logrus.Errorf("Failed to parse IPAddress=%s, %s", c.Network.IPAddress, c.Network.Port)
		return nil, err
	}

	_, err = share.ParsePort(c.Network.Port)
	if err != nil {
		logrus.Errorf("Failed to parse conf file Port=%s", c.Network.Port)
		return nil, err
	}

	logrus.Debugf("Parsed IPAddress=%s and Port=%s", c.Network.IPAddress, c.Network.Port)

	return &c, nil
}
