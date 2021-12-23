// SPDX-License-Identifier: Apache-2.0

package server

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
	"golang.org/x/sys/unix"

	"github.com/distro-management-api/pkg/conf"
	"github.com/distro-management-api/pkg/share"
	"github.com/distro-management-api/pkg/system"
	"github.com/distro-management-api/plugins/management"
	"github.com/distro-management-api/plugins/network"
	"github.com/distro-management-api/plugins/proc"
	"github.com/distro-management-api/plugins/systemd"
)

var httpSrv *http.Server

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	systemd.InitSystemd()
	systemd.RegisterRouterSystemd(s)

	management.RegisterRouterManagement(s)
	network.RegisterRouterNetwork(s)

	proc.RegisterRouterProc(s)

	return r
}

func runUnixDomainHttpServer(c *conf.Config, r *mux.Router) error {
	var credentialsContextKey = struct{}{}

	if c.System.UseAuthentication {
		r.Use(UnixDomainPeerCredential)
	}

	httpSrv = &http.Server{
		Handler: r,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			file, _ := c.(*net.UnixConn).File()
			credentials, _ := unix.GetsockoptUcred(int(file.Fd()), unix.SOL_SOCKET, unix.SO_PEERCRED)
			return context.WithValue(ctx, credentialsContextKey, credentials)
		},
	}

	log.Infof("Starting distro-management-apid server at unix domain socket='%s' in HTTP mode pid=%d", conf.UnixDomainSocketPath, os.Getpid())

	os.Remove(conf.UnixDomainSocketPath)
	unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: conf.UnixDomainSocketPath, Net: "unix"})
	if err != nil {
		log.Fatalf("Unable to listen on unix domain socket='%s': %v", conf.UnixDomainSocketPath, err)
	}
	defer unixListener.Close()

	if err := system.ChangeUnixDomainSocketPermission(conf.UnixDomainSocketPath); err != nil {
		log.Errorf("Failed to change unix domain socket permissions: %v", err)
		return err
	}

	log.Fatal(httpSrv.Serve(unixListener))

	return nil
}

func runWebHttpServer(c *conf.Config, r *mux.Router) error {
	if c.System.UseAuthentication {
		amw, err := InitAuthMiddleware()
		if err != nil {
			log.Fatalf("Failed to init auth DB existing: %v", err)
			return err
		}

		r.Use(amw.AuthMiddleware)
	}

	ip, port, _ := share.ParseIpPort(c.Network.Listen)

	if system.TLSFilePathExits() {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: false,
		}
		httpSrv = &http.Server{
			Addr:         ip + ":" + port,
			Handler:      r,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}

		log.Infof("Starting distro-management-apid server at %s:%s in HTTPS mode pid=%d", ip, port, os.Getpid())

		httpSrv.ListenAndServeTLS(path.Join(conf.ConfPath, conf.TLSCert), path.Join(conf.ConfPath, conf.TLSKey))
	} else {
		httpSrv = &http.Server{
			Addr:    ip + ":" + port,
			Handler: r,
		}

		log.Infof("Starting distro-management-apid server at %s:%s in HTTP mode pid=%d", ip, port, os.Getpid())

		httpSrv.ListenAndServe()
	}

	return nil
}

func Run(c *conf.Config) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Printf("Signal received='%v'. Shutting down distro-management-apid ...", sig)

			if err := httpSrv.Shutdown(ctx); err != nil {
				os.Exit(1)
			}

			cancel()
		case <-ctx.Done():
		}
	}()

	r := NewRouter()
	if c.Network.ListenUnixSocket {
		runUnixDomainHttpServer(c, r)
	} else {
		runWebHttpServer(c, r)
	}

	return nil
}
