// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package wrapper

import (
	"context"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"time"
)

func WithCall(level zapcore.Level) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := utils.Log(ctx)
		start := time.Now()
		resp, err = handler(ctx, req)

		if l := log.Check(level, "server call"); nil != l {
			l.Write(
				zap.String("method", info.FullMethod),
				zap.Duration("latency", time.Now().Sub(start)),
				zap.Error(err),
			)
		}

		return resp, err
	}
}

func WithSteamCall(level zapcore.Level) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := utils.Log(ss.Context())
		start := time.Now()
		err = handler(srv, ss)

		if l := log.Check(level, "server call"); nil != l {
			l.Write(
				zap.String("method", info.FullMethod),
				zap.Duration("latency", time.Now().Sub(start)),
				zap.Error(err),
			)
		}
		return err
	}
}
