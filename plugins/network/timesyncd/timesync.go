// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timesyncd

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type NTPServer struct {
	NTPServerName     string   `json:"NTPServerName `
	NTPServerIpFamily int32    `json:"NTPServerIpFamily`
	ServerAddress     string   `json:"ServerAddress"`
	SystemNTPServers  []string `json:"SystemNTPServers"`
	LinkNTPServers    []string `json:"LinkNTPServers"`
}

func AcquireNTPServer(kind string, ctx context.Context) (*NTPServer, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	var s *NTPServer
	switch kind {
	case "currentntpserver":
		s, err = c.DBusAcquireCurrentNTPServerFromTimeSync(ctx)
	case "systemntpservers":
		s, err = c.DBusAcquireSystemNTPServerFromTimeSync(ctx)
	case "linkntpservers":
		s, err = c.DBusAcquireLinkNTPServerFromTimeSync(ctx)
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}
