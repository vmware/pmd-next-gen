// SPDX-License-Identifier: Apache-2.0

package router

import (
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/pmd/pkg/system"
)

const (
	authConfPath = "/etc/pm-web/pmweb-auth.conf"
)

// TokenDB token DB
type TokenDB struct {
	tokenUsers map[string]string
}

// AuthMiddleware Authenticate the User
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
	})
}

// InitAuthMiddleware init middleware
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
