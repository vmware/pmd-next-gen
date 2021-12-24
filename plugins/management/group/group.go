// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package group

import (
	"fmt"
	"net/http"
	"os/user"

	log "github.com/sirupsen/logrus"

	"github.com/distro-management-api/pkg/system"
	"github.com/distro-management-api/pkg/web"
)

type Group struct {
	Gid     string `json:"Gid"`
	Name    string `json:"Name"`
	NewName string `json:"NewName"`
}

func (g *Group) GroupAdd(w http.ResponseWriter) error {
	var grp *user.Group
	var err error

	if grp, err = user.LookupGroup(g.Name); err != nil {
		_, ok := err.(user.UnknownGroupError)
		if !ok {
			return err
		}
	}
	if grp != nil {
		return fmt.Errorf("group %s already exists", grp.Name)
	}

	if g.Gid != "" {
		id, err := user.LookupGroupId(g.Gid)
		if err != nil {
			_, ok := err.(user.UnknownGroupIdError)
			if !ok {
				return err
			}
		}
		if id != nil {
			return fmt.Errorf(" gid '%v' already exists", id.Gid)
		}
	}

	if g.Gid != "" {
		if s, err := system.ExecAndCapture("groupadd", g.Name, "-g", g.Gid); err != nil {
			return fmt.Errorf("%s (%v)", s, err)
		}
	} else {
		if s, err := system.ExecAndCapture("groupadd", g.Name); err != nil {
			return fmt.Errorf("%s (%v)", s, err)
		}
	}
	return web.JSONResponse("group added", w)
}

func (g *Group) GroupRemove(w http.ResponseWriter) error {
	if _, err := system.GetGroupCredentials(g.Name); err != nil {
		return err
	}

	if s, err := system.ExecAndCapture("groupdel", g.Name); err != nil {
		log.Errorf("Failed to remove group '%s': %s (%v)", g.Name, s, err)
		return fmt.Errorf("%s (%v)", s, err)
	}

	return web.JSONResponse("group removed", w)
}

func (g *Group) GroupModify(w http.ResponseWriter) error {
	if _, err := system.GetGroupCredentials(g.Name); err != nil {
		return err
	}

	if s, err := system.ExecAndCapture("groupmod", "-n", g.NewName, g.Name); err != nil {
		log.Errorf("Failed to modify group '%s': %s (%v)", g.Name, s, err)
		return fmt.Errorf("%s (%v)", s, err)
	}

	return web.JSONResponse("group modified", w)
}
