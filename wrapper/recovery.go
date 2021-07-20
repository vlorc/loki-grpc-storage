// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package wrapper

import (
	"context"
	"errors"
	"fmt"
	"github/vlorc/loki-grpc-storage/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithRecovery(level zapcore.Level) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := utils.Log(ctx)
		defer func() {
			if r := recover(); r != nil {
				err = __recovery(log, level, info.FullMethod, r)
			}
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func WithSteamRecovery(level zapcore.Level) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := utils.Log(ss.Context())
		defer func() {
			if r := recover(); r != nil {
				err = __recovery(log, level, info.FullMethod, r)
			}
		}()

		return handler(srv, ss)
	}
}

func __recovery(log *zap.Logger, level zapcore.Level, method string, r interface{}) error {
	if ce := log.Check(level, "call crash"); ce != nil {
		stack := utils.Stack(3)
		ce.Write(
			zap.String("method", method),
			zap.Error(errors.New(fmt.Sprint(r))),
			zap.String("stack", string(stack)),
		)
	}
	return status.Errorf(codes.Internal, "%v", r)
}
