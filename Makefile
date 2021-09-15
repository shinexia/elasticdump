.PHONY: all elasticdump docker-build img clean

DOCKER         ?= docker
GOLANG_VERSION ?= 1.17
IMAGE          ?= elasticdump
VERSION        ?= latest
OS_VERSION     ?= linux
ARCH_VERSION   ?= amd64

all: elasticdump

elasticdump:
	cd cmd/$@; CGO_ENABLED=0 go build -v -o ../../$@

docker-build:
	$(DOCKER) run --rm --name elasticdump-build -it \
		-v $(shell pwd):/go/src/github.com/shinexia/elasticdump \
		--workdir /go/src/github.com/shinexia/elasticdump/cmd/elasticdump \
		--user $(shell id -u):$(shell id -g) \
		--env XDG_CACHE_HOME=/tmp/.cache \
		--env GOOS=$(OS_VERSION) \
		--env GOARCH=amd64 \
		--env CGO_ENABLED=0 \
		golang:$(GOLANG_VERSION)-alpine \
		go build -v -o ../../elasticdump-$(OS_VERSION)-$(ARCH_VERSION)

img: 
	$(DOCKER) build --pull \
		--build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
		--tag $(IMAGE):$(VERSION) \
		--file docker/Dockerfile \
		.

clean:
	rm -fr elasticdump*
