# Copyright 2026 flc1125
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0

GO ?= go
MODULE_DIRS := . crosslink gotmpl multimod
TIMEOUT ?= 120
LINT_GO_VERSION ?= go1.25.14
GOLANGCI_LINT_VERSION ?= v2.11.4
GOLANGCI_LINT := GOTOOLCHAIN=$(LINT_GO_VERSION) $(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
TOOLS_DIR := $(CURDIR)/.tools
CROSSLINK := $(TOOLS_DIR)/crosslink
GOTMPL := $(TOOLS_DIR)/gotmpl
MULTIMOD := $(TOOLS_DIR)/multimod
COMMIT ?= HEAD
REMOTE ?= origin

.DEFAULT_GOAL := precommit

.PHONY: precommit ci fmt fmt-check tidy generate build vet lint test test-race \
	license-check check-clean-work-tree tools crosslink gowork \
	multimod-verify multimod-prerelease print-tags push-tags

precommit: fmt tidy generate vet lint test multimod-verify
ci: fmt-check generate vet lint test license-check

$(TOOLS_DIR):
	mkdir -p $@

$(CROSSLINK): | $(TOOLS_DIR)
	cd crosslink && $(GO) build -o $@ .

$(GOTMPL): | $(TOOLS_DIR)
	cd gotmpl && $(GO) build -o $@ .

$(MULTIMOD): | $(TOOLS_DIR)
	cd multimod && $(GO) build -o $@ .

tools: $(CROSSLINK) $(GOTMPL) $(MULTIMOD)

fmt:
	gofmt -w $$(find . -type f -name '*.go' ! -path './.git/*')

fmt-check:
	@files=$$(gofmt -l $$(find . -type f -name '*.go' ! -path './.git/*')); \
	if [ -n "$$files" ]; then echo "Unformatted Go files:"; echo "$$files"; exit 1; fi

tidy:
	@set -e; for dir in $(MODULE_DIRS); do \
		echo "$(GO) mod tidy in $$dir"; \
		(cd "$$dir" && $(GO) mod tidy); \
	done

generate:
	@set -e; for dir in $(MODULE_DIRS); do \
		echo "$(GO) generate in $$dir"; \
		(cd "$$dir" && $(GO) generate ./...); \
	done

build: tools
	$(GO) build ./...

vet:
	@set -e; for dir in $(MODULE_DIRS); do \
		echo "$(GO) vet in $$dir"; \
		(cd "$$dir" && $(GO) vet ./...); \
	done

lint:
	@set -e; for dir in $(MODULE_DIRS); do \
		echo "golangci-lint in $$dir"; \
		(cd "$$dir" && $(GOLANGCI_LINT) run --config $(CURDIR)/.golangci.yml); \
	done

test:
	@set -e; for dir in $(MODULE_DIRS); do \
		echo "$(GO) test in $$dir"; \
		(cd "$$dir" && $(GO) test -timeout $(TIMEOUT)s $(ARGS) ./...); \
	done

test-race:
	$(MAKE) test ARGS=-race

license-check:
	@missing=$$(find . -type f -name '*.go' ! -path './.git/*' -exec \
		awk 'NR <= 3 && /Copyright The OpenTelemetry Authors|Copyright 2026 flc1125|generated|GENERATED/ { found=1 } END { if (!found) print FILENAME }' {} \;); \
	if [ -n "$$missing" ]; then echo "Missing license headers:"; echo "$$missing"; exit 1; fi

check-clean-work-tree:
	git diff --exit-code

crosslink: $(CROSSLINK)
	$(CROSSLINK) --root=$(CURDIR) --prune

gowork: $(CROSSLINK)
	$(CROSSLINK) work --root=$(CURDIR) --go=1.25

multimod-verify: $(MULTIMOD)
	$(MULTIMOD) verify --versioning-file $(CURDIR)/versions.yaml

multimod-prerelease: multimod-verify $(MULTIMOD)
	$(MULTIMOD) prerelease --module-set-name tools --versioning-file $(CURDIR)/versions.yaml

print-tags: multimod-verify $(MULTIMOD)
	$(MULTIMOD) tag --module-set-name tools --commit-hash $(COMMIT) --print-tags

push-tags: multimod-verify $(MULTIMOD)
	@set -e; for tag in $$($(MULTIMOD) tag --module-set-name tools --commit-hash $(COMMIT) --print-tags); do \
		echo "Pushing $$tag"; \
		git push $(REMOTE) "$$tag"; \
	done
