SHELL := /bin/bash

GOPATH ?= $(shell go env GOPATH)
# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

GO                  := GO111MODULE=on go
GOBUILD             := $(GO) build $(BUILD_FLAG)  

PACKAGE_LIST        := go list ./...| grep -vE "cmd"
PACKAGES            := $$($(PACKAGE_LIST))

CURDIR := $(shell pwd)
export PATH := $(CURDIR)/bin/:$(PATH)

proto:
	cd proto && ./generate_go.sh

http: 
	$(GOBUILD) -o bin ./cmd/http

rpc: 
	$(GOBUILD) -o bin ./cmd/rpc

all: proto http rpc