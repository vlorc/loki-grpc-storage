// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package wrapper

import (
	"context"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *serverStream) Context() context.Context {
	return ss.ctx
}

func WithLog(log *zap.Logger) grpc.UnaryServerInterceptor {
	return WithContext(func(ctx context.Context) context.Context {
		return utils.WithLog(ctx, log)
	})
}

func WithStreamLog(log *zap.Logger) grpc.StreamServerInterceptor {
	return WithStreamContext(func(ctx context.Context) context.Context {
		return utils.WithLog(ctx, log)
	})
}

func WithValue(key, value interface{}) grpc.UnaryServerInterceptor {
	return WithContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})
}

func WithStreamValue(key, value interface{}) grpc.StreamServerInterceptor {
	return WithStreamContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})
}

func WithContext(inject func(context.Context) context.Context) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(inject(ctx), req)
		return resp, err
	}
}

func WithStreamContext(inject func(context.Context) context.Context) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ss = &serverStream{
			ServerStream: ss,
			ctx:          inject(ss.Context()),
		}
		return handler(srv, ss)
	}
}
