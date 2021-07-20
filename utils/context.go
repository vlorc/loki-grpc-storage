// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Stopper interface {
	Stop()
}

func Process(sig ...os.Signal) context.Context {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sig...)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signals
		cancel()
	}()

	return ctx
}

func OnExit(ss ...Stopper) context.Context {
	ctx := Process(syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	if len(ss) == 0 {
		return ctx
	}

	go func() {
		<-ctx.Done()
		for _, s := range ss {
			s.Stop()
		}
	}()

	return ctx
}
