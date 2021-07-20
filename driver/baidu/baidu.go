// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package baidu

import (
	"context"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
)

type Baidu struct {
	log    *zap.Logger
	client *bos.Client
	bucket string
	domain string
}

var _ types.ObjectClient = &Baidu{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	qn, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return qn
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	endpoint := "bj.bcebos.com"
	if "" != config.Region {
		endpoint = config.Region + ".bcebos.com"
	}

	cli, err := bos.NewClient(config.Access, config.Secret, endpoint)
	if nil != err {
		return nil, err
	}
	bd := &Baidu{
		log:    log,
		bucket: config.Bucket,
		domain: config.Url,
		client: cli,
	}

	return bd, bd.Ping()
}

func (bd *Baidu) PutObject(ctx context.Context, key string, object []byte) error {
	return bd.write(ctx, key, object)
}

func (bd *Baidu) GetObject(ctx context.Context, key string) ([]byte, error) {
	return bd.read(ctx, key)
}

func (bd *Baidu) DeleteObject(ctx context.Context, key string) error {
	return bd.remove(ctx, key)
}

func (bd *Baidu) Ping() error {
	return nil
}

func (bd *Baidu) remove(ctx context.Context, key string) error {
	err := bd.client.DeleteObject(bd.bucket, key)

	return err
}

func (bd *Baidu) write(ctx context.Context, key string, buf []byte) error {
	_, err := bd.client.PutObjectFromBytes(bd.bucket, key, buf, nil)

	return err
}

func (bd *Baidu) read(ctx context.Context, key string) ([]byte, error) {
	resp, err := bd.client.GetObject(bd.bucket, key, nil)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	return utils.ReadAll(resp.Body)
}
