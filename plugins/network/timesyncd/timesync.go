// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/web"
)

type NTPServer struct {
	ServerName       string   `json:"ServerName`
	Family           int32    `json:"Family`
	ServerAddress    string   `json:"ServerAddress"`
	SystemNTPServers []string `json:"SystemNTPServers"`
	LinkNTPServers   []string `json:"LinkNTPServers"`
}

func AcquireNTPServer(kind string, ctx context.Context, w http.ResponseWriter) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	var s interface{}
	switch kind {
	case "currentntpserver":
		s, err = c.DBusAcquireCurrentNTPServerFromTimeSync(ctx)
	case "systemntpservers":
		s, err = c.DBusAcquireSystemNTPServerFromTimeSync(ctx)
	case "linkntpservers":
		s, err = c.DBusAcquireLinkNTPServerFromTimeSync(ctx)
	}

	if err != nil {
		return err
	}

	return web.JSONResponse(s, w)
}
