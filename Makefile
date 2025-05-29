############################################################
#                      Global Configs                      #
############################################################

SHELL := /bin/bash -euo pipefail

include _makefiles/nfpm.mk
include _makefiles/go-test.mk
include _makefiles/go-licenses.mk
include _makefiles/scanoss.mk

.DEFAULT_GOAL := help

# Load .env file if exist.
ifneq (,$(wildcard .env))
    include .env
endif

############################################################
#                          Build                           #
############################################################

# Check build flags with "go help build" or https://pkg.go.dev/cmd/go.
# Check compile flags with "go tool compile -help" or https://pkg.go.dev/cmd/compile.
BUILD_FLAGS ?= -trimpath
BUILD_LDFLAGS ?= -w -s -extldflags \"-static\"
BUILD_GCFLAGS ?= #-m=0
BUILD_TAGS ?= #netgo,osusergo
CGO_ENABLED ?= 0

.PHONY: build
build:
	$(info INFO: GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED))
	export CGO_ENABLED=$(CGO_ENABLED) && \
	go build $(BUILD_FLAGS) -ldflags="$(BUILD_LDFLAGS)" -tags="$(BUILD_TAGS)" -gcflags="$(BUILD_GCFLAGS)" -o ./ ./cmd/aileron/

.PHONY: rpm deb apk archlinux
rpm deb apk archlinux:
	$(MAKE) -C packaging $@

.PHONY: proto proto-clean proto-lint
proto proto-clean proto-lint:
	$(MAKE) -C proto $@

############################################################
#                          Testing                         #
############################################################

OUTPUT_PATH := ./_output/
WHAT ?= ./...

.PHONY: test
test:
	mkdir -p $(OUTPUT_PATH)
	go test -v -cover -timeout 60s -covermode=atomic -coverprofile=$(OUTPUT_PATH)coverage.out $(WHAT)
	sed -i.bak -E '/(testutil|apis|examples)/d' $(OUTPUT_PATH)coverage.out
	go tool cover -html=$(OUTPUT_PATH)coverage.out -o $(OUTPUT_PATH)coverage.html
	go tool cover -func=$(OUTPUT_PATH)coverage.out -o $(OUTPUT_PATH)coverage.txt
	@echo ==================================================
	@cat $(OUTPUT_PATH)coverage.txt
	@echo ==================================================

.PHONY: integration
integration:
	go test -v -tags=integration -timeout 180s ./test/integration/...

.PHONY: example
example:
	go test -v -tags=example -timeout 60s ./test/example/...

############################################################
#                         Analysis                         #
############################################################

.PHONY: lint
lint:
ifeq (,$(shell which golangci-lint 2>/dev/null))
	go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
endif
	go mod tidy
	golangci-lint run ./...

.PHONY: spell
spell:
ifeq (,$(shell which misspell 2>/dev/null))
	go install "github.com/client9/misspell/cmd/misspell@latest"
endif
	misspell -i importas ./

.PHONY: vuln
vuln:
ifeq (,$(shell which govulncheck 2>/dev/null))
	go install "golang.org/x/vuln/cmd/govulncheck@latest"
endif
	govulncheck ./...

.PHONY: sbom
sbom:
	mkdir -p $(OUTPUT_PATH)
ifeq (,$(shell which cyclonedx-gomod 2>/dev/null))
	go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
endif
	cyclonedx-gomod mod -licenses -type library -json -output $(OUTPUT_PATH)sbom.json ./

.PHONY: licenseheader
licenseheader:
ifeq (,$(shell which addlicense 2>/dev/null))
	go install github.com/google/addlicense@latest
endif
	addlicense -check -f ./LICENSE_HEADER.txt -ignore "apis/**" $(shell find . -name "*.go")


############################################################
#                           Help                           #
############################################################

define HELP_MESSAGE
RULES
 help        : show this message.
 build       : build binary.
 test        : run unit tests.
 integration : run integration tests.
 e2e         : run e2e tests.
 example     : test examples in _example.
 lint        : lint with golang-lint.
 spell       : spell check with misspell.
 vuln        : vulnerbility check with govulncheck.
 sbom        : generate cyclonedx SBOM.

RULES imported from proto/Makefile
 proto       : generate go codes from proto files.
 proto-clean : delete generated go codes.
 proto-lint  : run linter for proto files.

RULES imported from packaging/Makefile
 rpm         : generate .rpm package.
 deb         : generate .deb package.
 apk         : generate .apk package.
 archlinux   : generate archlinux package.

VARIABLES
 BUILD_FLAGS   : go build option (Default '-trimpath').
 BUILD_LDFLAGS : go build option given to -ldflags (Default '-w -s -extldflags \"-static\"').
 BUILD_GCFLAGS : go build option given to -gcflags (Default '').
 BUILD_TAGS    : go build option given to -tags (Default '').
endef
export HELP_MESSAGE

.PHONY: help
help: 
	@echo "$${HELP_MESSAGE}"
