// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package service

import (
	"context"
	"github.com/baidubce/bce-sdk-go/util/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github/vlorc/loki-grpc-storage/api"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type StoreService struct {
	api.UnimplementedGrpcStoreServer
	store types.ObjectClient
	level zapcore.Level
	log   *zap.Logger
}

func NewStoreService(log *zap.Logger, level zapcore.Level, store types.ObjectClient) api.GrpcStoreServer {
	s := &StoreService{
		store: store,
		log:   log,
		level: level,
	}

	go s.ping()

	return s
}

func (s *StoreService) PutChunks(ctx context.Context, req *api.PutChunksRequest) (*empty.Empty, error) {
	if nil == s.store {
		return &empty.Empty{}, status.Errorf(codes.Unimplemented, "method PutChunks not implemented")
	}

	var last error
	var cache [64]byte

	log := utils.Log(ctx, s.log)
	chunks := req.GetChunks()
	count := 0

	for _, c := range chunks {

		key := c.GetKey()
		buf := c.GetEncoded()
		now := time.Now()

		if err := s.store.PutObject(ctx, utils.AppendKey(key, cache[:]), buf); nil != err {
			last = err
			log.Error("putObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
		} else {
			count++
			log.Debug("putObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)))
		}
	}

	s.print("putChunks", last, count, chunks)

	return &empty.Empty{}, last
}

func (s *StoreService) GetChunks(req *api.GetChunksRequest, srv api.GrpcStore_GetChunksServer) error {
	if nil == s.store {
		return status.Errorf(codes.Unimplemented, "method GetChunks not implemented")
	}

	var last error
	var cache [64]byte

	log := utils.Log(srv.Context(), s.log)
	chunks := req.GetChunks()
	ctx := srv.Context()
	count := 0

	for _, c := range chunks {

		key := c.GetKey()
		now := time.Now()

		if buf, err := s.store.GetObject(ctx, utils.AppendKey(key, cache[:])); nil != err {
			last = err
			log.Error("getObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
		} else {
			log.Debug("getObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)))
			if err = srv.Send(&api.GetChunksResponse{Chunks: []*api.Chunk{{Key: key, Encoded: buf}}}); nil != err {
				last = err
				log.Error("sendObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
			} else {
				count++
				log.Debug("sendObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)))
			}
		}
	}

	s.print("getChunks", last, count, chunks)

	return last
}

func (s *StoreService) DeleteChunks(ctx context.Context, req *api.ChunkID) (*empty.Empty, error) {
	if nil == s.store {
		return &empty.Empty{}, status.Errorf(codes.Unimplemented, "method DeleteChunks not implemented")
	}

	var cache [64]byte

	log := utils.Log(ctx, s.log)
	key := req.GetChunkID()
	now := time.Now()

	err := s.store.DeleteObject(ctx, utils.AppendKey(key, cache[:]))
	if nil != err {
		log.Error("delObject", zap.String("key", key), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
	} else {
		log.Debug("delObject", zap.String("key", key), zap.Duration("latency", time.Now().Sub(now)))
	}

	s.printId("deleteChunks", err, key)

	return &empty.Empty{}, err
}

func (s *StoreService) ping() {
	for range time.NewTicker(time.Hour).C {
		if err := s.store.Ping(); nil != err {
			log.Error("driver ping", zap.Error(err))
		} else {
			log.Debug("driver ping")
		}
	}
}

func (s *StoreService) printId(msg string, err error, id ...string) {
	if l := s.log.Check(s.level, msg); nil != l {
		l.Write(zap.Strings("key", id), zap.Error(err))
	}
}

func (s *StoreService) print(msg string, err error, count int, chunks []*api.Chunk) {
	l := s.log.Check(s.level, msg)
	if nil == l {
		return
	}

	var hit [8]string
	var keys []string

	if len(chunks) > len(hit) {
		keys = make([]string, len(chunks))
	} else {
		keys = hit[:len(chunks)]
	}
	for i := range keys {
		keys[i] = chunks[i].Key
	}

	l.Write(zap.Strings("keys", keys), zap.Int("count", count), zap.Int("total", len(keys)), zap.Error(err))
}
