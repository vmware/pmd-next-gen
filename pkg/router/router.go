// SPDX-License-Identifier: Apache-2.0

package router

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/share"
	"github.com/pmd/pkg/systemd"
)

// StartRouter Init and start Gorilla mux router
func StartRouter(ip string, port string, tlsCertPath string, tlsKeyPath string) error {
	var srv http.Server

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	// Register services

	systemd.InitSystemd()
	systemd.RegisterRouterSystemd(s)

	// Authenticate users
	amw, err := InitAuthMiddleware()
	if err != nil {
		log.Fatalf("Failed to init auth DB existing: %s", err)
		return fmt.Errorf("Failed to init Auth DB: %s", err)
	}

	r.Use(amw.AuthMiddleware)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop

		log.Printf("Received signal: %+v", sig)
		log.Println("Shutting down pm-webd ...")

		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Failed to shutdown server gracefully: %s", err)
		}

		os.Exit(0)
	}()

	if share.PathExists(tlsCertPath) && share.PathExists(tlsKeyPath) {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
		}
		srv = http.Server{
			Addr:         ip + ":" + port,
			Handler:      r,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Info("Starting pm-webd in TLS mode")

		log.Fatal(srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath))
	} else {
		srv = http.Server{
			Addr:    ip + ":" + port,
			Handler: r,
		}

		log.Info("Starting pm-webd in plain text mode")

		log.Fatal(srv.ListenAndServe())
	}

	return nil
}
