// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package server

import (
	"github/vlorc/loki-grpc-storage/api"
	"github/vlorc/loki-grpc-storage/driver"
	"github/vlorc/loki-grpc-storage/service"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/wrapper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	log    *zap.Logger
	config *types.Config
	server *grpc.Server
}

func NewServer(config *types.Config) *Server {
	return &Server{
		log:    types.NewLog(&config.Log),
		config: config,
	}
}

func (s *Server) Serve() error {
	log := types.NewLog(&s.config.Log)

	l, err := net.Listen("tcp", net.JoinHostPort(s.config.Server.Host, s.config.Server.Port))
	if err != nil {
		log.Error("listen failed ", zap.Error(err))
		return err
	}

	ss := grpc.NewServer(wrapper.Default(log)...)
	s.server = ss
	s.register(ss)

	log.Info("server listening at", zap.String("addr", l.Addr().String()))
	if err := ss.Serve(l); err != nil {
		log.Error("serve failed : %v", zap.Error(err))
	}

	return err
}

func (s *Server) Stop() {
	if nil != s.server {
		s.log.Info("server is being stopped")
		s.server.Stop()
	}
}

func (s *Server) register(ss *grpc.Server) {
	api.RegisterGrpcStoreServer(ss, service.NewStoreService(s.log, types.Level(s.config.Chunk.Level), driver.New(s.log, &s.config.Store)))
}
