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
	"sync"
	"time"
)

type StoreService struct {
	api.UnimplementedGrpcStoreServer
	store    types.ObjectClient
	level    zapcore.Level
	log      *zap.Logger
	parallel int
	min      int
}

type chunkResult struct {
	key   string
	data  []byte
	err   error
	begin time.Time
	end   time.Time
}

func NewStoreService(log *zap.Logger, conf *types.ChunkConfig, store types.ObjectClient) api.GrpcStoreServer {
	s := &StoreService{
		store:    store,
		log:      log,
		level:    types.Level(conf.Level),
		parallel: conf.Parallel,
		min:      conf.Min,
	}

	go s.ping()

	return s
}

func (s *StoreService) PutChunks(ctx context.Context, req *api.PutChunksRequest) (*empty.Empty, error) {
	if nil == s.store {
		return &empty.Empty{}, status.Errorf(codes.Unimplemented, "method PutChunks not implemented")
	}

	var err error
	if chunks := req.GetChunks(); s.parallel > s.min && len(chunks) > s.min {
		err = s.putChunksParallel(ctx, chunks)
	} else {
		err = s.putChunks(ctx, chunks)
	}

	return &empty.Empty{}, err
}

func (s *StoreService) GetChunks(req *api.GetChunksRequest, srv api.GrpcStore_GetChunksServer) (err error) {
	if nil == s.store {
		err = status.Errorf(codes.Unimplemented, "method GetChunks not implemented")
		return
	}

	var count int
	chunks := req.GetChunks()
	if s.parallel > s.min && len(chunks) > s.min {
		count, err = s.getChunksParallel(srv, chunks)
	} else {
		count, err = s.getChunks(srv, chunks)
	}

	s.print("getChunks", err, count, chunks)

	return err
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

func (s *StoreService) getChunks(srv api.GrpcStore_GetChunksServer, chunks []*api.Chunk) (int, error) {
	var last error
	var cache [64]byte

	ctx := srv.Context()
	log := utils.Log(ctx, s.log)
	count := 0

	for _, c := range chunks {
		r := &chunkResult{key: c.GetKey(), begin: time.Now()}
		r.data, r.err = s.store.GetObject(ctx, utils.AppendKey(r.key, cache[:]))
		r.end = time.Now()
		if err := s.sendChunk(log, srv, r); nil != err {
			last = err
		} else {
			count++
		}
	}

	return count, last
}

func (s *StoreService) getChunksParallel(srv api.GrpcStore_GetChunksServer, chunks []*api.Chunk) (int, error) {
	w := make(chan string)
	q := make(chan *chunkResult, s.parallel)
	g := &sync.WaitGroup{}

	g.Add(s.parallel)
	go __wait(chunks, w, q, g)

	ctx := srv.Context()
	for i := 0; i < s.parallel; i++ {
		go s.getChunkWork(ctx, w, q, g)
	}

	var last error
	log := utils.Log(ctx, s.log)
	count := 0

	for r := range q {
		if err := s.sendChunk(log, srv, r); nil != err {
			last = err
		} else {
			count++
		}
	}

	return count, last
}

func (s *StoreService) getChunkWork(ctx context.Context, queue chan string, result chan *chunkResult, group *sync.WaitGroup) {
	defer group.Done()

	var cache [64]byte

	for key := range queue {
		r := &chunkResult{
			key:   key,
			begin: time.Now(),
		}
		r.data, r.err = s.store.GetObject(ctx, utils.AppendKey(key, cache[:]))
		r.end = time.Now()
		result <- r
	}
}

func (s *StoreService) sendChunk(log *zap.Logger, srv api.GrpcStore_GetChunksServer, r *chunkResult) error {
	if nil != r.err {
		log.Error("getObject", zap.String("key", r.key), zap.Int("length", len(r.data)), zap.Duration("latency", r.end.Sub(r.begin)), zap.Error(r.err))
		return r.err
	}

	log.Debug("getObject", zap.String("key", r.key), zap.Int("length", len(r.data)), zap.Duration("latency", r.end.Sub(r.begin)))

	if nil != srv {
		now := time.Now()
		if err := srv.Send(&api.GetChunksResponse{Chunks: []*api.Chunk{{Key: r.key, Encoded: r.data}}}); nil != err {
			log.Error("sendObject", zap.String("key", r.key), zap.Int("length", len(r.data)), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
			return err
		}
		log.Debug("sendObject", zap.String("key", r.key), zap.Int("length", len(r.data)), zap.Duration("latency", time.Now().Sub(now)))
	}

	return nil
}

func __wait(chunks []*api.Chunk, work chan string, queue chan *chunkResult, group *sync.WaitGroup) {
	for _, c := range chunks {
		work <- c.GetKey()
	}
	close(work)
	group.Wait()
	close(queue)
}

func (s *StoreService) putChunks(ctx context.Context, chunks []*api.Chunk) error {
	var last error

	log := utils.Log(ctx, s.log)
	count := 0

	for _, c := range chunks {
		if err := s.putChunk(ctx, log, c); nil != err {
			last = err
		} else {
			count++
		}
	}

	s.print("putChunks", last, count, chunks)

	return last
}

func (s *StoreService) putChunk(ctx context.Context, log *zap.Logger, chunk *api.Chunk) error {
	var cache [64]byte

	key := chunk.GetKey()
	buf := chunk.GetEncoded()
	now := time.Now()
	err := s.store.PutObject(ctx, utils.AppendKey(key, cache[:]), buf)
	if nil != err {
		log.Error("putObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)), zap.Error(err))
	} else {
		log.Debug("putObject", zap.String("key", key), zap.Int("length", len(buf)), zap.Duration("latency", time.Now().Sub(now)))
	}

	return err
}

func (s *StoreService) putChunkWork(ctx context.Context, chunks chan *api.Chunk, group *sync.WaitGroup) {
	defer group.Done()

	log := utils.Log(ctx, s.log)

	for c := range chunks {
		_ = s.putChunk(ctx, log, c)
	}
}

func (s *StoreService) putChunksParallel(ctx context.Context, chunks []*api.Chunk) error {
	w := make(chan *api.Chunk)
	g := &sync.WaitGroup{}

	g.Add(s.parallel)
	go func() {
		for _, c := range chunks {
			w <- c
		}
		close(w)
	}()

	for i := 0; i < s.parallel; i++ {
		go s.putChunkWork(ctx, w, g)
	}

	return nil
}
