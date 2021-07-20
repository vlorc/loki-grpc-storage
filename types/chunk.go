// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package types

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type ObjectClient interface {
	PutObject(ctx context.Context, key string, object []byte) error
	GetObject(ctx context.Context, key string) ([]byte, error)
	DeleteObject(ctx context.Context, key string) error
	Ping() error
}

func errInvalidChunkID(s string) error {
	return errors.Errorf("invalid chunk ID %q", s)
}

func ParseCheckId(key string) (*CheckInfo, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return nil, errInvalidChunkID(key)
	}
	userID := parts[0]
	hexParts := strings.Split(parts[1], ":")
	if len(hexParts) != 4 {
		return nil, errInvalidChunkID(key)
	}
	fingerprint, err := strconv.ParseUint(hexParts[0], 16, 64)
	if err != nil {
		return nil, err
	}
	from, err := strconv.ParseInt(hexParts[1], 16, 64)
	if err != nil {
		return nil, err
	}
	through, err := strconv.ParseInt(hexParts[2], 16, 64)
	if err != nil {
		return nil, err
	}
	checksum, err := strconv.ParseUint(hexParts[3], 16, 32)
	if err != nil {
		return nil, err
	}
	return &CheckInfo{
		Id:          key,
		UserID:      userID,
		Fingerprint: fingerprint,
		From:        time.Unix(0, from),
		Through:     time.Unix(0, through),
		Checksum:    uint32(checksum),
		ChecksumSet: true,
	}, nil
}

type CheckInfo struct {
	Id          string
	UserID      string
	Fingerprint uint64
	From        time.Time
	Through     time.Time
	Checksum    uint32
	ChecksumSet bool
}
