// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timedate

import (
	"net/http"
	"sync"

	"github.com/pmd-nextgen/pkg/web"
	log "github.com/sirupsen/logrus"
)

type TimeDate struct {
	Method string `json:"Method"`
	Value  string `json:"Value"`
}

type Describe struct {
	Timezone        string `json:"Timezone"`
	LocalRTC        bool   `json:"LocalRTC"`
	CanNTP          bool   `json:"CanNTP"`
	NTP             string `json:"NTP"`
	NTPSynchronized bool   `json:"NTPSynchronized"`
	TimeUSec        uint64 `json:"TimeUSec"`
	RTCTimeUSec     uint64 `json:"RTCTimeUSec"`
}

func (t *TimeDate) ConfigureTimeDate(w http.ResponseWriter) error {
	conn, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %v", err)
		return err
	}
	defer conn.Close()

	err = conn.dBusConfigureTimeDate(t.Method, t.Value)
	if err != nil {
		log.Errorf("Failed to set timedate property: %s", err)
		return err
	}

	web.JSONResponse("configured", w)
	return nil
}

func AcquireTimeDate(w http.ResponseWriter) error {
	c, err := NewSDConnection()
	if err != nil {
		return err
	}
	defer c.Close()

	h := Describe{}

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("Timezone")
		if err == nil {
			h.Timezone = s.Value().(string)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("LocalRTC")
		if err == nil {
			h.LocalRTC = s.Value().(bool)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("CanNTP")
		if err == nil {
			h.CanNTP = s.Value().(bool)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("NTPSynchronized")
		if err == nil {
			h.NTPSynchronized = s.Value().(bool)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("TimeUSec")
		if err == nil {
			h.TimeUSec = s.Value().(uint64)
		}
	}()

	go func() {
		defer wg.Done()
		s, err := c.dbusAcquire("RTCTimeUSec")
		if err == nil {
			h.RTCTimeUSec = s.Value().(uint64)
		}
	}()

	wg.Wait()

	web.JSONResponse(h, w)

	return nil
}
