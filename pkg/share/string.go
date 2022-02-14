// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package share

import (
	"errors"
	"strings"
)

func StringContains(list []string, s string) bool {
	set := make(map[string]int)

	for k, v := range list {
		set[v] = k
	}

	return set[s] != 0
}

func StringDeleteSlice(list []string, s string) ([]string, error) {
	set := make(map[string]int)

	for k, v := range list {
		set[v] = k
	}

	i, v := set[s]
	if v {
		list = append(list[:i], list[i+1:]...)
		return list, nil
	}

	return nil, errors.New("slice not found")
}

func UniqueString(s []string, t []string) []string {
	set := NewSet()

	list := []string{}
	for _, e := range s {
		if e == "" {
			continue
		}
		if v := set.Contains(e); !v {
			set.Add(e)
			list = append(list, strings.TrimSpace(e))
		}
	}

	for _, e := range t {
		if e == "" {
			continue
		}
		if v := set.Contains(e); !v {
			set.Add(e)
			list = append(list, strings.TrimSpace(e))
		}
	}

	return list
}
