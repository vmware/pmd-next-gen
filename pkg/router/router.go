// SPDX-License-Identifier: Apache-2.0

package router

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/conf"
	"github.com/pmd/pkg/proc"
	"github.com/pmd/pkg/system"
	"github.com/pmd/pkg/systemd"
)

func StartRouter(c *conf.Config) error {
	var srv http.Server

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	systemd.InitSystemd()
	systemd.RegisterRouterSystemd(s)
	proc.RegisterRouterProc(s)

	if c.System.UseAuthentication {
		amw, err := InitAuthMiddleware()
		if err != nil {
			log.Fatalf("Failed to init auth DB existing: %v", err)
			return err
		}

		r.Use(amw.AuthMiddleware)
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop

		log.Printf("Received signal: %+v", sig)
		log.Println("Shutting down pm-webd ...")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Errorf("Failed to shutdown server gracefully: %v", err)
		}

		os.Exit(0)
	}()

	if c.Network.ListenUnixSocket {
		system.CreateDirectory("/run/pmwebd/", 0755)

		server := http.Server{
			Handler: r,
		}

		log.Infof("Starting pm-webd server at unix domain socket '/run/pmwebd/pmwebd.sock' in HTTP mode")

		os.Remove("/run/pmwebd/pmwebd.sock")

		unixListener, err := net.Listen("unix", "/run/pmwebd/pmwebd.sock")
		if err != nil {
			log.Fatalf("Unable to listen on unix domain socket file '/run/pmwebd/pmwebd.sock': %v", err)
		}

		defer unixListener.Close()

		log.Fatal(server.Serve(unixListener))

	} else {
		if system.PathExists(path.Join(conf.ConfPath, conf.TLSCert)) && system.PathExists(path.Join(conf.ConfPath, conf.TLSKey)) {
			cfg := &tls.Config{
				MinVersion:               tls.VersionTLS12,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				PreferServerCipherSuites: false,
			}
			srv = http.Server{
				Addr:         c.Network.IPAddress + ":" + c.Network.Port,
				Handler:      r,
				TLSConfig:    cfg,
				TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
			}

			log.Infof("Starting pm-webd server at %s:%s in HTTPS mode", c.Network.IPAddress, c.Network.Port)

			log.Fatal(srv.ListenAndServeTLS(path.Join(conf.ConfPath, conf.TLSCert), path.Join(conf.ConfPath, conf.TLSKey)))
		} else {
			srv = http.Server{
				Addr:    c.Network.IPAddress + ":" + c.Network.Port,
				Handler: r,
			}

			log.Infof("Starting pm-webd server at %s:%s in HTTP mode", c.Network.IPAddress, c.Network.Port)

			log.Fatal(srv.ListenAndServe())
		}
	}

	return nil
}
