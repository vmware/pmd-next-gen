package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/distro-management-api/pkg/web"
	"github.com/distro-management-api/plugins/management/group"
	"github.com/distro-management-api/plugins/management/user"
)

func groupAdd(name string, gid string, host string, token map[string]string) {
	var resp []byte
	var err error

	g := group.Group{
		Name: name,
		Gid:  gid,
	}

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/system/group/add", token, g)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/system/group/add", g)
	}
	if err != nil {
		fmt.Printf("Failed add group: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to add group: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}

func groupRemove(name string, gid string, host string, token map[string]string) {
	var resp []byte
	var err error

	g := group.Group{
		Name: name,
		Gid:  gid,
	}

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodDelete, host+"/api/v1/system/group/remove", token, g)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodDelete, "http://localhost/api/v1/system/group/remove", g)
	}
	if err != nil {
		fmt.Printf("Failed remove group: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to remove group: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}

func userAdd(name string, uid string, groups string, gid string, shell string, homeDir string, password, string, host string, token map[string]string) {
	var resp []byte
	var err error

	u := user.User{
		Name:          name,
		Uid:           uid,
		Gid:           gid,
		Groups:        strings.Split("groups", ","),
		Shell:         shell,
		HomeDirectory: homeDir,
		Password:      password,
	}

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodPost, host+"/api/v1/system/user/add", token, u)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodPost, "http://localhost/api/v1/system/user/add", u)
	}
	if err != nil {
		fmt.Printf("Failed add user: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to add user: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}

func userRemove(name string, host string, token map[string]string) {
	var resp []byte
	var err error

	u := user.User{
		Name: name,
	}

	if host != "" {
		resp, err = web.DispatchSocket(http.MethodDelete, host+"/api/v1/system/user/add", token, u)
	} else {
		resp, err = web.DispatchUnixDomainSocket(http.MethodDelete, "http://localhost/api/v1/system/user/add", u)
	}
	if err != nil {
		fmt.Printf("Failed remove user: %v\n", err)
		return
	}

	m := web.JSONResponseMessage{}
	if err := json.Unmarshal(resp, &m); err != nil {
		fmt.Printf("Failed to decode json message: %v\n", err)
		os.Exit(1)
	}

	if !m.Success {
		fmt.Printf("Failed to remove user: %v\n", m.Errors)
		os.Exit(1)
	}

	fmt.Println(m.Message)
}
