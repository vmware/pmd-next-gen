// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/pm-web/pkg/conf"
)

const (
	defaultRequestTimeout = 5 * time.Second
)

func decodeHttpResponse(resp *http.Response, err error) ([]byte, error) {
	if resp.StatusCode != 200 {
		return nil, errors.Wrap(err, "non-200 status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode body")
	}

	return body, nil
}

func buildHttpRequest(ctx context.Context, method string, url string, headers map[string]string, data interface{}) (*http.Request, error){
	j := new(bytes.Buffer)
	if err := json.NewEncoder(j).Encode(data);err != nil {
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

func DispatchSocket(method string, url string, headers map[string]string, data interface{}) ([]byte, error) {
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

	return decodeHttpResponse(resp, err)
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

	return decodeHttpResponse(resp, err)
}
