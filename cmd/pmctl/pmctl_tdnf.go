// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
	"github.com/pmd-nextgen/plugins/tdnf"
)

type ItemListDesc struct {
	Success bool            `json:"success"`
	Message []tdnf.ListItem `json:"message"`
	Errors  string          `json:"errors"`
}

type ItemSearchDesc struct {
	Success bool              `json:"success"`
	Message []tdnf.SearchItem `json:"message"`
	Errors  string            `json:"errors"`
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

type AlterResultDesc struct {
	Success bool             `json:"success"`
	Message tdnf.AlterResult `json:"message"`
	Errors  string           `json:"errors"`
}

type NilDesc struct {
	Success bool   `json:"success"`
	Errors  string `json:"errors"`
}

type StatusDesc struct {
	Success bool               `json:"success"`
	Message web.StatusResponse `json:"message"`
	Errors  string             `json:"errors"`
}

func tdnfParseFlagsInterface(c *cli.Context, optType reflect.Type) interface{} {
	options := reflect.New(optType)
	v := options.Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		value := v.Field(i).Interface()
		switch value.(type) {
		case bool:
			v.Field(i).SetBool(c.Bool(name))
		case string:
			v.Field(i).SetString(c.String(name))
		case []string:
			str := c.String(name)
			if !validator.IsEmpty(str) {
				list := strings.Split(str, ",")
				size := len(list)
				if size > 0 {
					v.Field(i).Set(reflect.MakeSlice(reflect.TypeOf([]string{}), size, size))
					for j, s := range list {
						v.Field(i).Index(j).Set(reflect.ValueOf(s))
					}
				}
			}
		}
	}
	return options.Interface()
}

func tdnfParseFlags(c *cli.Context) tdnf.Options {
	var o tdnf.Options
	o = *tdnfParseFlagsInterface(c, reflect.TypeOf(o)).(*tdnf.Options)
	return o
}

func tdnfParseScopeFlags(c *cli.Context) tdnf.ScopeOptions {
	var o tdnf.ScopeOptions
	o = *tdnfParseFlagsInterface(c, reflect.TypeOf(o)).(*tdnf.ScopeOptions)
	return o
}

func tdnfCreateFlagsInterface(optType reflect.Type) []cli.Flag {
	var flags []cli.Flag

	options := reflect.New(optType)
	v := options.Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		value := v.Field(i).Interface()
		switch value.(type) {
		case bool:
			flags = append(flags, &cli.BoolFlag{Name: name})
		case string:
			flags = append(flags, &cli.StringFlag{Name: name})
		case []string:
			flags = append(flags, &cli.StringFlag{Name: name, Usage: "Separate by ,"})
		}
	}
	return flags
}

func tdnfCreateFlags() []cli.Flag {
	var o tdnf.Options
	return tdnfCreateFlagsInterface(reflect.TypeOf(o))
}

func tdnfCreateScopeFlags() []cli.Flag {
	var o tdnf.ScopeOptions
	return tdnfCreateFlagsInterface(reflect.TypeOf(o))
}

func tdnfCreateAlterCommand(cmd string, aliases []string, desc string, pkgRequired bool, token map[string]string) *cli.Command {
	return &cli.Command{
		Name:        cmd,
		Aliases:     aliases,
		Description: desc,

		Action: func(c *cli.Context) error {
			options := tdnfParseFlags(c)
			if c.NArg() >= 1 {
				tdnfAlterCmd(&options, cmd, c.Args().First(), c.String("url"), token)
			} else {
				if pkgRequired {
					fmt.Printf("Needs a package name\n")
					return nil
				}
				tdnfAlterCmd(&options, cmd, "", c.String("url"), token)
			}
			return nil
		},
	}
}

func tdnfOptionsMap(options interface{}) url.Values {
	m := url.Values{}

	v := reflect.ValueOf(options).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		switch v.Field(i).Kind() {
		case reflect.Struct:
			value := v.Field(i).Addr().Interface()
			m1 := tdnfOptionsMap(value)
			for k, v := range m1 {
				m[k] = v
			}
		default:
			value := v.Field(i).Interface()
			switch value.(type) {
			case bool:
				if value.(bool) {
					m.Add(name, "true")
				}
			case string:
				str := value.(string)
				if !validator.IsEmpty(str) {
					m.Add(name, str)
				}
			case []string:
				list := value.([]string)
				if len(list) != 0 {
					m[name] = list
				}
			}
		}
	}
	return m
}

func tdnfOptionsQuery(options interface{}) string {
	if m := tdnfOptionsMap(options); len(m) != 0 {
		return "?" + m.Encode()
	}
	return ""
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

func displayTdnfSearch(l *ItemSearchDesc) {
	for _, i := range l.Message {
		fmt.Printf("%v %v\n", color.HiBlueString("   Name:"), i.Name)
		fmt.Printf("%v %v\n", color.HiBlueString("Summary:"), i.Summary)
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

func displayAlterList(l []tdnf.ListItem, header string) {
	if len(l) > 0 {
		fmt.Printf("%s:\n\n", header)
		for _, i := range l {
			fmt.Printf("%v %v\n", color.HiBlueString("Name:"), i.Name)
			fmt.Printf("%v %v\n", color.HiBlueString("Arch:"), i.Arch)
			fmt.Printf("%v %v\n", color.HiBlueString(" Evr:"), i.Evr)
			fmt.Printf("%v %v\n", color.HiBlueString("Repo:"), i.Repo)
			fmt.Printf("%v %v\n", color.HiBlueString("Size:"), i.Size)
			fmt.Printf("\n")
		}
	}
}

func displayTdnfAlterResult(rDesc *AlterResultDesc) {
	r := rDesc.Message
	displayAlterList(r.Exist, "Existing Packages")
	displayAlterList(r.Unavailable, "Unavailable Packages")
	displayAlterList(r.Install, "Packages to Install")
	displayAlterList(r.Upgrade, "Packages to Upgrade")
	displayAlterList(r.Downgrade, "Packages to Downgrade")
	displayAlterList(r.Remove, "Packages to Remove")
	displayAlterList(r.UnNeeded, "Unneeded Packages")
	displayAlterList(r.Reinstall, "Packages to Reinstall")
	displayAlterList(r.Obsolete, "Packages to be Obsoleted")
}

func acquireTdnfList(options *tdnf.ListOptions, pkg string, host string, token map[string]string) (*ItemListDesc, error) {
	var path string
	if !validator.IsEmpty(pkg) {
		path = "/api/v1/tdnf/list/" + pkg
	} else {
		path = "/api/v1/tdnf/list"
	}
	path = path + tdnfOptionsQuery(options)

	resp, err := web.DispatchAndWait(http.MethodGet, host, path, token, nil)
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

func acquireTdnfRepoList(options *tdnf.Options, host string, token map[string]string) (*RepoListDesc, error) {
	resp, err := web.DispatchAndWait(http.MethodGet, host, "/api/v1/tdnf/repolist"+tdnfOptionsQuery(options), token, nil)
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

func acquireTdnfInfoList(options *tdnf.Options, pkg string, host string, token map[string]string) (*InfoListDesc, error) {
	var path string
	if !validator.IsEmpty(pkg) {
		path = "/api/v1/tdnf/info/" + pkg
	} else {
		path = "/api/v1/tdnf/info"
	}
	path = path + tdnfOptionsQuery(options)

	resp, err := web.DispatchAndWait(http.MethodGet, host, path, token, nil)
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

func acquireTdnfCheckUpdate(options *tdnf.Options, pkg string, host string, token map[string]string) (*ItemListDesc, error) {
	var path string
	if !validator.IsEmpty(pkg) {
		path = "/api/v1/tdnf/check-update/" + pkg
	} else {
		path = "/api/v1/tdnf/check-update"
	}
	path = path + tdnfOptionsQuery(options)

	resp, err := web.DispatchAndWait(http.MethodGet, host, path, token, nil)
	if err != nil {
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

func acquireTdnfSearch(options *tdnf.Options, q string, host string, token map[string]string) (*ItemSearchDesc, error) {
	var path string

	v := tdnfOptionsMap(options)
	v.Add("q", q)
	path = "/api/v1/tdnf/search?" + v.Encode()

	resp, err := web.DispatchAndWait(http.MethodGet, host, path, token, nil)
	if err != nil {
		return nil, err
	}

	m := ItemSearchDesc{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func acquireTdnfSimpleCommand(options *tdnf.Options, cmd string, host string, token map[string]string) (*NilDesc, error) {
	var msg []byte

	msg, err := web.DispatchAndWait(http.MethodGet, host, "/api/v1/tdnf/"+cmd+tdnfOptionsQuery(options), token, nil)
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

func acquireTdnfAlterCmd(options *tdnf.Options, cmd string, pkg string, host string, token map[string]string) (*AlterResultDesc, error) {
	var msg []byte

	msg, err := web.DispatchAndWait(http.MethodGet, host, "/api/v1/tdnf/"+cmd+"/"+pkg+tdnfOptionsQuery(options), token, nil)
	if err != nil {
		return nil, err
	}

	m := AlterResultDesc{}
	if err := json.Unmarshal(msg, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if m.Success {
		return &m, nil
	}

	return nil, errors.New(m.Errors)
}

func tdnfClean(options *tdnf.Options, host string, token map[string]string) {
	_, err := acquireTdnfSimpleCommand(options, "clean", host, token)
	if err != nil {
		fmt.Printf("Failed execute tdnf clean: %v\n", err)
		return
	}
	fmt.Printf("package cache cleaned\n")
}

func tdnfCheckUpdate(options *tdnf.Options, pkg string, host string, token map[string]string) {
	l, err := acquireTdnfCheckUpdate(options, pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire check-update: %v\n", err)
		return
	}
	displayTdnfList(l)
}

func tdnfList(options *tdnf.Options, scOptions *tdnf.ScopeOptions, pkg string, host string, token map[string]string) {
	listOptions := tdnf.ListOptions{*options, *scOptions}
	l, err := acquireTdnfList(&listOptions, pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf list: %v\n", err)
		return
	}
	displayTdnfList(l)
}

func tdnfMakeCache(options *tdnf.Options, host string, token map[string]string) {
	_, err := acquireTdnfSimpleCommand(options, "makecache", host, token)
	if err != nil {
		fmt.Printf("Failed tdnf makecache: %v\n", err)
		return
	}
	fmt.Printf("package cache acquired\n")
}

func tdnfRepoList(options *tdnf.Options, host string, token map[string]string) {
	l, err := acquireTdnfRepoList(options, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf repolist: %v\n", err)
		return
	}
	displayTdnfRepoList(l)
}

func tdnfSearch(options *tdnf.Options, pkg string, host string, token map[string]string) {
	l, err := acquireTdnfSearch(options, pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf search: %v\n", err)
		return
	}
	displayTdnfSearch(l)
}

func tdnfInfoList(options *tdnf.Options, pkg string, host string, token map[string]string) {
	l, err := acquireTdnfInfoList(options, pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf info: %v\n", err)
		return
	}
	displayTdnfInfoList(l)
}

func tdnfAlterCmd(options *tdnf.Options, cmd string, pkg string, host string, token map[string]string) {
	l, err := acquireTdnfAlterCmd(options, cmd, pkg, host, token)
	if err != nil {
		fmt.Printf("Failed to acquire tdnf %s: %v\n", cmd, err)
		return
	}
	displayTdnfAlterResult(l)
}
