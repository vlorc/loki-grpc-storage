// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Flag(conf interface{}) {
	v := reflect.ValueOf(conf)
	if reflect.Ptr != v.Kind() {
		return
	}
	if v = v.Elem(); reflect.Struct == v.Kind() {
		__flag(v, "")
		flag.Parse()
	}
}

func __flag(val reflect.Value, parent string) {
	for t, i, n := val.Type(), 0, val.NumField(); i < n; i++ {
		f := t.Field(i)
		s := f.Tag.Get("flag")

		if "" == s || "-" == s {
			continue
		}
		if reflect.Struct == f.Type.Kind() {
			__flag(val.Field(i), parent+s+".")
			continue
		}

		tags := strings.Split(s, ",")
		name := parent + tags[0]
		usage := tags[0]

		if len(tags) >= 3 {
			usage = tags[2]
		}

		switch f.Type.Kind() {
		case reflect.String:
			v := ""
			if len(tags) >= 2 {
				v = __value(tags[1])
			}
			flag.StringVar(val.Field(i).Addr().Interface().(*string), name, v, usage)
		case reflect.Int:
			v := 0
			if len(tags) >= 2 {
				v, _ = strconv.Atoi(tags[1])
			}
			flag.IntVar(val.Field(i).Addr().Interface().(*int), name, v, usage)
		case reflect.Bool:
			flag.BoolVar(val.Field(i).Addr().Interface().(*bool), name, len(tags) >= 2 && "true" == tags[1], usage)
		}
	}
}

func __value(val string) string {
	if len(val) < 3 || '{' != val[0] || '}' != val[len(val)-1] {
		return val
	}

	val = val[1 : len(val)-1]
	switch val {
	case "hostname":
		val, _ = os.Hostname()
	case "tmpdir":
		val = os.TempDir()
	case "workdir":
		val, _ = os.Getwd()
	case "timestamp":
		val = strconv.FormatInt(time.Now().Unix(), 10)
	}

	return val
}
