// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/pm-web/pkg/conf"
)

func Fetch(url string, headers map[string]string) ([]byte, error) {
	httpClient, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		httpClient.Header.Set(k, v)
	}

	httpClient.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpClient)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Dispatch(method string, url string, headers map[string]string, data interface{}) ([]byte, error) {
	j := new(bytes.Buffer)
	err := json.NewEncoder(j).Encode(data)
	if err != nil {
		return nil, err
	}

	httpClient, err := http.NewRequest(method, url, j)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		httpClient.Header.Set(k, v)
	}

	httpClient.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpClient)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchUnixDomainSocket(url string) ([]byte, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", conf.UnixDomainSocketPath)
			},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err == nil && resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %+v", resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DispatchUnixDomainSocket(method string, url string, data interface{}) ([]byte, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", conf.UnixDomainSocketPath)
			},
		},
	}

	j := new(bytes.Buffer)
	err := json.NewEncoder(j).Encode(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, j)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err == nil && resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %+v", resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
