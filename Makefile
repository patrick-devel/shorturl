.PHONY: tools
PROJECT_PATH=$(shell pwd)
PROJECT_BIN_PATH=$(PROJECT_PATH)/bin

export GOBIN := $(PROJECT_BIN_PATH)

MOCKGEN_VERSION=v1.6.0

-include .env

tools:
	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
