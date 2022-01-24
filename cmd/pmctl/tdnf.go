// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"

	"github.com/pmd-nextgen/pkg/web"
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

func DispatchSocket(method, host string, url string, token map[string]string, body io.Reader) ([]byte, error) {
	var resp []byte
	var err error
	if host != "" {
		resp, err = web.DispatchSocket(method, host+url, token, body)
	} else {
		resp, err = web.DispatchUnixDomainSocket(method, "http://localhost"+url, body)
	}
	return resp, err
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
	if pkg != "" {
		path = "/api/v1/tdnf/list/" + pkg
	} else {
		path = "/api/v1/tdnf/list"
	}
	resp, err := DispatchSocket(http.MethodGet, host, path, token, nil)
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

func acquireTdnfRepoList(host string, token map[string]string) (*RepoListDesc, error) {
	resp, err := DispatchSocket(http.MethodGet, host, "/api/v1/tdnf/repolist", token, nil)
	if err != nil {
		fmt.Printf("tdnf command failed: %v\n", err)
		return nil, err
	}

	m := RepoListDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
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

	resp, err := DispatchSocket(http.MethodGet, host, path, token, nil)
	if err != nil {
		fmt.Printf("tdnf command failed: %v\n", err)
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
	resp, err := DispatchSocket(http.MethodGet, host, "/api/v1/tdnf/"+cmd, token, nil)
	if err != nil {
		fmt.Printf("tdnf command failed: %v\n", err)
		return nil, err
	}

	m := NilDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
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
		fmt.Printf("Failed tdnf clean: %v\n", err)
		return
	}
}

func tdnfList(pkg string, host string, token map[string]string) {
	l, err := acquireTdnfList(pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to fetch tdnf list: %v\n", err)
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
}

func tdnfRepoList(host string, token map[string]string) {
	l, err := acquireTdnfRepoList(host, token)
	if err != nil {
		fmt.Printf("Failed to fetch tdnf repolist: %v\n", err)
		return
	}
	displayTdnfRepoList(l)
}

func tdnfInfoList(pkg string, host string, token map[string]string) {
	l, err := acquireTdnfInfoList(pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to fetch tdnf info: %v\n", err)
		return
	}
	displayTdnfInfoList(l)
}
