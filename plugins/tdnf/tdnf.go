// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/web"
)

type ListItem struct {
	Name string `json:"name"`
	Arch string `json:"arch"`
	Evr  string `json:"evr"`
	Repo string `json:"repo"`
}

type Repo struct {
	Repo     string `json:"repo"`
	RepoName string `json:"repo_name"`
	Enabled  bool   `json:"enabled"`
}

type Info struct {
	Name        string `json:"name"`
	Arch        string `json:"arch"`
	Evr         string `json:"evr"`
	InstallSize int    `json:"install_size"`
	Repo        string `json:"repo"`
	Summary     string `json:"summary"`
	Url         string `json:"url"`
	License     string `json:"license"`
	Description string `json:"description"`
}

func TdnfExec(args ...string) (string, error) {
	args = append([]string{"-j"}, args...)
	s, err := system.ExecAndCapture("tdnf", args...)
	if err != nil {
		werr := err.(*exec.ExitError)
		log.Errorf("tdnf returned %d\n", werr.Error())
	}
	return s, err
}

func AcquireList(w http.ResponseWriter, pkg string) error {
	var s string
	var err error
	if pkg != "" {
		s, err = TdnfExec("list", pkg)
	} else {
		s, err = TdnfExec("list")
	}
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var listData []ListItem
	json.Unmarshal([]byte(s), &listData)
	return web.JSONResponse(listData, w)
}

func AcquireRepoList(w http.ResponseWriter) error {
	s, err := TdnfExec("repolist")
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var repoList []Repo
	json.Unmarshal([]byte(s), &repoList)
	return web.JSONResponse(repoList, w)
}

func AcquireInfoList(w http.ResponseWriter, pkg string) error {
	var s string
	var err error
	if pkg != "" {
		s, err = TdnfExec("info", pkg)
	} else {
		s, err = TdnfExec("info")
	}
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	var infoList []Info
	json.Unmarshal([]byte(s), &infoList)
	return web.JSONResponse(infoList, w)
}

func AcquireMakeCache(w http.ResponseWriter) error {
	s, err := TdnfExec("makecache")
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	return web.JSONResponse(nil, w)
}

func AcquireClean(w http.ResponseWriter) error {
	s, err := TdnfExec("clean", "all")
	if err != nil {
		return fmt.Errorf("tdnf failed: '%s'", s)
	}
	return web.JSONResponse(nil, w)
}
