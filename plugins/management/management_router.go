// SPDX-License-Identifier: Apache-2.0

package management

import (
	"github.com/gorilla/mux"

	"github.com/pm-web/plugins/management/group"
	"github.com/pm-web/plugins/management/hostname"
	"github.com/pm-web/plugins/management/login"
	"github.com/pm-web/plugins/management/user"
)

func RegisterRouterManagement(router *mux.Router) {
	n := router.PathPrefix("/system").Subrouter()

	group.RegisterRouterGroup(n)
	user.RegisterRouterUser(n)

	hostname.RegisterRouterHostname(n)
	login.RegisterRouterLogin(n)
}
