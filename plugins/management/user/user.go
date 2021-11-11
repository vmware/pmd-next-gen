// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
	u, err := user.Lookup(r.UserName)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if u != nil {
		return fmt.Errorf("user '%s' already exists", r.UserName)
	}

	if r.Uid != "" {
		id, err := user.LookupId(r.Uid)
		if err != nil {
			_, ok := err.(user.UnknownUserIdError)
			if !ok {
				return err
			}
		}
		if id != nil {
			return fmt.Errorf("user '%s': gid '%s' already exists", r.UserName, r.Gid)
		}
	}

	// <UserName>:<Password>:<UID>:<GID>:<User Info>:<Home Dir>:<Default Shell>
	line := r.UserName + ":" + r.Password + ":" + r.Uid + ":" + r.Gid + ":" + r.Comment + ":" + r.HomeDirectory + ":" + r.Shell
	if err := system.WriteOneLineFile(userFile, line); err != nil {
		return err
	}
	defer os.Remove(userFile)

	path, err := exec.LookPath("newusers")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, userFile)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to add user %s: %s", r.UserName, stdout)
		return fmt.Errorf("failed to add user '%s': %s", r.UserName, stdout)
	}

	return web.JSONResponse("user added", w)
}

func (r *User) Remove(w http.ResponseWriter) error {
	g, err := user.Lookup(r.UserName)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("user '%s' does not exists", r.UserName)
	}

	path, err := exec.LookPath("userdel")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.UserName)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to delete user %s: %s", r.UserName, stdout)
		return fmt.Errorf("user '%s': %s", r.UserName, stdout)
	}

	return web.JSONResponse("user removed", w)
}

func (r *User) Modify(w http.ResponseWriter) error {
	g, err := user.Lookup(r.UserName)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("user %s does not exists", r.UserName)
	}

	path, err := exec.LookPath("usermod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-G", r.Groups[0], r.UserName)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to modify user %s: %s", r.UserName, stdout)
		return fmt.Errorf("unable to modify user %s: %s", r.UserName, stdout)
	}

	return web.JSONResponse("user modified", w)
}
