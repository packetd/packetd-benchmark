GO ?= go
SHELL = bash
PKG = github.com/packetd/packetd

VERSION := $(shell cat VERSION)

.PHONY: help
help:
	@echo "Make Targets: "
	@echo " mod: Download and tidy dependencies"
	@echo " lint: Lint Go code"
	@echo " test: Run unit tests"
	@echo " build: Build Go package"
	@echo " tools: Install dev tools"

.PHONY: license
license:
	find ./ -type f \( -iname \*.go -o -iname \*.sh \) | xargs addlicense -v -f LICENSE

.PHONY: lint
lint: license
	gofumpt -w .
	goimports-reviser -project-name "github.com/packetd/packetd-benchmark" ./...

.PHONY: mod
mod:
	$(GO) mod download
	$(GO) mod tidy

.PHONY: tools
tools:
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/incu6us/goimports-reviser/v3@v3.1.1
	$(GO) install github.com/google/addlicense@latest
