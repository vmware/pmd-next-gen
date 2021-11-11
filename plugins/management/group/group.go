// SPDX-License-Identifier: Apache-2.0

package group

import (
	"errors"
	"fmt"
	"net/http"
	"os/user"

	log "github.com/sirupsen/logrus"

	"github.com/pm-web/pkg/system"
	"github.com/pm-web/pkg/web"
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
		_, ok := err.(user.UnknownUserError)
		if !ok {
			return err
		}
	}
	if grp != nil {
		return fmt.Errorf("group %s gid %s already exists", grp.Name, grp.Gid)
	}

	if g.Gid != "" {
		id, err := user.LookupGroupId(g.Gid)
		if err != nil {
			_, ok := err.(user.UnknownUserError)
			if !ok {
				return err
			}
		}
		if id != nil {
			return fmt.Errorf("group %s gid %s already exists", g.Name, g.Gid)
		}
	}

	if s, err := system.ExecAndCapture("groupadd", g.Name, "-g", g.Gid); err != nil {
		log.Errorf("Failed to add group %s: %s (%v)", g.Name, s, err)
		return fmt.Errorf("%s (%v)", s, err)
	}

	return web.JSONResponse("group added", w)
}

func (g *Group) GroupRemove(w http.ResponseWriter) error {
	if _, err := system.GetUserCredentials(g.Name); err != nil {
		return err
	}

	if s, err := system.ExecAndCapture("groupdel", g.Name); err != nil {
		log.Errorf("Failed to remove group %s: %s (%v)", g.Name, s, err)
		return fmt.Errorf("%s (%v)", s, err)
	}

	return web.JSONResponse("group removed", w)
}

func (g *Group) GroupModify(w http.ResponseWriter) error {
	if _, err := system.GetUserCredentials(g.Name); err != nil {
		return err
	}

	if g, err := user.LookupGroup(g.NewName); err != nil || g != nil {
		return errors.New("new group exists")
	}

	if s, err := system.ExecAndCapture("groupmod", "-n", g.NewName, g.Name); err != nil {
		log.Errorf("Failed to modify group %s: %s (%v)", g.Name, s, err)
		return fmt.Errorf("%s (%v)", s, err)
	}

	return web.JSONResponse("group modified", w)
}
