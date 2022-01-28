// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/pmd-nextgen/pkg/conf"
	"github.com/pmd-nextgen/pkg/validator"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

func decodeHttpResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not parse body")
	}

	return body, nil
}

func buildHttpRequest(ctx context.Context, method string, url string, headers map[string]string, data interface{}) (*http.Request, error) {
	j := new(bytes.Buffer)
	if err := json.NewEncoder(j).Encode(data); err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, method, url, j)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpRequest.Header.Set(k, v)
	}

	return httpRequest, nil
}

func DispatchWebSocket(method string, url string, headers map[string]string, data interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	req, err := buildHttpRequest(ctx, method, url, headers, data)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete HTTP request")
	}
	defer resp.Body.Close()

	return decodeHttpResponse(resp)
}

func DispatchUnixDomainSocket(method string, url string, data interface{}) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", conf.UnixDomainSocketPath)
			},
		},
	}

	req, err := buildHttpRequest(ctx, method, url, nil, data)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete HTTP request")
	}

	return decodeHttpResponse(resp)
}

func DispatchSocket(method, host string, url string, token map[string]string, data interface{}) ([]byte, error) {
	var resp []byte
	var err error

	if !validator.IsEmpty(host) {
		resp, err = DispatchWebSocket(method, host+url, token, data)
	} else {
		resp, err = DispatchUnixDomainSocket(method, "http://localhost"+url, data)
	}
	return resp, err
}

func BuildAuthTokenFromEnv() (map[string]string, error) {
	token := os.Getenv("PHOTON_MGMT_AUTH_TOKEN")
	if token == "" {
		return nil, errors.New("authentication token not found")
	}

	headers := make(map[string]string)
	headers["X-Session-Token"] = token

	return headers, nil
}
