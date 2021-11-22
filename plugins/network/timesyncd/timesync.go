// SPDX-License-Identifier: Apache-2.0

package timesyncd

import (
	"context"
	"net/http"

	"github.com/pm-web/pkg/web"
)

type NTPServer struct {
	ServerName       string   `json:"ServerName`
	Family           int32    `json:"Family`
	ServerAddress    string   `json:"ServerAddress"`
	SystemNTPServers []string `json:"SystemNTPServers"`
}

func AcquireCurrentNTPServer(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusAcquireCurrentNTPServerFromTImeSync(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireSystemNTPServers(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusAcquireSystemNTPServerFromTImeSync(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}
