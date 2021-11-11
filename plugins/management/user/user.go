// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"
	"os"
	"os/user"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/system"
	"github.com/pm-web/pkg/web"
)

const (
	userFile = "/run/pmwebd/users"
)

type User struct {
	Uid           string   `json:"Uid"`
	Gid           string   `json:"Gid"`
	Groups        []string `json:"Groups"`
	Comment       string   `json:"Comment"`
	HomeDirectory string   `json:"HomeDir"`
	Shell         string   `json:"Shell"`
	UserName      string   `json:"UserName"`
	Password      string   `json:"Password"`
}

func (r *User) Add(w http.ResponseWriter) error {
	if _, err := system.GetUserCredentials(r.UserName); err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}

	if r.Uid != "" {
		id, err := user.LookupId(r.Uid)
		if err != nil {
			_, ok := err.(user.UnknownUserError)
			if !ok {
				return err
			}
		}
		if id != nil {
			return fmt.Errorf("user %s gid %s already exists", r.UserName, r.Gid)
		}
	}

	// <UserName>:<Password>:<UID>:<GID>:<User Info>:<Home Dir>:<Default Shell>
	line := r.UserName + ":" + r.Password + ":" + r.Uid + ":" + r.Gid + ":" + r.Comment + ":" + r.HomeDirectory + ":" + r.Shell
	if err := system.WriteOneLineFile(userFile, line); err != nil {
		return err
	}
	defer os.Remove(userFile)

	if s, err := system.ExecAndCapture("newusers", userFile); err != nil {
		log.Errorf("Failed to add user %s: %s (%v)", r.UserName, s, err)
		return fmt.Errorf("failed to add user '%s': %s (%v)", r.UserName, s, err)
	}

	return web.JSONResponse("user added", w)
}

func (r *User) Remove(w http.ResponseWriter) error {
	if _, err := system.GetUserCredentials(r.UserName); err != nil {
		return err
	}

	if s, err := system.ExecAndCapture("userdel", r.UserName); err != nil {
		log.Errorf("Failed to delete user %s: %s (%v)", r.UserName, s)
		return fmt.Errorf("user '%s': %s (%v)", r.UserName, s, err)
	}

	return web.JSONResponse("user removed", w)
}

func (r *User) Modify(w http.ResponseWriter) error {
	if _, err := system.GetUserCredentials(r.UserName); err != nil {
		return err
	}

	if s, err := system.ExecAndCapture("usermod", "-G", r.Groups[0], r.UserName); err != nil {
		log.Errorf("Failed to modify user %s: %s (%v)", r.UserName, s, err)
		return fmt.Errorf("unable to modify user %s: %s (%v)", r.UserName, s, err)
	}

	return web.JSONResponse("user modified", w)
}
