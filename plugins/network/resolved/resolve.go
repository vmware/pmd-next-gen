// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package resolved

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"

	"github.com/pmd-nextgen/pkg/web"
)

type Dns struct {
	Index  int32  `json:"Index"`
	Link   string `json:"Link"`
	Family int32  `json:"Family"`
	Dns    string `json:"Dns"`
}

type Domains struct {
	Index  int32  `json:"Index"`
	Link   string `json:"Link"`
	Domain string `json:"Domain"`
}

type Describe struct {
	CurrentDNS string    `json:"CurrentDNS"`
	DnsServers []Dns     `json:"DnsServers"`
	Domains    []Domains `json:"Domains"`
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


func AcquireLinkCurrentDns(ctx context.Context, link string, w http.ResponseWriter) error {
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

	links, err := c.DBusAcquireCurrentDnsFromResolveLink(ctx, l.Attrs().Index)
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

func DescribeDns(ctx context.Context) (*Describe, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return nil, err
	}
	defer c.Close()

	d := Describe{}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		dns, err := c.DBusAcquireDnsFromResolveManager(ctx)
		if err == nil {
			d.DnsServers = dns
		}
	}()

	go func() {
		defer wg.Done()
		domains, err := c.DBusAcquireDomainsFromResolveManager(ctx)
		if err == nil {
			d.Domains = domains
		}
	}()

	wg.Wait()
	return &d, nil
}
