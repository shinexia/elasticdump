DOCKER         ?= docker
GOLANG_VERSION ?= 1.23-alpine
IMAGE          ?= elasticdump
VERSION        ?= latest
GOOS           ?= linux
GOARCH         ?= amd64
GOPROXY        ?= $(shell printenv GOPROXY)
LDFLAGS        ?= "-s -w"

.PHONY: all
all: elasticdump

.PHONY: elasticdump
elasticdump:
	go build -v -o $@

.PHONY: docker-build
docker-build:
	$(DOCKER) run --rm --name elasticdump-build -it \
		-v $(shell pwd):/build \
		--workdir /build \
		--user $(shell id -u):$(shell id -g) \
		--env XDG_CACHE_HOME=/tmp/.cache \
		--env GOOS=$(GOOS) \
		--env GOARCH=$(GOARCH) \
		--env GOPROXY=$(GOPROXY) \
		--env CGO_ENABLED=0 \
		golang:$(GOLANG_VERSION) \
		go build -v -ldflags=$(LDFLAGS) -buildvcs=false -o elasticdump-$(GOOS)-$(GOARCH)

.PHONY: img
img: 
	$(DOCKER) build \
		--build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
		--build-arg GOPROXY=$(GOPROXY) \
		--build-arg LDFLAGS=$(LDFLAGS) \
		--tag $(IMAGE):$(VERSION) \
		--file docker/Dockerfile \
		.

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run -v --timeout 10m ./cmd/... ./pkg/...

.PHONY: test
test:
	go test -count=1 ./cmd/... ./pkg/...

.PHONY: clean
clean:
	rm -fr elasticdump*
