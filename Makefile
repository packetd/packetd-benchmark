GO ?= go
SHELL = bash

.PHONY: help
help:
	@echo "Make Targets: "
	@echo " mod: Download and tidy dependencies"
	@echo " lint: Lint Go code"
	@echo " install-tools: Install dev tools"

.PHONY: lint
lint:
	gofumpt -w .
	goimports-reviser -project-name "github.com/packetd/packetd-benchmark" ./...
	find ./ -type f \( -iname \*.go -o -iname \*.sh \) | xargs addlicense -v -f LICENSE

.PHONY: mod
mod:
	$(GO) mod download
	$(GO) mod tidy

.PHONY: install-tools
tools:
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/incu6us/goimports-reviser/v3@v3.1.1
	$(GO) install github.com/google/addlicense@latest
