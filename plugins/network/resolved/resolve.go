// SPDX-License-Identifier: Apache-2.0

package resolved

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pm-web/pkg/web"
)

type Dns struct {
	Link   string `json:"Link"`
	Family int32  `json:"Family"`
	Dns    string `json:"Dns"`
}

type Domains struct {
	Link   string `json:"Link"`
	Domain string `json:"Domain"`
}

func AcquireLinkDns(ctx context.Context, link string, w http.ResponseWriter) error {
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

	links, err := c.DBusAcquireDnsFromResolveLink(ctx, l.Attrs().Index)
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

func AcquireDns(ctx context.Context) ([]Dns, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	dns, err := c.DBusAcquireDnsFromResolveManager(ctx)
	if err != nil {
		return nil, err
	}

	return dns, nil
}

func AcquireDomains(ctx context.Context) ([]Domains, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	domains, err := c.DBusAcquireDomainsFromResolveManager(ctx)
	if err != nil {
		return nil, err
	}

	return domains, nil
}
