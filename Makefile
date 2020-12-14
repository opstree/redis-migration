.PHONY: all build lint check-fmt vet golangci-lint

GIT_COMMIT := $(shell git rev-parse --short HEAD)
IMAGE_NAME ?= quay.io/opstree/redis-migrator
VERSION := $(shell cat ./VERSION)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -o redis-migrator

golangci-lint:
	golangci-lint run ./...

image:
	docker build -t ${IMAGE_NAME}:${VERSION} -f Dockerfile .
