SHELL=/bin/bash
ROOT=$(shell pwd)
PACKAGE=github.com/dantin/microservice-go

GOMOD := -mod=vendor

GOVER_MAJOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\1/")
GOVER_MINOR := $(shell go version | sed -E -e "s/.*go([0-9]+)[.]([0-9]+).*/\2/")
GO111 := $(shell [ $(GOVER_MAJOR) -gt 1  ] || [ $(GOVER_MAJOR) -eq 1  ] && [ $(GOVER_MINOR) -ge 11  ]; echo $$?)
ifeq ($(GO111), 1)
	$(warning "go below 1.11 does not support modules")
GOMOD :=
endif

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
	@bash $(ROOT)/scripts/clean_vendor.sh

.PHONY: clean
clean:
	@echo "clean up"
