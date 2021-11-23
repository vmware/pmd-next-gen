// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

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
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	l, err := netlink.LinkByName(link)
	if err != nil {
		return err
	}

	links, err := c.DBusAcquireDNSFromResolveLink(ctx, l.Attrs().Index)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireLinkDomains(ctx context.Context, link string, w http.ResponseWriter) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	l, err := netlink.LinkByName(link)
	if err != nil {
		return err
	}

	links, err := c.DBusAcquireDomainsFromResolveLink(ctx, l.Attrs().Index)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireDNSFromResolveManager(ctx context.Context, w http.ResponseWriter) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	links, err := c.DBusAcquireDNSFromResolveManager(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}

func AcquireDomainsFromResolveManager(ctx context.Context, w http.ResponseWriter) error {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer c.Close()

	links, err := c.DBusAcquireDomainsFromResolveManager(ctx)
	if err != nil {
		return web.JSONResponseError(err, w)
	}

	return web.JSONResponse(links, w)
}
