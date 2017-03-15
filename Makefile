SHELL          := /bin/bash
PROGRAM        := sisito-api
VERSION        := v0.1.3
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC            := $(wildcard *.go) $(wildcard src/*/*.go)

UBUNTU_IMAGE          := docker-go-pkg-build-ubuntu
UBUNTU_CONTAINER_NAME := docker-go-pkg-build-ubuntu-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get gopkg.in/gin-gonic/gin.v1
	go get github.com/BurntSushi/toml
	go get github.com/go-sql-driver/mysql
	go get gopkg.in/gorp.v1

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) go build -a -tags netgo -installsuffix netgo -o $(PROGRAM)
else
	GOPATH=$(RUNTIME_GOPATH) go build -o $(PROGRAM)
endif

.PHONY: clean
clean:
	rm -f $(PROGRAM) pkg/*

.PHONY: package
package: clean $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(UBUNTU_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: docker/build/ubuntu
docker/build/ubuntu: docker/Dockerfile.ubuntu
	docker build -f docker/Dockerfile.ubuntu -t $(UBUNTU_IMAGE) .
