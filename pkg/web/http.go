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
)

func Fetch(url string, headers map[string]string) ([]byte, error) {
	httpClient, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		httpClient.Header.Set(k, v)
	}

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
				return net.Dial("unix", "/run/pmwebd/pmwebd.sock")
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
				return net.Dial("unix", "/run/pmwebd/pmwebd.sock")
			},
		},
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, buf)
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
