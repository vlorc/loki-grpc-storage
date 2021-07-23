// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package types

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var __level = []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

func Level(v string) zapcore.Level {
	return __zapLevel(v)
}

func __zapLevel(v string) zapcore.Level {
	for i := range __level {
		if __level[i] == v {
			return zapcore.Level(i - 1)
		}
	}
	return zap.InfoLevel
}

func __zapConfig(config *LogConfig) zap.Config {
	if "dev" == config.Mode {
		return zap.NewDevelopmentConfig()
	}

	conf := zap.NewProductionConfig()
	if !config.Caller {
		conf.DisableCaller = true
	} else {
		conf.DisableCaller = false
	}
	if !config.Trace {
		conf.DisableStacktrace = true
	} else {
		conf.DisableStacktrace = false
	}
	conf.Level = zap.NewAtomicLevelAt(zapcore.Level(__zapLevel(config.Level)))

	return conf
}

func __zapLog(config *LogConfig) *zap.Logger {
	var opts []zap.Option
	sink := zapcore.Lock(os.Stdout)

	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	if config.Caller {
		opts = append(opts, zap.AddCaller())
	}
	if config.Trace {
		opts = append(opts, zap.AddStacktrace(zapcore.DPanicLevel))
	}

	log := zap.New(zapcore.NewCore(enc, sink, __zapLevel(config.Level)), opts...)
	return log
}

func NewLog(config *LogConfig) *zap.Logger {
	if "prod" == config.Mode {
		return __zapLog(config)
	}

	opts := __zapConfig(config)
	log, _ := opts.Build()
	return log
}
