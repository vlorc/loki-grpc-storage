// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package wrapper

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/grpclog"
)

func setLoggerV2(log *zap.Logger) {
	grpclog.SetLoggerV2(&zapLog{
		log:   log.Sugar(),
		level: 0,
	})
}

type zapLog struct {
	log   *zap.SugaredLogger
	level int
}

func (l *zapLog) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *zapLog) Infoln(args ...interface{}) {
	l.log.Info(args...)
}

func (l *zapLog) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *zapLog) Warning(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *zapLog) Warningln(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *zapLog) Warningf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *zapLog) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *zapLog) Errorln(args ...interface{}) {
	l.log.Error(args...)
}

func (l *zapLog) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *zapLog) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *zapLog) Fatalln(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *zapLog) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *zapLog) V(level int) bool {
	return l.level <= level
}
