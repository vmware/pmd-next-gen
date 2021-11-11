// SPDX-License-Identifier: Apache-2.0

package group

import (
	"fmt"
	"net/http"
	"os/exec"
	"os/user"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/web"
)

type Group struct {
	Gid     string `json:"Gid"`
	Name    string `json:"Name"`
	NewName string `json:"NewName"`
}

func (r *Group) GroupAdd(w http.ResponseWriter) error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}
	if g != nil {
		return fmt.Errorf("group %s already exists", r.Name)
	}

	id, err := user.LookupGroupId(r.Gid)
	if err != nil {
		_, ok := err.(user.UnknownGroupIdError)
		if !ok {
			return err
		}
	}
	if id != nil {
		return fmt.Errorf("group %s: Gid %s already exists", r.Name, r.Gid)
	}

	path, err := exec.LookPath("groupadd")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name, "-g", r.Gid)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to add group %s: %s", r.Name, stdout)
		return fmt.Errorf("group '%s': %s", r.Name, stdout)
	}

	return web.JSONResponse("group added", w)
}

func (r *Group) GroupRemove(w http.ResponseWriter) error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("group %s does not exists", r.Name)
	}

	path, err := exec.LookPath("groupdel")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, r.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to delete group %s: %s", r.Name, stdout)
		return fmt.Errorf("group '%s': %s", r.Name, stdout)
	}

	return web.JSONResponse("group removed", w)
}

func (r *Group) GroupModify(w http.ResponseWriter) error {
	g, err := user.LookupGroup(r.Name)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}
	if g == nil {
		return fmt.Errorf("existing group '%s' does not exists", r.Name)
	}

	g, err = user.LookupGroup(r.NewName)
	if err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}
	if g != nil {
		return fmt.Errorf("new group '%s' does not exists", r.NewName)
	}

	path, err := exec.LookPath("groupmod")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-n", r.NewName, r.Name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to modify group %s: %s", r.Name, stdout)
		return fmt.Errorf("group '%s': %s", r.Name, stdout)
	}

	return web.JSONResponse("group modified", w)
}
