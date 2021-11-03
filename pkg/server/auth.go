// SPDX-License-Identifier: Apache-2.0

package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/pm-web/pkg/share"
	"github.com/pm-web/pkg/system"
	"github.com/pm-web/pkg/web"
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

func authenticateLocalUser(credentials *unix.Ucred) error {
	if credentials.Uid != 0 {
		pmUser, err := system.GetUserCredentials("pm-web")
		if err != nil {
			log.Infof("Failed to get user 'pm-web' credentials: %+v", err)
			return err
		}

		u, _ := system.GetUserCredentialsByUid(credentials.Uid)

		groups, _ := u.GroupIds()
		if !share.StringContains(groups,strconv.Itoa(int(pmUser.Gid))) {
			return errors.New("user's gid not same as pm-web's gid")
		}

		log.Infof("Connection credentials: pid=%v, user='%s' uid=%v, gid=%v belongs to groups='%v'", credentials.Pid, u.Username, credentials.Gid, credentials.Uid, groups)
	}

	return nil
}

func UnixDomainPeerCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var credentialsContextKey = struct{}{}

		credentials := r.Context().Value(credentialsContextKey).(*unix.Ucred)

		if err := authenticateLocalUser(credentials); err != nil {
			web.JSONResponseError(err, w)
			log.Infof("Unauthorized connection. Credentials: pid=%v, uid=%v, gid=%v: %v", credentials.Pid, credentials.Gid, credentials.Uid, err)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
