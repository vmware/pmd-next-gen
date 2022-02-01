// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package ethtool

import (
	"net/http"

	"github.com/pmd-nextgen/pkg/parser"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Ethtool struct {
	Action   string `json:"action"`
	Link     string `json:"link"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

func (r *Ethtool) AcquireEthTool(w http.ResponseWriter) error {
	_, err := netlink.LinkByName(r.Link)
	if err != nil {
		log.Errorf("Failed to find link='%s': %v", r.Link, err)
		return err
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link='%s': %v", r.Link, err)
		return err
	}
	defer e.Close()

	switch r.Action {
	case "statistics":
		stats, err := e.Stats(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool statitics for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(stats, w)

	case "features":
		features, err := e.Features(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool features for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(features, w)

	case "bus":
		bus, err := e.BusInfo(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool bus for link='%s': %v", r.Link, err)
			return err
		}

		b := struct {
			Bus string
		}{
			bus,
		}

		return web.JSONResponse(b, w)

	case "drivername":
		driver, err := e.DriverName(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool driver name for link='%s': %v", r.Link, err)
			return err
		}

		d := struct {
			Driver string
		}{
			driver,
		}

		return web.JSONResponse(d, w)

	case "driverinfo":
		d, err := e.DriverInfo(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool driver name for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(d, w)

	case "permaddr":
		permaddr, err := e.PermAddr(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool Perm Addr for link='%s': %v", r.Link, err)
			return err
		}

		p := struct {
			PermAddr string
		}{
			permaddr,
		}

		return web.JSONResponse(p, w)

	case "eeprom":
		eeprom, err := e.ModuleEepromHex(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool eeprom for link='%s': %v", r.Link, err)
			return err
		}

		e := struct {
			ModuleEeprom string
		}{
			eeprom,
		}

		return web.JSONResponse(e, w)

	case "msglvl":
		msglvl, err := e.MsglvlGet(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool msglvl for link='%s': %v", r.Link, err)
			return err
		}

		g := struct {
			ModuleMsglv uint32
		}{
			msglvl,
		}

		return web.JSONResponse(g, w)

	case "mapped":
		mapped, err := e.CmdGetMapped(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool msglvl for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(mapped, w)

	case "channels":
		c, err := e.GetChannels(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool channels for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(c, w)

	case "coalesce":
		c, err := e.GetCoalesce(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool coalesce for link='%s': %v", r.Link, err)
			return err
		}

		return web.JSONResponse(c, w)

	case "linkstate":
		c, err := e.LinkState(r.Link)
		if err != nil {
			log.Errorf("Failed to acquire ethtool linkstate for link='%s': %v", r.Link, err)
			return err
		}

		g := struct {
			LinkState uint32
		}{
			c,
		}

		return web.JSONResponse(g, w)
	}

	return nil
}

func (r *Ethtool) ConfigureEthTool(w http.ResponseWriter) error {
	_, err := netlink.LinkByName(r.Link)
	if err != nil {
		log.Errorf("Failed to find link='%s': %v", r.Link, err)
		return err
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		log.Errorf("Failed to init ethtool for link='%s': %v", r.Link, err)
		return err
	}
	defer e.Close()

	switch r.Action {
	case "setfeature":
		feature := make(map[string]bool)

		b, err := parser.ParseBool(r.Value)
		if err != nil {
			return err
		}

		feature[r.Property] = b
		if err := e.Change(r.Link, feature); err != nil {
			return err
		}
	}

	return nil
}
