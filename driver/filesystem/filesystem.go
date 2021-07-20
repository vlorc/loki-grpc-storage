// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package filesystem

import (
	"context"
	"github.com/pkg/errors"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

type FS struct {
	Directory string
	log       *zap.Logger
}

var _ types.ObjectClient = &FS{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	fs, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return fs
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	root, err := filepath.Abs(filepath.Clean(config.Url))
	if nil != err {
		return nil, err
	}
	fs := &FS{Directory: root, log: log}

	return fs, fs.Ping()
}

func (fs *FS) PutObject(ctx context.Context, key string, object []byte) error {
	return fs.write(ctx, key, object)
}

func (fs *FS) GetObject(ctx context.Context, key string) ([]byte, error) {
	return fs.read(ctx, key)
}

func (fs *FS) DeleteObject(ctx context.Context, key string) error {
	return fs.remove(ctx, key)
}

func (fs *FS) Ping() error {
	return fs.ping()
}

func (fs *FS) ping() error {
	stat, err := os.Stat(fs.Directory)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(fs.Directory, 0755)
		}
	} else if !stat.IsDir() {
		err = errors.Errorf("the path must be a directory '%s'", err)
	}

	return err
}

func (fs *FS) remove(ctx context.Context, key string) error {
	p, err := realpath(fs.Directory, key)
	if nil != err {
		return err
	}

	stat, err := os.Stat(p)
	if nil != err {
		return err
	}
	if stat.IsDir() {
		err = os.RemoveAll(p)
	} else {
		err = os.Remove(p)
	}

	return err
}

func (fs *FS) write(ctx context.Context, key string, buf []byte) error {
	p, err := realpath(fs.Directory, key)
	if nil != err {
		return err
	}

	if err = utils.WriteFile(p, buf); nil != err {
		if os.IsNotExist(err) && fs.mkdir(p) {
			err = utils.WriteFile(p, buf)
		}
	}

	return err
}

func (fs *FS) read(ctx context.Context, key string) ([]byte, error) {
	p, err := realpath(fs.Directory, key)
	if nil != err {
		return nil, err
	}

	return utils.ReadFile(p)
}

func (fs *FS) mkdir(p string) bool {
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); nil != err {
		fs.log.Error("make directory", zap.String("path", dir), zap.Error(err))
		return false
	}
	fs.log.Debug("make directory", zap.String("path", dir))

	return true
}

func realpath(root, name string) (string, error) {
	p := filepath.Join(root, filepath.FromSlash(name))
	p, err := filepath.Abs(p)
	if nil != err {
		return "", err
	}
	if !strings.HasPrefix(p, root) {
		return "", errors.Errorf("invalid path %s to %s", name, p)
	}

	return p, nil
}
