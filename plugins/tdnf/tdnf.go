// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
//        "context"
//        "errors"
        "net/http"
//        "strconv"
//        "strings"
	"fmt"
	"encoding/json"

        "github.com/pmd-nextgen/pkg/system"
        "github.com/pmd-nextgen/pkg/web"
)

type ListItem struct {
	Name string `json:"name"`
	Arch string `json:"arch"`
	Evr string  `json:"evr"`
	Repo string `json:"repo"`
}

type Repo struct {
	Repo string     `json:"repo"`
	RepoName string `json:"repo_name"`
	Enabled bool    `json:"enabled"`
}

func AcquireList(w http.ResponseWriter) error {
	s, err := system.ExecAndCapture("tdnf", "-j", "list");
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var listData []ListItem
	json.Unmarshal([]byte(s), &listData)
	return web.JSONResponse(listData, w)
}

func AcquireRepoList(w http.ResponseWriter) error {
	s, err := system.ExecAndCapture("tdnf", "-j", "repolist");
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var repoList []Repo
	json.Unmarshal([]byte(s), &repoList)
	return web.JSONResponse(repoList, w)
}

