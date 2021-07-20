// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package http

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type HTTP struct {
	log    *zap.Logger
	url    string
	client *http.Client
}

var _ types.ObjectClient = &HTTP{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	fs, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return fs
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	h := &HTTP{
		url:    config.Url,
		log:    log,
		client: http.DefaultClient,
	}
	if !strings.HasSuffix(h.url, "/") {
		h.url += "/"
	}
	return h, h.Ping()
}

func (h *HTTP) PutObject(ctx context.Context, key string, object []byte) error {
	return h.write(ctx, key, object)
}

func (h *HTTP) GetObject(ctx context.Context, key string) ([]byte, error) {
	return h.read(ctx, key)
}

func (h *HTTP) DeleteObject(ctx context.Context, key string) error {
	return h.remove(ctx, key)
}

func (h *HTTP) Ping() error {
	return nil
}

func (h *HTTP) remove(ctx context.Context, key string) error {
	_, err := h.request(ctx, http.MethodDelete, key, nil, utils.ReadNop)

	return err
}

func (h *HTTP) write(ctx context.Context, key string, buf []byte) error {
	_, err := h.request(ctx, http.MethodPost, key, bytes.NewReader(buf), utils.ReadNop)

	return err
}

func (h *HTTP) read(ctx context.Context, key string) ([]byte, error) {
	return h.request(ctx, http.MethodGet, key, nil, utils.ReadAll)
}

func (h *HTTP) request(ctx context.Context, method string, key string, body io.Reader, read func(io.Reader) ([]byte, error)) ([]byte, error) {
	rawurl := h.url + strings.ReplaceAll(key, ":", "_")

	req, err := http.NewRequestWithContext(ctx, method, rawurl, body)
	if nil != err {
		return nil, err
	}
	req.Header.Set("User-Agent", types.UserAgent)

	h.log.Debug("request waiting", zap.String("path", key), zap.String("url", rawurl), zap.String("method", method))

	resp, err := h.client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("http status %d", resp.StatusCode)
		return nil, err
	}
	return read(resp.Body)
}
