.PHONY: all elasticdump img clean

DOCKER   ?= docker
GOLANG_VERSION ?= 1.17
IMAGE    ?= elasticdump
VERSION  ?= latest

all: elasticdump

elasticdump:
	cd cmd/$@; go build -v -o ../../$@

img: 
	$(DOCKER) build --pull \
		--build-arg GOLANG_VERSION=$(GOLANG_VERSION) \
		--tag $(IMAGE):$(VERSION) \
		--file docker/Dockerfile \
		.

clean:
	rm -fr elasticdump
