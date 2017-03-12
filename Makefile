SHELL          := /bin/bash
VERSION        := v0.1.0
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC            := $(wildcard *.go) $(wildcard src/*/*.go)

.PHONY: all go-get

all: sisito-api

go-get:
	go get gopkg.in/gin-gonic/gin.v1
	go get github.com/BurntSushi/toml
	go get github.com/go-sql-driver/mysql
	go get gopkg.in/gorp.v1

sisito-api: $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) go build -a -tags netgo -installsuffix netgo -o sisito-api
else
	GOPATH=$(RUNTIME_GOPATH) go build -o sisito-api
endif
