// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/management/login"
)

type LoginSessionStats struct {
	Success bool            `json:"success"`
	Message []login.Session `json:"message"`
	Errors  string          `json:"errors"`
}

type LoginUserStats struct {
	Success bool         `json:"success"`
	Message []login.User `json:"message"`
	Errors  string       `json:"errors"`
}

func acquireLoginUserStatus(host string, token map[string]string) {
	resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/system/login/listusers", token, nil)
	if err != nil {
		fmt.Printf("Failed to acquire login user info: %v\n", err)
		return
	}

	u := LoginUserStats{}
	if err := json.Unmarshal(resp, &u); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !u.Success {
		fmt.Printf("Failed to acquire login user info: %v\n", err)
		return
	}

	for _, usr := range u.Message {
		fmt.Printf("           %v %v\n", color.HiBlueString("Uid:"), usr.UID)
		fmt.Printf("          %v %v\n", color.HiBlueString("Name:"), usr.Name)
		fmt.Printf("          %v %v\n\n", color.HiBlueString("Path:"), usr.Path)
	}
}

func acquireLoginSessionStatus(host string, token map[string]string) {
	resp, err := web.DispatchSocket(http.MethodGet, host, "/api/v1/system/login/listsessions", token, nil)
	if err != nil {
		fmt.Printf("Failed to acquire login session info: %v\n", err)
		return
	}

	s := LoginSessionStats{}
	if err := json.Unmarshal(resp, &s); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !s.Success {
		fmt.Printf("Failed to acquire login session info: %v\n", err)
		return
	}

	for _, session := range s.Message {
		fmt.Printf("            %v %v\n", color.HiBlueString("Id:"), session.ID)
		fmt.Printf("           %v %v\n", color.HiBlueString("Uid:"), session.UID)
		fmt.Printf("          %v %v\n", color.HiBlueString("User:"), session.User)
		fmt.Printf("          %v %v\n", color.HiBlueString("Seat:"), session.Seat)
		fmt.Printf("          %v %v\n\n", color.HiBlueString("Path:"), session.Path)
	}
}
