# Copyright 2021 Shine Xia <shine.xgh@gmail.com>.
#
# Licensed under the MIT License.

ARG GOLANG_VERSION=1.23-alpine
FROM golang:$GOLANG_VERSION as build

COPY ./ /build

ARG GOPROXY=
ENV GOPROXY=$GOPROXY
ENV CGO_ENABLED=0

ARG LDFLAGS="-s -w"
RUN cd /build && go build -v -ldflags="$LDFLAGS" -buildvcs=false -o /build/elasticdump

FROM alpine:latest

COPY --from=build /build/elasticdump /usr/bin/elasticdump

ENTRYPOINT ["elasticdump"]
