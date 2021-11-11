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
	Uid           string   `json:"uid"`
	Gid           string   `json:"gid"`
	Groups        []string `json:"groups"`
	Comment       string   `json:"comment"`
	HomeDirectory string   `json:"home_dir"`
	Shell         string   `json:"shell"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
}

func (r *User) Add(w http.ResponseWriter) error {
	u, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if u != nil {
		return fmt.Errorf("user '%s' already exists", r.Username)
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
			return fmt.Errorf("user '%s': gid '%s' already exists", r.Username, r.Gid)
		}
	}

	// <Username>:<Password>:<UID>:<GID>:<User Info>:<Home Dir>:<Default Shell>
	line := r.Username + ":" + r.Password + ":" + r.Uid + ":" + r.Gid + ":" + r.Comment + ":" + r.HomeDirectory + ":" + r.Shell
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
		log.Errorf("Failed to add user %s: %s", r.Username, stdout)
		return fmt.Errorf("failed to add user '%s': %s", r.Username, stdout)
	}

	return web.JSONResponse("user added", w)
}

func (r *User) Remove(w http.ResponseWriter) error {
	g, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("user '%s' does not exists", r.Username)
	}

	path, err := exec.LookPath("userdel")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Username)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to delete user %s: %s", r.Username, stdout)
		return fmt.Errorf("user '%s': %s", r.Username, stdout)
	}

	return web.JSONResponse("user removed", w)
}

func (r *User) Modify(w http.ResponseWriter) error {
	g, err := user.Lookup(r.Username)
	if err != nil {
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("user %s does not exists", r.Username)
	}

	path, err := exec.LookPath("usermod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-G", r.Groups[0], r.Username)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to modify user %s: %s", r.Username, stdout)
		return fmt.Errorf("unable to modify user %s: %s", r.Username, stdout)
	}

	return web.JSONResponse("user modified", w)
}
