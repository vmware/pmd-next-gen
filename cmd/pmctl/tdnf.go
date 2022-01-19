// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
//	"strings"
	"github.com/fatih/color"

	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/tdnf"
)

type ItemListDesc struct {
	Success bool		`json:"success"`
	Message []tdnf.ListItem `json:"message"`
	Errors  string		`json:"errors"`
}

func displayTdnfList(l *ItemListDesc) {
	for _, i := range l.Message {
		fmt.Printf("%v %v\n", color.HiBlueString("Name:"), i.Name)
		fmt.Printf("%v %v\n", color.HiBlueString("Arch:"), i.Arch)
		fmt.Printf("%v %v\n", color.HiBlueString("Evr:"), i.Evr)
		fmt.Printf("%v %v\n", color.HiBlueString("Repo:"), i.Repo)
		fmt.Printf("\n")
	}
}

func acquireTdnfList(host string, token map[string]string) (*ItemListDesc, error) {
	var resp []byte
	var err error
	if host != "" {
		resp, err = web.DispatchSocket(http.MethodGet, host+"/api/v1/tdnf/list", token, nil)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost/api/v1/tdnf/list", nil)
	}
	if err != nil {
		fmt.Printf("tdnf command failed: %v\n", err)
		return nil, err
	}

	m := ItemListDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func commandTdnf(cmd string, host string, token map[string]string) {
	switch cmd {
	case "list":
		l, err := acquireTdnfList(host, token)
		if err != nil {
	        	fmt.Printf("Failed to fetch tdnf list: %v\n", err)
	        	return
		}
		displayTdnfList(l)
	}
}

