// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package driver

import (
	"github.com/pkg/errors"
	"github/vlorc/loki-grpc-storage/driver/aliyun"
	"github/vlorc/loki-grpc-storage/driver/baidu"
	"github/vlorc/loki-grpc-storage/driver/filesystem"
	"github/vlorc/loki-grpc-storage/driver/http"
	"github/vlorc/loki-grpc-storage/driver/memory"
	"github/vlorc/loki-grpc-storage/driver/qiniu"
	"github/vlorc/loki-grpc-storage/types"
	"go.uber.org/zap"
)

var driver = map[string]func(*zap.Logger, *types.StoreConfig) (types.ObjectClient, error){
	"fs":     filesystem.Factory,
	"qiniu":  qiniu.Factory,
	"baidu":  baidu.Factory,
	"aliyun": aliyun.Factory,
	"memory": memory.Factory,
	"http":   http.Factory,
	"empty": func(*zap.Logger, *types.StoreConfig) (types.ObjectClient, error) {
		return empty{}, nil
	},
}

func Register(name string, factory func(*zap.Logger, *types.StoreConfig) (types.ObjectClient, error)) {
	driver[name] = factory
}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	d, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return d
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	if factory, ok := driver[config.Driver]; ok {
		return factory(log.With(zap.String("driver", config.Driver), zap.String("name", config.Name)), config)
	}
	return nil, errors.Errorf("can not support driver '%s'", config.Driver)
}
