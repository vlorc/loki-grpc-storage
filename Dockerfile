// Copyright 2021 vlorc. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

FROM golang:latest as builder
WORKDIR /app

ARG GOPROXY="https://goproxy.cn"

COPY . .
RUN go mod download
RUN go build main/main.go -o storage

FROM alpine:latest
WORKDIR /opt/storage

COPY --from=builder /app/storage .
RUN chmod 755 /opt/storage/storage

EXPOSE 5783

ENTRYPOINT ["/opt/storage/storage"]