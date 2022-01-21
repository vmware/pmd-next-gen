// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/management/sysctl"
)

type SysctlStats struct {
	Success bool              `json:"success"`
	Message map[string]string `json:"message"`
	Errors  string            `json:"errors"`
}

func acquireSysctlParamStatus(key string, host string, token map[string]string) {
	var resp []byte
	var err error
	url := "/api/v1/system/sysctl/status"

	s := sysctl.Sysctl{
		Key: key,
	}

	if !validator.IsEmpty(host) {
		resp, err = web.DispatchSocket(http.MethodGet, host+url, token, s)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost"+url, s)
	}
	if err != nil {
		fmt.Printf("Failed to acquire sysctl info: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to fetch sysctl status: %v\n", m.Errors)
		return
	}

	fmt.Printf("             %v %v\n", color.HiBlueString(key+": "), m.Message)
}

func acquireSysctlStatus(urlSuffix string, pattern string, host string, token map[string]string) {
	var resp []byte
	var err error
	url := "/api/v1/system/sysctl/"

	s := sysctl.Sysctl{
		Pattern: pattern,
	}
	url = url + urlSuffix

	if !validator.IsEmpty(host) {
		resp, err = web.DispatchSocket(http.MethodGet, host+url, token, s)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodGet, "http://localhost"+url, s)
	}
	if err != nil {
		fmt.Printf("Failed to acquire sysctl info: %v\n", err)
		return
	}

	ss := SysctlStats{}
	if err := json.Unmarshal(resp, &ss); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !ss.Success {
		fmt.Printf("Failed to fetch sysctl status: %v\n", ss.Errors)
		return
	}

	jsonStr, err := json.Marshal(ss.Message)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("%v\n", color.HiBlueString(string(jsonStr)))
	}
}

func sysctlUpdateConfig(key string, value string, filename string, host string, token map[string]string) {
	var resp []byte
	var err error

	s := sysctl.Sysctl{
		Key:      key,
		Value:    value,
		FileName: filename,
		Apply:    true,
	}

	if !validator.IsEmpty(host) {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/system/sysctl/update", token, s)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/system/sysctl/update", s)
	}
	if err != nil {
		fmt.Printf("Failed to update sysctl configuration: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to update sysctl configuration: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}

func sysctlRemoveConfig(key string, filename string, host string, token map[string]string) {
	var resp []byte
	var err error

	s := sysctl.Sysctl{
		Key:      key,
		FileName: filename,
		Apply:    true,
	}

	if !validator.IsEmpty(host) {
		resp, err = web.DispatchSocket(http.MethodDelete, host+"/api/v1/system/sysctl/remove", token, s)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodDelete, "http://localhost/api/v1/system/sysctl/remove", s)
	}
	if err != nil {
		fmt.Printf("Failed to remove sysctl configuration: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to remove sysctl configuration: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}

func sysctlLoadConfig(files string, host string, token map[string]string) {
	var resp []byte
	var err error

	s := sysctl.Sysctl{
		Apply: true,
	}
	if !validator.IsEmpty(files) {
		s.Files = strings.Split(files, ",")
	}
	if !validator.IsEmpty(host) {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/system/sysctl/load", token, s)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/system/sysctl/load", s)
	}
	if err != nil {
		fmt.Printf("Failed to load sysctl configuration: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to load sysctl configuration: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}
