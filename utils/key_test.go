// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"bytes"
	"strings"
	"testing"
	"unsafe"
)

var __id1 = "fake/a70ecbaeaa65a26a:17ab9b3875f:17ab9b3889b:d8c9fe60"
var __id2 = "fake/a70ecbaeaa65a26a_17ab9b3875f_17ab9b3889b_d8c9fe60"

func BenchmarkAppend(b *testing.B) {
	var v [64]byte
	for i := 0; i < b.N; i++ {
		if s := AppendKey(__id1, v[:]); s != __id2 {
			b.Error("compare failed")
		}
	}
}

func BenchmarkReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if s := FormatKey(__id1); s != __id2 {
			b.Error("compare failed")
		}
	}
}

func BenchmarkBytesReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if s := string(bytes.ReplaceAll([]byte(__id1), []byte(":"), []byte("_"))); s != __id2 {
			b.Error("compare failed")
		}
	}
}

func BenchmarkStringsReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if s := strings.Replace(__id1, ":", "_", -1); s != __id2 {
			b.Error("compare failed")
		}
	}
}

func BenchmarkUnsafeBytesReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		k := bytes.ReplaceAll(*(*[]byte)(unsafe.Pointer(&__id1)), []byte(":"), []byte("_"))
		if s := *(*string)(unsafe.Pointer(&k)); s != __id2 {
			b.Error("compare failed")
		}
	}
}
