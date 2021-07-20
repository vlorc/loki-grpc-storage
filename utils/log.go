// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logKey struct{}

func Log(ctx context.Context, log ...*zap.Logger) *zap.Logger {
	if l, ok := ctx.Value(logKey{}).(*zap.Logger); ok {
		return l
	}
	if len(log) > 0 {
		return log[0]
	}
	return zap.L()
}

func WithLog(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, log)
}

func Println(log *zap.Logger, level zapcore.Level, msg string, field ...zap.Field) {
	if l := log.Check(level, msg); nil != l {
		l.Write(field...)
	}
}
