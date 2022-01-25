// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
	"encoding/json"
	"net/http"

	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"

	log "github.com/sirupsen/logrus"
)

type ListItem struct {
	Name string `json:"Name"`
	Arch string `json:"Arch"`
	Evr  string `json:"Evr"`
	Repo string `json:"Repo"`
}

type Repo struct {
	Repo     string `json:"Repo"`
	RepoName string `json:"RepoName"`
	Enabled  bool   `json:"Enabled"`
}

type Info struct {
	Name        string `json:"Name"`
	Arch        string `json:"Arch"`
	Evr         string `json:"Evr"`
	InstallSize int    `json:"InstallSize"`
	Repo        string `json:"Repo"`
	Summary     string `json:"Summary"`
	Url         string `json:"Url"`
	License     string `json:"License"`
	Description string `json:"Description"`
}

func TdnfExec(args ...string) (string, error) {
	args = append([]string{"-j"}, args...)
	s, err := system.ExecAndCapture("tdnf", args...)
	if err != nil {
		return "", err
	}

	return s, nil
}

func AcquireList(w http.ResponseWriter, pkg string) error {
	var s string
	var err error
	if !validator.IsEmpty(pkg) {
		log.Errorf("Failed execute tdnf list for '%s': %v", pkg, err)
		s, err = TdnfExec("list", pkg)
	} else {
		log.Errorf("Failed execute tdnf list: %v", pkg, err)
		s, err = TdnfExec("list")
	}
	if err != nil {
		return err
	}

	var listData []ListItem
	json.Unmarshal([]byte(s), &listData)

	return web.JSONResponse(listData, w)
}

func AcquireRepoList(w http.ResponseWriter) error {
	s, err := TdnfExec("repolist")
	if err != nil {
		log.Errorf("Failed execute tdnf repolist: %v", err)
		return err
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
		log.Errorf("Failed execute tdnf info for '%s': %v", pkg, err)

	} else {
		s, err = TdnfExec("info")
	}
	if err != nil {
		log.Errorf("Failed execute tdnf info: %v", err)
		return err
	}

	var infoList []Info
	json.Unmarshal([]byte(s), &infoList)

	return web.JSONResponse(infoList, w)
}

func AcquireMakeCache(w http.ResponseWriter) error {
	_, err := TdnfExec("makecache")
	if err != nil {
		log.Errorf("Failed execute tdnf makecache': %v", err)
		return err
	}

	return web.JSONResponse("successfully executed makecache", w)
}

func AcquireClean(w http.ResponseWriter) error {
	_, err := TdnfExec("clean", "all")
	if err != nil {
		log.Errorf("Failed execute tdnf clean all': %v", err)
		return err
	}

	return web.JSONResponse("cleaned", w)
}
