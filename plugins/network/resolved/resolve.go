// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"context"
	"net/http"

	"github.com/pm-web/pkg/web"
)

type DNS struct {
	Link   string `json:"Link"`
	Family int32  `json:"Family"`
	DNS    string `json:"DNS"`
}

func AcquireLinkDNS(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusNetworkLinkProperty(ctx)
	if err != nil {
		return err
	}

	return web.JSONResponse(links, w)
}

func AcquireManagerDNS(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusResolveManagerDNS(ctx)
	if err != nil {
		return err
	}

	return web.JSONResponse(links, w)
}
