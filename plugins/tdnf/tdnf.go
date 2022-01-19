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
	"os/exec"
	"encoding/json"

	log "github.com/sirupsen/logrus"

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

func TdnfExec(cmd string) (string, error) {
	s, err := system.ExecAndCapture("tdnf", "-j", cmd)
	if err != nil {
		werr := err.(*exec.ExitError)
		log.Errorf("tdnf returned %d\n", werr.Error())
	}
	return s, err
}

func AcquireList(w http.ResponseWriter) error {
	s, err := TdnfExec("list");
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var listData []ListItem
	json.Unmarshal([]byte(s), &listData)
	return web.JSONResponse(listData, w)
}

func AcquireRepoList(w http.ResponseWriter) error {
	s, err := TdnfExec("repolist");
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var repoList []Repo
	json.Unmarshal([]byte(s), &repoList)
	return web.JSONResponse(repoList, w)
}

