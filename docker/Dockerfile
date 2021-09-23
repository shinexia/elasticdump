# Copyright 2021 Shine Xia <shine.xgh@gmail.com>.
#
# Licensed under the MIT License.

ARG GOLANG_VERSION=1.17
FROM golang:$GOLANG_VERSION-alpine as build

WORKDIR /build
COPY . .

ARG GOPROXY=
ENV GOPROXY=$GOPROXY
ENV CGO_ENABLED=0

ARG LDFLAGS="-s -w"
RUN go build -v -ldflags="$LDFLAGS" -o /build/elasticdump

FROM alpine:latest

COPY --from=build /build/elasticdump /usr/bin/elasticdump

ENTRYPOINT ["elasticdump"]
