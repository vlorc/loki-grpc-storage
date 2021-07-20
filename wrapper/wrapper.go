// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package wrapper

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Default(log *zap.Logger, opts ...grpc.ServerOption) []grpc.ServerOption {

	setLoggerV2(log.With(zap.String("driver", "grpc")))

	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			WithLog(log),
			WithRecovery(zap.ErrorLevel),
			WithCall(zap.InfoLevel),
		),
		grpc.ChainStreamInterceptor(
			WithStreamLog(log),
			WithSteamRecovery(zap.ErrorLevel),
			WithSteamCall(zap.InfoLevel),
		),
	)

	return opts
}
