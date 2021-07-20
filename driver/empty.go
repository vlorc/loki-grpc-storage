// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package driver

import (
	"context"
	"github/vlorc/loki-grpc-storage/types"
)

type empty struct{}

var _ types.ObjectClient = empty{}

func (e empty) PutObject(ctx context.Context, key string, object []byte) error {
	return nil
}

func (e empty) GetObject(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (e empty) DeleteObject(ctx context.Context, key string) error {
	return nil
}

func (e empty) Ping() error {
	return nil
}
