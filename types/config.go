// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package types

const UserAgent = "storage"

type Config struct {
	Log    LogConfig    `flag:"log"`
	Chunk  ChunkConfig  `flag:"chunk"`
	Store  StoreConfig  `flag:"store"`
	Server ServerConfig `flag:"server"`
}

type ServerConfig struct {
	Host string `flag:"host,0.0.0.0,server host"`
	Port string `flag:"port,5783,server port"`
}

type LogConfig struct {
	Level  string `flag:"level,info,log level"`
	Caller bool   `flag:"caller,,log caller"`
	Trace  bool   `flag:"trace,,log trace"`
	Mode   string `flag:"mode,prod,log mode"`
}

type ChunkConfig struct {
	Level    string `flag:"level,debug,chunk level"`
	Mode     string `flag:"mode,prod,chunk mode"`
	Min      int    `flag:"min,12,chunk minimum"`
	Parallel int    `flag:"parallel,0,chunk parallel"`
}

type StoreConfig struct {
	Level  string `flag:"level,debug,store level"`
	Mode   string `flag:"mode,prod,store mode"`
	Driver string `flag:"driver,fs,store driver"`
	Name   string `flag:"name,,store name"`
	Url    string `flag:"url,{tmpdir},store url"`
	Access string `flag:"access,,store access"`
	Secret string `flag:"secret,,store secret"`
	Token  string `flag:"token,,store token"`
	Bucket string `flag:"bucket,,store bucket"`
	Region string `flag:"region,,store region"`
	Flag   string `flag:"flag,,store flag"`
}
