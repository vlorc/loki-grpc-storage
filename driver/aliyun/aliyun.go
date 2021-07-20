// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package aliyun

import (
	"bytes"
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"strings"
)

type Aliyun struct {
	log    *zap.Logger
	client *oss.Client
	bucket *oss.Bucket
	option []oss.Option
}

var _ types.ObjectClient = &Aliyun{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	qn, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return qn
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	var opts []oss.ClientOption

	if strings.Index(config.Flag, "cname") >= 0 {
		opts = append(opts, oss.UseCname(true))
	}
	if strings.Index(config.Flag, "md5") >= 0 {
		opts = append(opts, oss.EnableMD5(true))
	}
	if strings.Index(config.Flag, "crc") >= 0 {
		opts = append(opts, oss.EnableCRC(false))
	}
	if "" != config.Token {
		opts = append(opts, oss.SecurityToken(config.Token))
	}

	client, err := oss.New(config.Url, config.Access, config.Secret, opts...)
	if nil != err {
		return nil, err
	}

	bucket, err := client.Bucket(config.Bucket)
	if nil != err {
		return nil, err
	}

	al := &Aliyun{
		log:    log,
		bucket: bucket,
		client: client,
	}

	if strings.Index(config.Flag, "archive") >= 0 {
		al.option = append(al.option, oss.ObjectStorageClass(oss.StorageArchive))
	}
	if strings.Index(config.Flag, "cold") >= 0 {
		al.option = append(al.option, oss.ObjectStorageClass(oss.StorageColdArchive))
	}
	if strings.Index(config.Flag, "ia") >= 0 {
		al.option = append(al.option, oss.ObjectStorageClass(oss.StorageIA))
	}
	if strings.Index(config.Flag, "private") >= 0 {
		al.option = append(al.option, oss.ObjectACL(oss.ACLPrivate))
	}

	return al, al.Ping()
}

func (bd *Aliyun) PutObject(ctx context.Context, key string, object []byte) error {
	return bd.write(ctx, key, object)
}

func (al *Aliyun) GetObject(ctx context.Context, key string) ([]byte, error) {
	return al.read(ctx, key)
}

func (al *Aliyun) DeleteObject(ctx context.Context, key string) error {
	return al.remove(ctx, key)
}

func (al *Aliyun) Ping() error {
	return nil
}

func (al *Aliyun) remove(ctx context.Context, key string) error {
	err := al.bucket.DeleteObject(key)

	return err
}

func (al *Aliyun) write(ctx context.Context, key string, buf []byte) error {
	err := al.bucket.PutObject(key, bytes.NewReader(buf), al.option...)

	return err
}

func (al *Aliyun) read(ctx context.Context, key string) ([]byte, error) {
	body, err := al.bucket.GetObject(key)
	if nil != err {
		return nil, err
	}
	defer body.Close()

	return utils.ReadAll(body)
}
