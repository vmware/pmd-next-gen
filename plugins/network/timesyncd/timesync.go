// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package timesyncd

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Describe struct {
	Name             string   `json:"Name `
	IpFamily         int32    `json:"IpFamily`
	Address          string   `json:"Address"`
	SystemNTPServers []string `json:"SystemNTPServers"`
	LinkNTPServers   []string `json:"LinkNTPServers"`
}

func AcquireNTPServer(kind string, ctx context.Context) (*Describe, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %v", err)
		return nil, err
	}
	defer c.Close()

	s := Describe{}
	switch kind {
	case "currentntpserver":
		s.Name, s.IpFamily, s.Address, err = c.DBusAcquireCurrentNTPServerFromTimeSync(ctx)
	case "systemntpservers":
		s.SystemNTPServers, err = c.DBusAcquireSystemNTPServersFromTimeSync(ctx)
	case "linkntpservers":
		s.LinkNTPServers, err = c.DBusAcquireLinkNTPServersFromTimeSync(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func DescribeNTPServers(ctx context.Context) (*Describe, error) {
	c, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %v", err)
		return nil, err
	}
	defer c.Close()

	s := Describe{}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		s.Name, s.IpFamily, s.Address, err = c.DBusAcquireCurrentNTPServerFromTimeSync(ctx)
	}()

	go func() {
		defer wg.Done()
		s.SystemNTPServers, err = c.DBusAcquireSystemNTPServersFromTimeSync(ctx)
	}()

	go func() {
		defer wg.Done()
		s.LinkNTPServers, err = c.DBusAcquireLinkNTPServersFromTimeSync(ctx)
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return &s, nil
}
