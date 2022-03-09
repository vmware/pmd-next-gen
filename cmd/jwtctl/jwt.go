// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 VMware, Inc.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

func readFromStdIn() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return []byte{}, fmt.Errorf("could not get StdIn")
	}
	if info.Mode()&os.ModeCharDevice == os.ModeCharDevice { // || info.Size() <= 0 {
		return []byte{}, fmt.Errorf("could not read the data")
	}

	reader := bufio.NewReader(os.Stdin)
	var output []byte

	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return output, nil
}

func encode() {
	var signAlgorithm *jwt.SigningMethodHMAC

	signMethod := "H256"
	switch signMethod {
	case "H256":
		signAlgorithm = jwt.SigningMethodHS256
	case "H384":
		signAlgorithm = jwt.SigningMethodHS384
	case "H512":
		signAlgorithm = jwt.SigningMethodHS512
	default:
		signAlgorithm = jwt.SigningMethodHS256
	}

	data := "{\"Hello\":\"world\"}"
	/*if data == "@-" {
		stdIn, err := readFromStdIn()
		if err != nil {
			return fmt.Errorf("could not read from stdIn: %w", err)
		}
		data = string(stdIn)
	}*/

	var dataJSON map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dataJSON); err != nil {
		fmt.Errorf("could not unmarshal the data: %w", err)
		return
	}

	if t, ok := dataJSON["exp"]; ok { // t is a unix timestamp
		dataJSON["exp"] = t
	} else {
		dataJSON["exp"] = time.Now().Add(5 * 24 * time.Hour).Unix()
	}

	if t, ok := dataJSON["iat"]; ok { // t is a unix timestamp
		dataJSON["iat"] = t
	} else {
		dataJSON["iat"] = time.Now().Unix()
	}

	claim := jwt.NewWithClaims(
		signAlgorithm, jwt.MapClaims(
			dataJSON,
		),
	)

	token, err := claim.SignedString([]byte("test123"))
	if err != nil {
		fmt.Printf("could not write token")
	} else {

		fmt.Printf("%s\n", token)
	}

}
