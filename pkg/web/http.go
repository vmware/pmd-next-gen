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

type Response struct {
	Body       []byte
	Status     string
	StatusCode int
	Header     http.Header
}

const (
	defaultRequestTimeout = 5 * time.Second
)

func decodeHttpResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse body")
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

func DispatchSocketWithStatus(method, host string, url string, headers map[string]string, data interface{}) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var httpClient *http.Client
	if validator.IsEmpty(host) {
		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", conf.UnixDomainSocketPath)
				},
			},
		}
		url = "http://localhost" + url
	} else {
		httpClient = http.DefaultClient
		url = host + url
	}

	req, err := buildHttpRequest(ctx, method, url, headers, data)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete HTTP request")
	}
	defer resp.Body.Close()

	var body []byte
	if resp.StatusCode == 200 {
		body, err = decodeHttpResponse(resp)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Body:       body,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, nil
}

func DispatchSocket(method, host string, url string, headers map[string]string, data interface{}) ([]byte, error) {
	r, err := DispatchSocketWithStatus(method, host, url, headers, data)

	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}
	return r.Body, err
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
