VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main. version=$(VERSION)' -X 'main. revision=$(REVISION)'
SRC_DIR := ./
BIN_DIR := bin
BINARY := bin/alfred-k8s
WORKFLOW_DIR := "$${HOME}/Library/Application Support/Alfred/Alfred.alfredpreferences/workflows/user.workflow.2FA36703-F09E-4B6F-9CE1-F0E3944B7DDB"
ASSETS_DIR := assets
ARTIFACT_DIR := .artifact
ARTIFACT := ${ARTIFACT_DIR}/alfred-k8s.alfredworkflow

GOLANGCI_LINT_VERSION := v1.22.2
export GO111MODULE=on

## Build binaries on your environment
build:
	CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ${BINARY} ./${SRC_DIR}

## Format source codes
fmt:
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint:
	@(if ! type golangci-lint >/dev/null 2>&1; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION} ;fi)
	golangci-lint run ./...

## Build macos binaries
darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -ldflags "${LDFLAGS} -s -w" -o ${BINARY} ./${SRC_DIR}

## Run tests for my project
test:
	go test -v ./...

## Install Binary and Assets to Workflow Directory
install: build
	@(cp ${BINARY} ${WORKFLOW_DIR}/)
	@(cp ${ASSETS_DIR}/*  ${WORKFLOW_DIR}/)

## GitHub Release and uploads artifacts
release: darwin
	@(if ! type ghr >/dev/null 2>&1; then go get -u github.com/tcnksm/ghr ;fi)
	@(if [ ! -e ${ARTIFACT_DIR} ]; then mkdir ${ARTIFACT_DIR} ; fi)
	@(cp ${BINARY} ${ARTIFACT_DIR})
	@(cp ${ASSETS_DIR}/* ${ARTIFACT_DIR})
	@(zip -j ${ARTIFACT} ${ARTIFACT_DIR}/*)
	@ghr -replace ${VERSION} ${ARTIFACT}

## Clean Binary
clean:
	rm -f ${BIN_DIR}/*
	rm -f ${ARTIFACT_DIR}/*

## Show help
help:
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)
	@make2help $(MAKEFILE_LIST)

.PHONY: build test lint fmt darwin clean help
