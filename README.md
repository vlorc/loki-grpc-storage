# [Loki grpc Storage](https://github.com/vlorc/loki-grpc-storage)

[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/vlorc/loki-grpc-storage)](https://goreportcard.com/report/github.com/vlorc/loki-grpc-storage)
[![GoDoc](https://godoc.org/github.com/vlorc/loki-grpc-storage?status.svg)](https://godoc.org/github.com/vlorc/loki-grpc-storage)
[![Build Status](https://travis-ci.org/vlorc/loki-grpc-storage.svg?branch=master)](https://travis-ci.org/vlorc/loki-grpc-storage)
[![Coverage Status](https://coveralls.io/repos/github/vlorc/loki-grpc-storage/badge.svg?branch=master)](https://codecov.io/gh/vlorc/loki-grpc-storage)

A loki third party storage experimental project

# Features

+ filesystem
+ memory
+ http
+ qiniu
+ baidu
+ aliyun

# Quick Start

**filesystem**

```shell
./storage -store.url /tmp/loki/storage
```

**qiniu**

```shell
./storage -store.driver qiniu    \
    -store.url https://xxxx.cdn.com  \
    -store.access xxxx             \
    -store.secret xxxx            \
    -store.bucket log             \
    -store.flag https,cdn,private
```

## License

This project is under the apache License. See the LICENSE file for the full license text.

# Keyword

**loki storage,third party**

