// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/network/networkd"
	"github.com/urfave/cli/v2"
)

func networkConfigureLinkQueue(args cli.Args, host string, token map[string]string) {
	argStrings := args.Slice()
	
	l := networkd.Link{}
	for i := 0; i < len(argStrings); {
		switch argStrings[i] {
		case "dev":
			l.Link = argStrings[i+1]
		case "txqueue":
			if !validator.IsLinkQueue(argStrings[i+1]) {
				fmt.Printf("Failed to set link txqueue: Invalid txqueue=%s\n", argStrings[i+1])
				return
			}
			n, _ := strconv.ParseUint(argStrings[i+1], 10, 32)
			l.TransmitQueues= uint(n)
		}

		i++
	}

	if validator.IsEmpty(l.Link) {
		fmt.Printf("Failed to set link. Missing link name")
		return
	}

	resp, err := web.DispatchSocket(http.MethodPost, host, "/api/v1/network/networkd/link/configure", token, l)
	if err != nil {
		fmt.Printf("Failed to set link: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		return
	}

	if !m.Success {
		fmt.Printf("Failed to set link: %v\n", m.Errors)
	}
}
