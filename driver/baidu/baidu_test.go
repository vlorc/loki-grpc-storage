// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package baidu

import (
	"bytes"
	"context"
	"github/vlorc/loki-grpc-storage/types"
	"go.uber.org/zap"
	"testing"
)

var __id = "fake/a70ecbaeaa65a26a_17ab9b3875f_17ab9b3889b_d8c9fe60"

func __new() types.ObjectClient {
	log, _ := zap.NewDevelopment()
	return New(log, &types.StoreConfig{
		Driver: "baidu",
		Name:   "baidu",
		Access: "xxxxxxxxxxx",
		Secret: "xxxxxxxxxxx",
		Bucket: "test",
	})
}

func TestBaidu_Object(t *testing.T) {
	d := __new()

	src := []byte("ccccccccccccccccccccccccccccccccccccccccc")

	if err := d.PutObject(context.Background(), __id, src); nil != err {
		t.Error("putObject failed", err.Error())
	}
	dst, err := d.GetObject(context.Background(), __id)
	if nil != err {
		t.Error("getObject failed", err.Error())
	}
	if bytes.Compare(src, dst) != 0 {
		t.Error("compare failed")
	}
	if err := d.DeleteObject(context.Background(), __id); nil != err {
		t.Error("delObject", err.Error())
	}
}
