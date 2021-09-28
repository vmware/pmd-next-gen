// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pmd/pkg/share"
)

// App Version
const (
	Version  = "0.1"
	ConfPath = "/etc/pm-web"
	ConfFile = "pmweb"
	TLSCert  = "tls/server.crt"
	TLSKey   = "tls/server.key"
)

// flag
var (
	IPFlag   string
	PortFlag string
)

//Config config file key value
type Config struct {
	Server Network `mapstructure:"Network"`
}

//Network IP Address and Port
type Network struct {
	IPAddress string
	Port      string
}

func init() {
	const (
		defaultIP   = "0.0.0.0"
		defaultPort = "8080"
	)

	flag.StringVar(&IPFlag, "ip", defaultIP, "The server IP address.")
	flag.StringVar(&PortFlag, "port", defaultPort, "The server port.")
}

func parseConfFile() (Config, error) {
	var conf Config

	viper.SetConfigName(ConfFile)
	viper.AddConfigPath(ConfPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Faild to parse  config file, %v", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Errorf("Failed to decode config into struct, %v", err)
	}

	_, err = share.ParseIP(conf.Server.IPAddress)
	if err != nil {
		log.Errorf("Failed to parse IPAddress=%s, %s", conf.Server.IPAddress, conf.Server.Port)
		return conf, err
	}

	_, err = share.ParsePort(conf.Server.Port)
	if err != nil {
		log.Errorf("Failed to parse conf file Port=%s", conf.Server.Port)
		return conf, err
	}

	log.Debugf("Conf file: Parsed IPAddress=%s and Port=%s", conf.Server.IPAddress, conf.Server.Port)

	return conf, nil
}

// InitConf Init the config from conf file
func InitConf() error {

	conf, err := parseConfFile()
	if err != nil {
		log.Fatalf("Failed to read conf file of '%s'. Using defaults: %v", ConfFile, err)
		flag.Parse()
	} else {
		IPFlag = conf.Server.IPAddress
		PortFlag = conf.Server.Port
	}

	return nil
}
