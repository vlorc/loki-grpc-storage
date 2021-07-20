// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package main

import (
	"github/vlorc/loki-grpc-storage/server"
	"github/vlorc/loki-grpc-storage/types"
	"github/vlorc/loki-grpc-storage/utils"
)

func main() {
	conf := &types.Config{}

	utils.Flag(conf)

	srv := server.NewServer(conf)

	utils.OnExit(srv)
	_ = srv.Serve()
}
