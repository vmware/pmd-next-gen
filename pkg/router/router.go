// SPDX-License-Identifier: Apache-2.0

package router

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/system"
	"github.com/pmd/pkg/systemd"
)

func StartRouter(ip string, port string, tlsCertPath string, tlsKeyPath string) error {
	var srv http.Server

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	// Register services
	systemd.InitSystemd()
	systemd.RegisterRouterSystemd(s)

	// Authenticate users
	amw, err := InitAuthMiddleware()
	if err != nil {
		log.Fatalf("Failed to init auth DB existing: %s", err)
		return err
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

	if system.PathExists(tlsCertPath) && system.PathExists(tlsKeyPath) {
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

		log.Infof("Starting pm-webd server at %s:%s in HTTPS mode", ip, port)

		log.Fatal(srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath))
	} else {
		srv = http.Server{
			Addr:    ip + ":" + port,
			Handler: r,
		}

		log.Infof("Starting pm-webd server at %s:%s in HTTP mode", ip, port)

		log.Fatal(srv.ListenAndServe())
	}

	return nil
}
