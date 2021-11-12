// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/web"
)

type Hostname struct {
	Method string `json:"Method"`
	Value  string `json:"Value"`
}

func (h *Hostname) SetHostname(ctx context.Context) error {
	conn, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	return conn.ExecuteHostNameMethod(ctx, h.Method, h.Value)
}

func AcquireHostnameProperties(ctx context.Context, w http.ResponseWriter) error {
	conn, err := NewSDConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to the system bus: %s", err)
		return err
	}
	defer conn.Close()

	if p, err := conn.AcquireHostNameProperty(ctx, "Hostname"); err != nil {
		return err
	} else {
		return web.JSONResponse(p, w)
	}
}
