# Copyright 2026 flc1125
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0

GO ?= go
TIMEOUT ?= 120
TOOLS_MOD_DIR := ./internal/tools
ALL_GO_MOD_DIRS := $(filter-out $(TOOLS_MOD_DIR), $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort))
ALL_DOCS := $(shell find . -type f -name '*.md' | sort)
TOOLS_DIR := $(CURDIR)/.tools
CROSSLINK := $(TOOLS_DIR)/crosslink
GOTMPL := $(TOOLS_DIR)/gotmpl
MULTIMOD := $(TOOLS_DIR)/multimod
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint
GOLANGCI_LINT_MODULE := github.com/golangci/golangci-lint/v2
GOLANGCI_LINT_PACKAGE := $(GOLANGCI_LINT_MODULE)/cmd/golangci-lint
GOLANGCI_LINT_VERSION := $(shell cd $(TOOLS_MOD_DIR) && $(GO) list -m -f '{{.Version}}' $(GOLANGCI_LINT_MODULE))
MISSPELL := $(TOOLS_DIR)/misspell
GOVULNCHECK := $(TOOLS_DIR)/govulncheck
COMMIT ?= HEAD
REMOTE ?= origin

.DEFAULT_GOAL := precommit

.PHONY: precommit ci fmt fmt-check tidy generate build vet lint golangci-lint \
	misspell govulncheck test test-race \
	license-check check-clean-work-tree tools crosslink gowork \
	multimod-verify multimod-prerelease print-tags create-tags push-tags

precommit: fmt tidy generate vet lint test multimod-verify
ci: fmt-check generate vet lint test license-check check-clean-work-tree

$(TOOLS_DIR):
	mkdir -p $@

$(TOOLS_DIR)/%: $(TOOLS_MOD_DIR)/go.mod | $(TOOLS_DIR)
	cd $(TOOLS_MOD_DIR) && $(GO) build -o $@ $(PACKAGE)

$(CROSSLINK): PACKAGE=github.com/flc1125/go-build-tools/crosslink
$(GOTMPL): PACKAGE=github.com/flc1125/go-build-tools/gotmpl
$(MULTIMOD): PACKAGE=github.com/flc1125/go-build-tools/multimod
$(MISSPELL): PACKAGE=github.com/client9/misspell/cmd/misspell
$(GOVULNCHECK): PACKAGE=golang.org/x/vuln/cmd/govulncheck

$(GOLANGCI_LINT): $(TOOLS_MOD_DIR)/go.mod Makefile | $(TOOLS_DIR)
	GOBIN=$(TOOLS_DIR) $(GO) install $(GOLANGCI_LINT_PACKAGE)@$(GOLANGCI_LINT_VERSION)

tools: $(CROSSLINK) $(GOTMPL) $(MULTIMOD) $(GOLANGCI_LINT) $(MISSPELL) $(GOVULNCHECK)

fmt:
	gofmt -w $$(find . -type f -name '*.go' ! -path './.git/*')

fmt-check:
	@files=$$(gofmt -l $$(find . -type f -name '*.go' ! -path './.git/*')); \
	if [ -n "$$files" ]; then echo "Unformatted Go files:"; echo "$$files"; exit 1; fi

tidy:
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "$(GO) mod tidy in $$dir"; \
		(cd "$$dir" && $(GO) mod tidy); \
	done
	cd $(TOOLS_MOD_DIR) && $(GO) mod tidy

generate:
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "$(GO) generate in $$dir"; \
		(cd "$$dir" && $(GO) generate ./...); \
	done

build: generate
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "$(GO) build in $$dir"; \
		(cd "$$dir" && $(GO) build ./...); \
	done

vet:
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "$(GO) vet in $$dir"; \
		(cd "$$dir" && $(GO) vet ./...); \
	done

lint: misspell golangci-lint govulncheck

golangci-lint: generate | $(GOLANGCI_LINT)
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "golangci-lint in $$dir"; \
		(cd "$$dir" && $(GOLANGCI_LINT) run --fix && $(GOLANGCI_LINT) run); \
	done

misspell: | $(MISSPELL)
	$(MISSPELL) -w $(ALL_DOCS)

govulncheck: | $(GOVULNCHECK)
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
		echo "govulncheck in $$dir"; \
		(cd "$$dir" && $(GOVULNCHECK) ./...); \
	done

test:
	@set -e; for dir in $(ALL_GO_MOD_DIRS); do \
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
	$(CROSSLINK) work --root=$(CURDIR) --go=1.26

multimod-verify: $(MULTIMOD)
	$(MULTIMOD) verify --versioning-file $(CURDIR)/versions.yaml

multimod-prerelease: multimod-verify $(MULTIMOD)
	$(MULTIMOD) prerelease --module-set-name tools --versioning-file $(CURDIR)/versions.yaml

print-tags: multimod-verify $(MULTIMOD)
	$(MULTIMOD) tag --module-set-name tools --preview-tags

create-tags: multimod-verify $(MULTIMOD)
	$(MULTIMOD) tag --module-set-name tools --commit-hash $(COMMIT) --print-tags

push-tags: multimod-verify $(MULTIMOD)
	@set -e; target_commit=$$(git rev-parse $(COMMIT)); \
	for tag in $$($(MULTIMOD) tag --module-set-name tools --preview-tags); do \
		if ! git show-ref --verify --quiet "refs/tags/$$tag"; then \
			echo "Missing local tag: $$tag"; exit 1; \
		fi; \
		tag_commit=$$(git rev-list -n 1 "$$tag"); \
		if [ "$$tag_commit" != "$$target_commit" ]; then \
			echo "Tag $$tag points to $$tag_commit, expected $$target_commit"; exit 1; \
		fi; \
	done; \
	for tag in $$($(MULTIMOD) tag --module-set-name tools --preview-tags); do \
		echo "Pushing $$tag"; \
		git push $(REMOTE) "$$tag"; \
	done
