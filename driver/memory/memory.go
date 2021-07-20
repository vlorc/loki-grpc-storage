// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package memory

import (
	"context"
	"github/vlorc/loki-grpc-storage/types"
	"go.uber.org/zap"
	"sync"
)

type Memory struct {
	mtx     sync.RWMutex
	objects map[string][]byte
	log     *zap.Logger
}

var _ types.ObjectClient = &Memory{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	fs, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return fs
}

func Factory(log *zap.Logger, _ *types.StoreConfig) (types.ObjectClient, error) {
	return &Memory{objects: map[string][]byte{}, log: log}, nil
}

func (mm *Memory) PutObject(ctx context.Context, key string, object []byte) error {
	mm.mtx.RLock()
	defer mm.mtx.RUnlock()

	mm.objects[key] = object
	return nil
}

func (mm *Memory) GetObject(ctx context.Context, key string) ([]byte, error) {
	mm.mtx.RLock()
	defer mm.mtx.RUnlock()

	buf := mm.objects[key]
	return buf, nil
}

func (mm *Memory) DeleteObject(ctx context.Context, key string) error {
	mm.mtx.Lock()
	defer mm.mtx.Unlock()

	delete(mm.objects, key)
	return nil
}

func (mm *Memory) Ping() error {
	return nil
}
