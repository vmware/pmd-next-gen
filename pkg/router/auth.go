// SPDX-License-Identifier: Apache-2.0

package router

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/pm-web/pkg/system"
)

const (
	authConfPath = "/etc/pm-web/pmweb-auth.conf"
)

type TokenDB struct {
	tokenUsers map[string]string
}

func (db *TokenDB) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("X-Session-Token")

		if user, found := db.tokenUsers[token]; found {
			log.Printf("Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			log.Infof("Unauthorized user")
		}

		next.ServeHTTP(w, r)
	})
}

func UnixDomainPeerCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var credentialsContextKey = struct{}{}

		credentials := r.Context().Value(credentialsContextKey).(*unix.Ucred)

		log.Infof("Connection credentials: Pid=%v, Uid=%v, Gid=%v", credentials.Pid, credentials.Gid, credentials.Uid)

		next.ServeHTTP(w, r)
	})
}

func InitAuthMiddleware() (TokenDB, error) {
	db := TokenDB{make(map[string]string)}

	lines, r := system.ReadFullFile(authConfPath)
	if r != nil {
		log.Fatal("Failed to read auth config file")
		return db, errors.New("Failed to read auth config file")
	}

	for _, line := range lines {
		authLine := strings.Fields(line)
		db.tokenUsers[authLine[1]] = authLine[0]
	}

	return db, nil
}
