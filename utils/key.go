// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"bytes"
	"strings"
	"unsafe"
)

func FormatKey(k string) string {
	return AppendKey(k, nil)
}

func AppendKey(k string, b []byte) string {
	i := strings.IndexByte(k, ':')
	if i < 0 {
		return k
	}

	if len(k) > len(b) {
		b = make([]byte, len(k))
	} else {
		b = b[:len(k)]
	}

	copy(b, *(*[]byte)(unsafe.Pointer(&k)))
	for t := b; i >= 0; i = bytes.IndexByte(t, ':') {
		t[i], t = '_', t[i+1:]
	}
	s := *(*string)(unsafe.Pointer(&b))

	return s
}
