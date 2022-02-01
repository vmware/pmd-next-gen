// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package networkd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pmd-nextgen/pkg/configfile"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	log "github.com/sirupsen/logrus"
)

type Link struct {
	Link         string       `json:"Link"`
	MatchSection MatchSection `json:"MatchSection"`

	// [Link]
	Name                                string `json:"Name"`
	Alias                               string `json:"Alias"`
	MACAddress                          string `json:"MACAddress"`
	Description                         string `json:"Description"`
	TransmitQueues                      uint   `json:"TransmitQueues"`
	ReceiveQueues                       uint   `json:"ReceiveQueues"`
	TransmitQueueLength                 uint   `json:"TransmitQueueLength"`
	MTUBytes                            string `json:"MTUBytes"`
	BitsPerSecond                       string `json:"BitsPerSecond"`
	Duplex                              string `json:"Duplex"`
	AutoNegotiation                     string `json:"AutoNegotiation"`
	WakeOnLan                           string `json:"WakeOnLan"`
	WakeOnLanPassword                   string `json:"WakeOnLanPassword"`
	Port                                string `json:"Port"`
	Advertise                           string `json:"Advertise"`
	ReceiveChecksumOffload              string `json:"ReceiveChecksumOffload"`
	TransmitChecksumOffload             string `json:"TransmitChecksumOffload"`
	TCPSegmentationOffload              string `json:"TCPSegmentationOffload"`
	TCP6SegmentationOffload             string `json:"TCP6SegmentationOffload"`
	GenericSegmentationOffload          string `json:"GenericSegmentationOffload"`
	GenericReceiveOffload               string `json:"GenericReceiveOffload"`
	GenericReceiveOffloadHardware       string `json:"GenericReceiveOffloadHardware"`
	LargeReceiveOffload                 string `json:"LargeReceiveOffload"`
	ReceiveVLANCTAGHardwareAcceleration string `json:"ReceiveVLANCTAGHardwareAcceleration"`
	RxChannels                          string `json:"RxChannels"`        // range 1…4294967295 or "max
	TxChannels                          string `json:"TxChannels"`        // range 1…4294967295 or "max
	OtherChannels                       string `json:"OtherChannels"`     // range 1…4294967295 or "max
	CombinedChannels                    string `json:"CombinedChannels"`  // range 1…4294967295 or "max
	RxBufferSize                        string `json:"RxBufferSize"`      // range 1…4294967295 or "max
	RxMiniBufferSize                    string `json:"RxMiniBufferSize"`  // range 1…4294967295 or "max
	TxBufferSize                        string `json:"TxBufferSize"`      // range 1…4294967295 or "max
	RxJumboBufferSize                   string `json:"RxJumboBufferSize"` // range 1…4294967295 or "max
	RxFlowControl                       string `json:"RxFlowControl"`
	AutoNegotiationFlowControl          string `json:"TAutoNegotiationFlowControl"`
	GenericSegmentOffloadMaxBytes       uint   `json:"GenericSegmentOffloadMaxBytes"`
	GenericSegmentOffloadMaxSegments    uint   `json:"GenericSegmentOffloadMaxSegments"`
}

func decodeLinkJSONRequest(r *http.Request) (*Link, error) {
	l := Link{}
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Link) BuildLinkSection(m *configfile.Meta) error {
	if !validator.IsEmpty(l.Description) {
		m.SetKeySectionString("Link", "Description", l.Description)
	}

	if !validator.IsEmpty(l.MACAddress) {
		if validator.IsNotMAC(l.MACAddress) {
			log.Errorf("Failed to create .link. Invalid MACAddress='%s': %v", l.Name)
			return fmt.Errorf("invalid MACAddress='%s'", l.MACAddress)
		}

		m.SetKeySectionString("Link","MACAddress", l.MACAddress)
	}

	if l.TransmitQueues > 0 {
		m.SetKeySectionUint("Link","TransmitQueues", l.TransmitQueues)
	}

	if !validator.IsEmpty(l.RxChannels) {
		if !validator.IsUintOrMax(l.RxChannels) {
			log.Errorf("Failed to create .link. Invalid RxChannels='%s': %v", l.Name)
			return fmt.Errorf("invalid RxChannels='%s'", l.RxChannels)
		}

		m.SetKeySectionString("Link","RxChannels", l.RxChannels)
	}

	if !validator.IsEmpty(l.TxChannels) {
		if !validator.IsUintOrMax(l.TxChannels) {
			log.Errorf("Failed to create .link. Invalid TxChannels='%s': %v", l.Name)
			return fmt.Errorf("invalid TxChannels='%s'", l.TxChannels)
		}

		m.SetKeySectionString("Link","TxChannels", l.TxChannels)
	}

	if !validator.IsEmpty(l.OtherChannels) {
		if !validator.IsUintOrMax(l.OtherChannels) {
			log.Errorf("Failed to create .link. Invalid OtherChannels='%s': %v", l.Name)
			return fmt.Errorf("invalid OtherChannels='%s'", l.OtherChannels)
		}

		m.SetKeySectionString("Link","OtherChannels", l.OtherChannels)
	}

	if !validator.IsEmpty(l.CombinedChannels) {
		if !validator.IsUintOrMax(l.CombinedChannels) {
			log.Errorf("Failed to create .link. Invalid CombinedChannels='%s': %v", l.Name)
			return fmt.Errorf("invalid CombinedChannels='%s'", l.CombinedChannels)
		}

		m.SetKeySectionString("Link","CombinedChannels", l.CombinedChannels)
	}

	if !validator.IsEmpty(l.RxBufferSize) {
		if !validator.IsUintOrMax(l.RxBufferSize) {
			log.Errorf("Failed to create .link. Invalid RxBufferSize='%s': %v", l.Name)
			return fmt.Errorf("invalid RxBufferSize='%s'", l.RxBufferSize)
		}

		m.SetKeySectionString("Link","RxBufferSize", l.RxBufferSize)
	}

	if !validator.IsEmpty(l.RxMiniBufferSize) {
		if !validator.IsUintOrMax(l.RxMiniBufferSize) {
			log.Errorf("Failed to create .link. Invalid RxMiniBufferSize='%s': %v", l.Name)
			return fmt.Errorf("invalid RxMiniBufferSize='%s'", l.RxMiniBufferSize)
		}

		m.SetKeySectionString("Link","RxMiniBufferSize", l.RxMiniBufferSize)
	}

	if !validator.IsEmpty(l.TxBufferSize) {
		if !validator.IsUintOrMax(l.TxBufferSize) {
			log.Errorf("Failed to create .link. Invalid TxBufferSize='%s': %v", l.Name)
			return fmt.Errorf("invalid TxBufferSize='%s'", l.TxBufferSize)
		}

		m.SetKeySectionString("Link","TxBufferSize", l.TxBufferSize)
	}

	if !validator.IsEmpty(l.RxJumboBufferSize) {
		if !validator.IsUintOrMax(l.RxJumboBufferSize) {
			log.Errorf("Failed to create .link. Invalid RxJumboBufferSize='%s': %v", l.Name)
			return fmt.Errorf("invalid RxJumboBufferSize='%s'", l.RxJumboBufferSize)
		}

		m.SetKeySectionString("Link","RxJumboBufferSize", l.RxJumboBufferSize)
	}

	if l.GenericSegmentOffloadMaxBytes > 0 {
		m.SetKeySectionUint("Link","GenericSegmentOffloadMaxBytes", l.GenericSegmentOffloadMaxBytes)
	}

	if l.GenericSegmentOffloadMaxSegments > 0 {
		m.SetKeySectionUint("Link", "GenericSegmentOffloadMaxSegments", l.GenericSegmentOffloadMaxSegments)
	}

	return nil
}

func (l *Link) ConfigureLink(ctx context.Context, w http.ResponseWriter) error {
	m, err := CreateOrParseLinkFile(l.Link)
	if err != nil {
		return err
	}

	if err := l.BuildLinkSection(m); err != nil {
		return err
	}

	if err := m.Save(); err != nil {
		log.Errorf("Failed to update config file='%s': %v", m.Path, err)
		return err
	}

	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection with the system bus: %v", err)
		return err
	}
	defer c.Close()

	if err := c.DBusNetworkReload(ctx); err != nil {
		return err
	}

	return web.JSONResponse("configured", w)
}
