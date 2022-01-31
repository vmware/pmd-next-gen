// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pmd-nextgen/pkg/jobs"
	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
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
	job := jobs.CreateJob(func() (string, error) {
		var s string
		var err error
		if !validator.IsEmpty(pkg) {
			s, err = TdnfExec("list", pkg)
		} else {
			s, err = TdnfExec("list")
		}
		return s, err
	})
	return jobs.AcceptedResponse(w, job)
}

func AcquireRepoList(w http.ResponseWriter) error {
	s, err := TdnfExec("repolist")
	if err != nil {
		log.Errorf("Failed to execute tdnf repolist: %v", err)
		return err
	}

	var repoList interface{}
	json.Unmarshal([]byte(s), &repoList)

	return web.JSONResponse(repoList, w)
}

func AcquireInfoList(w http.ResponseWriter, pkg string) error {
	job := jobs.CreateJob(func() (string, error) {
		var s string
		var err error
		if !validator.IsEmpty(pkg) {
			s, err = TdnfExec("info", pkg)
		} else {
			s, err = TdnfExec("info")
		}
		return s, err
	})
	return jobs.AcceptedResponse(w, job)
}

func AcquireMakeCache(w http.ResponseWriter) error {
	job := jobs.CreateJob(func() (string, error) {
		s, err := TdnfExec("makecache")
		return s, err
	})
	return jobs.AcceptedResponse(w, job)
}

func AcquireClean(w http.ResponseWriter) error {
	_, err := TdnfExec("clean", "all")
	if err != nil {
		log.Errorf("Failed to execute tdnf clean all': %v", err)
		return err
	}

	return web.JSONResponse("cleaned", w)
}
