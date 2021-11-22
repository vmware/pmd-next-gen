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

type Domains struct {
	Link   string `json:"Link"`
	Domain string `json:"DNS"`
}

func AcquireLinkDNS(ctx context.Context, link string, w http.ResponseWriter) error {
	links, err := DBusAcquireDNSFromResolveLink(ctx, link)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireLinkDomains(ctx context.Context, link string, w http.ResponseWriter) error {
	links, err := DBusAcquireDomainsFromResolveLink(ctx, link)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireDNSFromResolveManager(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusAcquireDNSFromResolveManager(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireDomainsFromResolveManager(ctx context.Context, w http.ResponseWriter) error {
	links, err := DBusAcquireDomainsFromResolveManager(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}
