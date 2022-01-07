// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package sysctl

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pmd-nextgen/pkg/share"
	"github.com/pmd-nextgen/pkg/system"
	"github.com/pmd-nextgen/pkg/web"
	log "github.com/sirupsen/logrus"
)

const (
	sysctlDirPath = "/etc/sysctl.d"
	sysctlPath    = "/etc/sysctl.conf"
	procSysPath   = "/proc/sys/"
)

// Sysctl json request
type Sysctl struct {
	Key      string   `json:"Key"`
	Value    string   `json:"Value"`
	Apply    string   `json:"Apply"`
	Pattern  string   `json:"Pattern"`
	FileName string   `json:"FileName"`
	Files    []string `json:"Files"`
}

// Get filepath from key
func pathFromKey(key string) string {
	return filepath.Join(procSysPath, strings.Replace(key, ".", "/", -1))
}

// Get key from filepath
func keyFromPath(path string) string {
	subPath := strings.TrimPrefix(path, procSysPath)
	return strings.Replace(subPath, "/", ".", -1)
}

// Apply sysctl configuration to system
func (s *Sysctl) apply(fileName string) error {
	b, err := share.ParseBool(s.Apply)
	if err != nil || !b {
		return fmt.Errorf("Failed to apply sysctl: '%s'", fileName)
	}

	path, err := exec.LookPath("sysctl")
	if err != nil {
		return err
	}

	cmd := exec.Command(path, "-p", fileName)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Failed to load sysctl variable: %s", stdout)
		return fmt.Errorf("Failed to load sysctl variable: %s", stdout)
	}

	return err
}

// Read configuration file and prepare sysctlMap
func readSysctlConfigFromFile(path string, sysctlMap map[string]string) error {
	lines, err := system.ReadFullFile(path)
	if err != nil {
		return fmt.Errorf("Failed to %v", err)
	}

	for _, line := range lines {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			log.Errorf("Could not parse line : '%s'", line)
			continue
		}

		k := strings.TrimSpace(tokens[0])
		v := strings.TrimSpace(tokens[1])
		sysctlMap[k] = v
	}

	return err
}

// Write sysctlMap entry in configuration file
func writeSysctlConfigInFile(confFile string, sysctlMap map[string]string) error {
	var lines []string
	var line string

	for k, v := range sysctlMap {
		line = k + "=" + v
		lines = append(lines, line)
	}

	return system.WriteFullFile(confFile, lines)
}

// Read /etc/sysctl.conf file and prepare sysctlMap
func createSysctlMapFromConfFile(sysctlMap map[string]string) error {
	err := readSysctlConfigFromFile(sysctlPath, sysctlMap)
	return err
}

// Traverse the baseDirPath and prepare sysctlMap
func createSysctlMapFromDir(baseDirPath string, sysctlMap map[string]string) error {
	err := filepath.Walk(baseDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Error accessing sysctl path: %v", err)
		}
		if info.IsDir() {
			return err
		}

		// Reading all files from procSysPath so need to create key and insert in map.
		if baseDirPath == procSysPath {
			key := keyFromPath(path)
			val, err := os.ReadFile(path)
			if err != nil {
				log.Errorf("Failed to read file '%s': %v", path, err)
			} else {
				sysctlMap[key] = strings.TrimSpace(string(val))
			}
		} else {
			err = readSysctlConfigFromFile(path, sysctlMap)
			if err != nil {
				log.Errorf("Failed to read file '%s': %v", path, err)
			}
		}

		return nil
	})

	return err
}

// Fetch key value from proc/sys directory
func getKeyValueFromProcSys(key string, sysctlMap map[string]string) error {
	data, err := os.ReadFile(pathFromKey(key))
	if err != nil {
		return err
	}

	sysctlMap[key] = strings.TrimSpace(string(data))
	return err
}

// Get sysctl key value from any of the following
// sysctl.conf, sysctl.d or proc/sys
func (s *Sysctl) Get(rw http.ResponseWriter) error {
	sysctlMap := make(map[string]string)
	if len(s.Key) == 0 {
		return fmt.Errorf("Failed to get sysctl parameter. Input key missing")
	}

	// First try to get from sysctl.conf.
	err := createSysctlMapFromConfFile(sysctlMap)
	_, ok := sysctlMap[s.Key]
	if ok {
		return web.JSONResponse(sysctlMap[s.Key], rw)
	}

	// Cant find the key from main sysctl.conf read sysctl.d dir files.
	err = createSysctlMapFromDir(sysctlDirPath, sysctlMap)
	_, ok = sysctlMap[s.Key]
	if ok {
		return web.JSONResponse(sysctlMap[s.Key], rw)
	}

	// Cant find the key from sysctl.d try to get from proc/sys.
	err = getKeyValueFromProcSys(s.Key, sysctlMap)
	_, ok = sysctlMap[s.Key]
	if ok {
		return web.JSONResponse(sysctlMap[s.Key], rw)
	}

	log.Errorf("Failed to get the sysctl key[%s] value from all config: %v", s.Key, err)
	return err
}

// GetPatern will return all the entry with matching pattern
// If pattern is empty it should return all values
func (s *Sysctl) GetPattern(rw http.ResponseWriter) error {
	sysctlMap := make(map[string]string)
	if len(s.Pattern) == 0 {
		log.Errorf("Input pattern is empty return all system configuration")
	}

	re, err := regexp.CompilePOSIX(s.Pattern)
	if err != nil {
		return fmt.Errorf("Failed to get sysctl parameter, Invalid pattern %s: %v", s.Pattern, err)
	}

	err = createSysctlMapFromConfFile(sysctlMap)
	if err != nil {
		log.Errorf("Failed reading configuration from %s: %v", sysctlPath, err)
	}

	err = createSysctlMapFromDir(sysctlDirPath, sysctlMap)
	if err != nil {
		log.Errorf("Failed reading configuration from %s: %v", sysctlDirPath, err)
	}

	err = createSysctlMapFromDir(procSysPath, sysctlMap)
	if err != nil {
		log.Errorf("Failed reading configuration from %s: %v", procSysPath, err)
		return err
	}

	result := make(map[string]string)
	for k, v := range sysctlMap {
		if !re.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return web.JSONResponse(result, rw)
}

// Update sysctl configuration file and apply
// Action can be SET, UPDATE or DELETE
func (s *Sysctl) Update() error {
	sysctlMap := make(map[string]string)
	if len(s.FileName) == 0 {
		s.FileName = sysctlPath
	} else {
		s.FileName = filepath.Join(sysctlDirPath, s.FileName)
	}

	if len(s.Key) == 0 {
		return fmt.Errorf("Input Key is missing in json data")
	}

	if len(s.Value) == 0 {
		return fmt.Errorf("Input Value is missing in json data")
	}

	err := readSysctlConfigFromFile(s.FileName, sysctlMap)
	if err != nil {
		return fmt.Errorf("could not parse file %s: %v", s.FileName, err)
	}

	if s.Value == "Delete" {
		_, ok := sysctlMap[s.Key]
		if !ok {
			return fmt.Errorf("Failed to delete sysctl parameter '%s'. Key not found", s.Key)
		}
		delete(sysctlMap, s.Key)
	} else {
		sysctlMap[s.Key] = s.Value
	}

	///Update config file and apply.
	err = writeSysctlConfigInFile(s.FileName, sysctlMap)
	if err != nil {
		return fmt.Errorf("Failed to update in file %s: %v", s.FileName, err)
	}

	return s.apply(s.FileName)
}

// Load all the configuration files and apply
func (s *Sysctl) Load() error {
	sysctlMap := make(map[string]string)
	if len(s.Files) == 0 {
		s.Files = []string{sysctlPath}
	}

	for _, f := range s.Files {
		if f != sysctlPath {
			f = filepath.Join(sysctlDirPath, f)
		}

		if err := readSysctlConfigFromFile(f, sysctlMap); err != nil {
			return fmt.Errorf("Failed to parse file %s: %v", f, err)
		}
	}

	err := writeSysctlConfigInFile(sysctlPath, sysctlMap)
	if err != nil {
		return err
	}

	return s.apply(sysctlPath)
}
