// SPDX-License-Identifier: Apache-2.0

package share

import (
	"os/exec"
)

// CheckBinaryExists verifies if binary exists
func CheckBinaryExists(binary string) error {
	_, err := exec.LookPath(binary)
	if err != nil {
		return err
	}

	return nil
}
