// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"

	"github.com/pmd-nextgen/pkg/jobs"
	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/plugins/tdnf"
)

type ItemListDesc struct {
	Success bool            `json:"success"`
	Message []tdnf.ListItem `json:"message"`
	Errors  string          `json:"errors"`
}

type RepoListDesc struct {
	Success bool        `json:"success"`
	Message []tdnf.Repo `json:"message"`
	Errors  string      `json:"errors"`
}

type InfoListDesc struct {
	Success bool        `json:"success"`
	Message []tdnf.Info `json:"message"`
	Errors  string      `json:"errors"`
}

type NilDesc struct {
	Success bool   `json:"success"`
	Errors  string `json:"errors"`
}

type StatusDesc struct {
	Success bool                `json:"success"`
	Message jobs.StatusResponse `json:"message"`
	Errors  string              `json:"errors"`
}

func displayTdnfList(l *ItemListDesc) {
	for _, i := range l.Message {
		fmt.Printf("%v %v\n", color.HiBlueString("Name:"), i.Name)
		fmt.Printf("%v %v\n", color.HiBlueString("Arch:"), i.Arch)
		fmt.Printf("%v %v\n", color.HiBlueString(" Evr:"), i.Evr)
		fmt.Printf("%v %v\n", color.HiBlueString("Repo:"), i.Repo)
		fmt.Printf("\n")
	}
}

func displayTdnfRepoList(l *RepoListDesc) {
	for _, r := range l.Message {
		fmt.Printf("%v %v\n", color.HiBlueString("   Repo:"), r.Repo)
		fmt.Printf("%v %v\n", color.HiBlueString("   Name:"), r.RepoName)
		fmt.Printf("%v %v\n", color.HiBlueString("Enabled:"), r.Enabled)
		fmt.Printf("\n")
	}
}

func displayTdnfInfoList(l *InfoListDesc) {
	for _, i := range l.Message {
		fmt.Printf("%v %v\n", color.HiBlueString("        Name:"), i.Name)
		fmt.Printf("%v %v\n", color.HiBlueString("        Arch:"), i.Arch)
		fmt.Printf("%v %v\n", color.HiBlueString("         Evr:"), i.Evr)
		fmt.Printf("%v %v\n", color.HiBlueString("Install Size:"), i.InstallSize)
		fmt.Printf("%v %v\n", color.HiBlueString("        Repo:"), i.Repo)
		fmt.Printf("%v %v\n", color.HiBlueString("     Summary:"), i.Summary)
		fmt.Printf("%v %v\n", color.HiBlueString("         Url:"), i.Url)
		fmt.Printf("%v %v\n", color.HiBlueString("     License:"), i.License)
		fmt.Printf("%v %v\n", color.HiBlueString(" Description:"), i.Description)
		fmt.Printf("\n")
	}
}

func acquireTdnfList(pkg string, host string, token map[string]string) (*ItemListDesc, error) {
	var path string
	if !validator.IsEmpty(pkg) {
		path = "/api/v1/tdnf/list/" + pkg
	} else {
		path = "/api/v1/tdnf/list"
	}
	resp, err := jobs.DispatchAndWait(http.MethodGet, host, path, token, nil)
	if err != nil {
		return nil, err
	}

	m := ItemListDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func acquireTdnfRepoList(host string, token map[string]string) (*RepoListDesc, error) {
	resp, err := jobs.DispatchAndWait(http.MethodGet, host, "/api/v1/tdnf/repolist", token, nil)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf repolist: %v\n", err)
		return nil, err
	}

	m := RepoListDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func acquireTdnfInfoList(pkg string, host string, token map[string]string) (*InfoListDesc, error) {
	var path string
	if pkg != "" {
		path = "/api/v1/tdnf/info/" + pkg
	} else {
		path = "/api/v1/tdnf/info"
	}

	resp, err := jobs.DispatchAndWait(http.MethodGet, host, path, token, nil)
	if err != nil {
		return nil, err
	}

	m := InfoListDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func acquireTdnfSimpleCommand(cmd string, host string, token map[string]string) (*NilDesc, error) {
	var msg []byte

	msg, err := jobs.DispatchAndWait(http.MethodGet, host, "/api/v1/tdnf/"+cmd, token, nil)
	if err != nil {
		return nil, err
	}

	m := NilDesc{}
	if err := json.Unmarshal(msg, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func tdnfClean(host string, token map[string]string) {
	_, err := acquireTdnfSimpleCommand("clean", host, token)
	if err != nil {
		fmt.Printf("Failed execute tdnf clean: %v\n", err)
		return
	}
	fmt.Printf("package cache cleaned\n")
}

func tdnfList(pkg string, host string, token map[string]string) {
	l, err := acquireTdnfList(pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf list: %v\n", err)
		return
	}
	displayTdnfList(l)
}

func tdnfMakeCache(host string, token map[string]string) {
	_, err := acquireTdnfSimpleCommand("makecache", host, token)
	if err != nil {
		fmt.Printf("Failed tdnf makecache: %v\n", err)
		return
	}
	fmt.Printf("package cache acquired\n")
}

func tdnfRepoList(host string, token map[string]string) {
	l, err := acquireTdnfRepoList(host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf repolist: %v\n", err)
		return
	}
	displayTdnfRepoList(l)
}

func tdnfInfoList(pkg string, host string, token map[string]string) {
	l, err := acquireTdnfInfoList(pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf info: %v\n", err)
		return
	}
	displayTdnfInfoList(l)
}
