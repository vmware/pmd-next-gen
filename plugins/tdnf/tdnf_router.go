// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package tdnf

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"

	"github.com/pmd-nextgen/pkg/validator"
	"github.com/pmd-nextgen/pkg/web"
)

func routerParseOptions(values map[string][]string) Options {

	isTrue := func(key string) bool {
		if v, ok := values[key]; ok {
			return validator.IsBool(v[0])
		}
		return false
	}

	getString := func(key string) string {
		if v, ok := values[key]; ok {
			return v[0]
		}
		return ""
	}

	var options Options

	v := reflect.ValueOf(&options).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		value := v.Field(i).Interface()
		switch value.(type) {
		case bool:
			v.Field(i).SetBool(isTrue(name))
		case string:
			v.Field(i).SetString(getString(name))
		case []string:
			size := len(values[name])
			if size > 0 {
				v.Field(i).Set(reflect.MakeSlice(reflect.TypeOf([]string{}), size, size))
				for j, s := range values[name] {
					v.Field(i).Index(j).Set(reflect.ValueOf(s))
				}
			}
		}
	}
	return options
}

func routeracquireCommand(w http.ResponseWriter, r *http.Request) {
	var err error

	if err = r.ParseForm(); err != nil {
		web.JSONResponseError(err, w)
	}
	options := routerParseOptions(r.Form)

	switch cmd := mux.Vars(r)["command"]; cmd {
	case "check-update":
		err = acquireCheckUpdate(w, "", options)
	case "clean":
		err = acquireClean(w, options)
	case "distro-sync":
		err = acquireAlterCmd(w, cmd, "", options)
	case "downgrade":
		err = acquireAlterCmd(w, cmd, "", options)
	case "info":
		err = acquireInfoList(w, "", options)
	case "list":
		err = acquireList(w, "", options)
	case "makecache":
		err = acquireMakeCache(w, options)
	case "repolist":
		err = acquireRepoList(w, options)
	case "search":
		q := r.FormValue("q")
		if q != "" {
			err = acquireSearch(w, q, options)
		} else {
			err = errors.New("search needs 'q=str' query")
		}
	case "update":
		err = acquireAlterCmd(w, cmd, "", options)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func routeracquireCommandPkg(w http.ResponseWriter, r *http.Request) {
	var err error

	pkg := mux.Vars(r)["pkg"]

	if err = r.ParseForm(); err != nil {
		web.JSONResponseError(err, w)
	}
	options := routerParseOptions(r.Form)

	switch cmd := mux.Vars(r)["command"]; cmd {
	case "autoremove":
		err = acquireAlterCmd(w, cmd, pkg, options)
	case "downgrade":
		err = acquireAlterCmd(w, cmd, pkg, options)
	case "check-update":
		err = acquireCheckUpdate(w, pkg, options)
	case "erase":
		err = acquireAlterCmd(w, cmd, pkg, options)
	case "info":
		err = acquireInfoList(w, pkg, options)
	case "install":
		err = acquireAlterCmd(w, cmd, pkg, options)
	case "list":
		err = acquireList(w, pkg, options)
	case "reinstall":
		err = acquireAlterCmd(w, cmd, pkg, options)
	case "update":
		err = acquireAlterCmd(w, cmd, pkg, options)
	default:
		err = errors.New("unsupported")
	}

	if err != nil {
		web.JSONResponseError(err, w)
	}
}

func RegisterRouterTdnf(router *mux.Router) {
	n := router.PathPrefix("/tdnf").Subrouter().StrictSlash(false)

	n.HandleFunc("/{command}/{pkg}", routeracquireCommandPkg).Methods("GET")
	n.HandleFunc("/{command}", routeracquireCommand).Methods("GET")
}
