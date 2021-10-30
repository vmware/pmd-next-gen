// SPDX-License-Identifier: Apache-2.0

package system

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pm-web/pkg/conf"
)

func PathExists(path string) bool {
	_, r := os.Stat(path)
	return !os.IsNotExist(r)
}

func ReadFullFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFullFile(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	w.Flush()

	return nil
}

func ReadOneLineFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	line := scanner.Text()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return line, nil
}

func WriteOneLineFile(path string, line string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)

	return w.Flush()
}

func CreateDirectory(directoryPath string, perm os.FileMode) error {
	if !PathExists(directoryPath) {
		if err := os.Mkdir(directoryPath, perm); err != nil {
			return err
		}
	}

	return nil
}

func CreateDirectoryNested(directoryPath string, perm os.FileMode) error {
	if !PathExists(directoryPath) {
		if err := os.MkdirAll(directoryPath, perm); err != nil {
			return err
		}
	}

	return nil
}

func TLSFilePathExits() bool {
	return PathExists(path.Join(conf.ConfPath, conf.TLSCert)) && PathExists(path.Join(conf.ConfPath, conf.TLSKey))
}
