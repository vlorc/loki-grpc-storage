// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"io"
	"io/ioutil"
)

func WriteFile(p string, b []byte) error {
	return ioutil.WriteFile(p, b, 0644)
}

func ReadFile(p string) ([]byte, error) {
	return ioutil.ReadFile(p)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func ReadNop(r io.Reader) ([]byte, error) {
	return nil, nil
}
