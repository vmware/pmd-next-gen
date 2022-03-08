// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package firewall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/nftables"
	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type Table struct {
	Family string `json:"Family"`
	Name   string `json:Name"`
}

type Nft struct {
	Table Table `json:"Table"`
}

const (
	nftFilePath = "/etc/nftables-pmd-nextgen.conf"
)

func decodeNftJSONRequest(r *http.Request) (*Nft, error) {
	n := Nft{}
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		return nil, err
	}

	return &n, nil
}

func newConnection() nftables.Conn {
	return nftables.Conn{}
}

func acquireTables() ([]*nftables.Table, error) {
	c := newConnection()
	return c.ListTables()
}

func createMapKey(name, family string) string {
	return name + "_" + family
}

func addressFamilyStringToByte(f string) nftables.TableFamily {
	var family nftables.TableFamily
	switch f {
	case "inet":
		family = unix.NFPROTO_INET
	case "ipv4":
		family = unix.NFPROTO_IPV4
	case "ipv6":
		family = unix.NFPROTO_IPV6
	case "arp":
		family = unix.NFPROTO_ARP
	case "netdev":
		family = unix.NFPROTO_NETDEV
	case "bridge":
		family = unix.NFPROTO_BRIDGE
	}

	return family
}

func addressFamilyByteToString(f nftables.TableFamily) string {
	var family string
	switch f {
	case unix.NFPROTO_INET:
		family = "inet"
	case unix.NFPROTO_IPV4:
		family = "ipv4"
	case unix.NFPROTO_IPV6:
		family = "ipv6"
	case unix.NFPROTO_ARP:
		family = "arp"
	case unix.NFPROTO_NETDEV:
		family = "netdev"
	case unix.NFPROTO_BRIDGE:
		family = "bridge"
	}

	return family
}

func (n *Nft) AddTable(w http.ResponseWriter) error {
	if validator.IsEmpty(n.Table.Name) {
		log.Errorf("Failed to add nft table, Missing table name")
		return fmt.Errorf("missing table name")
	}

	if !validator.IsEmpty(n.Table.Family) {
		if !validator.IsNFTFamily(n.Table.Family) {
			log.Errorf("Failed to add nft table, Invalid family")
			return fmt.Errorf("Invalid family")
		}
	} else {
		n.Table.Family = "ipv4"
	}

	c := newConnection()

	c.AddTable(&nftables.Table{
		Name:   n.Table.Name,
		Family: addressFamilyStringToByte(n.Table.Family),
	})

	if err := c.Flush(); err != nil {
		log.Errorf("Unable to flush connection %v", err)
		return err
	}

	return web.JSONResponse("added", w)
}

func (n *Nft) ShowTable(w http.ResponseWriter) error {
	tables, err := acquireTables()
	if err != nil {
		log.Errorf("Failed to get nft tables: %v", err)
		return err
	}

	tableMap := make(map[string]Table)
	for _, t := range tables {
		tt := Table{
			Name:   t.Name,
			Family: addressFamilyByteToString(t.Family),
		}

		key := createMapKey(tt.Name, tt.Family)
		tableMap[key] = tt
	}

	if !validator.IsEmpty(n.Table.Name) && !validator.IsEmpty(n.Table.Family) {
		key := createMapKey(n.Table.Name, n.Table.Family)
		v, ok := tableMap[key]
		if ok {
			result := make(map[string]Table)
			result[n.Table.Name] = v
			return web.JSONResponse(result, w)
		} else {
			return fmt.Errorf("Table not found='%s'", n.Table.Name)
		}
	}

	return web.JSONResponse(tableMap, w)
}

func (n *Nft) SaveTable(w http.ResponseWriter) error {
	stdout, err := system.ExecAndCapture("nftt", "list", "ruleset")
	if err != nil {
		log.Errorf("Failed to get command output=%v", err)
		return fmt.Errorf("Failed to get command output=%v", err)
	}

	if err := ioutil.WriteFile(nftFilePath, []byte(stdout), 0644); err != nil {
		log.Errorf("Failed to save table info: %v", err)
		return err
	}

	return web.JSONResponse("saved", w)
}

func (n *Nft) RemoveTable(w http.ResponseWriter) error {
	if validator.IsEmpty(n.Table.Name) {
		log.Errorf("Failed to remove nft table, Missing table name")
		return fmt.Errorf("missing table name")
	}

	if !validator.IsEmpty(n.Table.Family) {
		if !validator.IsNFTFamily(n.Table.Family) {
			log.Errorf("Failed to remove nft table, Invalid family")
			return fmt.Errorf("Invalid family")
		}
	} else {
		n.Table.Family = "ipv4"
	}

	c := newConnection()

	c.DelTable(&nftables.Table{
		Name:   n.Table.Name,
		Family: addressFamilyStringToByte(n.Table.Family),
	})

	if err := c.Flush(); err != nil {
		log.Errorf("Unable to flush connection %v", err)
		return err
	}

	return web.JSONResponse("removed", w)
}
