// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package qiniu

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Qiniu struct {
	log      *zap.Logger
	mac      *auth.Credentials
	manager  *storage.BucketManager
	uploader *storage.FormUploader
	bucket   string
	domain   string
	token    string
	client   *http.Client
	mkUrl    func(*auth.Credentials, string, string) string
}

var _ types.ObjectClient = &Qiniu{}

func New(log *zap.Logger, config *types.StoreConfig) types.ObjectClient {
	qn, err := Factory(log, config)
	if nil != err {
		panic(err)
	}
	return qn
}

func Factory(log *zap.Logger, config *types.StoreConfig) (types.ObjectClient, error) {
	qn := &Qiniu{
		log:    log,
		mac:    auth.New(config.Access, config.Secret),
		bucket: config.Bucket,
		domain: config.Url,
		client: http.DefaultClient,
		mkUrl: func(_ *auth.Credentials, domain string, key string) string {
			return storage.MakePublicURLv2(domain, key)
		},
	}

	cfg := &storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	if "" != config.Region {
		if zone, ok := storage.GetRegionByID(storage.RegionID(config.Region)); ok {
			cfg.Zone = &zone
		} else {
			log.Warn("can not found region", zap.String("region", config.Region))
		}
	}
	if strings.Index(config.Flag, "https") >= 0 {
		cfg.UseHTTPS = true
	}
	if strings.Index(config.Flag, "cdn") >= 0 {
		cfg.UseCdnDomains = true
	}
	if strings.Index(config.Flag, "private") >= 0 {
		qn.mkUrl = func(mac *auth.Credentials, domain string, key string) string {
			return storage.MakePrivateURLv2(mac, domain, key, time.Now().Add(time.Hour).Unix())
		}
	}

	qn.uploader = storage.NewFormUploader(cfg)
	qn.manager = storage.NewBucketManager(qn.mac, cfg)

	return qn, nil
}

func (qn *Qiniu) PutObject(ctx context.Context, key string, object []byte) error {
	return qn.write(ctx, key, object)
}

func (qn *Qiniu) GetObject(ctx context.Context, key string) ([]byte, error) {
	return qn.read(ctx, key)
}

func (qn *Qiniu) DeleteObject(ctx context.Context, key string) error {
	return qn.remove(ctx, key)
}

func (qn *Qiniu) Ping() error {
	policy := storage.PutPolicy{
		Scope:   qn.bucket,
		Expires: 7200,
	}
	qn.token = policy.UploadToken(qn.mac)

	return nil
}

func (qn *Qiniu) remove(ctx context.Context, key string) error {
	host, err := qn.manager.RsReqHost(qn.bucket)
	if err != nil {
		return err
	}

	rawurl := strings.Join([]string{host, storage.URIDelete(qn.bucket, key)}, "")

	return qn.manager.Client.CredentialedCall(ctx, qn.mac, auth.TokenQiniu, nil, "POST", rawurl, nil)
}

func (qn *Qiniu) write(ctx context.Context, key string, buf []byte) error {
	token := qn.token
	if "" == token {
		token = (&storage.PutPolicy{Scope: qn.bucket}).UploadToken(qn.mac)
	}

	return qn.uploader.Put(ctx, nil, token, key, bytes.NewReader(buf), int64(len(buf)), nil)
}

func (qn *Qiniu) read(ctx context.Context, key string) ([]byte, error) {
	rawurl := qn.mkUrl(qn.mac, qn.domain, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
	if nil != err {
		return nil, err
	}
	req.Header.Set("User-Agent", types.UserAgent)

	qn.log.Debug("request waiting", zap.String("path", key), zap.String("url", rawurl))

	resp, err := qn.client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("http status %d", resp.StatusCode)
		return nil, err
	}

	return utils.ReadAll(resp.Body)
}
