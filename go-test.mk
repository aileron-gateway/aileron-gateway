INCLUDED += go-test # Basename of this makefile.
.DEFAULT_GOAL := go-test-help # Basename + "-help"
################################################################################
define go-test.mk
REQUIREMENTS:
  - go   : `go` command must be available.
  - qemu : QEMU User space emulator must be available for `go-test-qemu` target.

TARGETS:
  - go-test-help : show help message.
  - go-test      : run go test.
  - go-test-bin  : generate test binary.
  - go-test-qemu : run go test using qemu emulator.

VARIABLES:
  - GO_TEST_TARGET   : go test target. (Default "./...")
  - GO_TEST_OUTPUT   : test binary output directory. (Default "_output/bin-test/$(GOOS)_$(GOARCH)/")
  - GO_TEST_FLAGS    : go test flags (Default "-timeout 300s -cover -covermode=atomic")
  - GO_TEST_TAGS     : tags passed to the -tags. (Default "")
  - GO_TEST_EXTRAS   : extra flags for go test. (Default "")
  - GO_TEST_COVERAGE : coverage profile output path. (Default "_output/coverage.txt")

REFERENCES:
  - https://pkg.go.dev/cmd/go
  - https://pkg.go.dev/cmd/go/internal/test
	- https://pkg.go.dev/internal/platform
  - https://www.qemu.org/docs/master/user/main.html
endef
################################################################################


export CGO_ENABLED ?= 0

export GOOS ?= $(shell go env GOOS)
ifeq ($(GOOS),)
GOOS := $(shell go env GOOS)
endif

export GOARCH ?= $(shell go env GOARCH)
ifeq ($(GOARCH),)
GOARCH := $(shell go env GOARCH)
endif

OS_ARCH := $(GOOS)-$(GOARCH)
ifeq ($(GOARCH),arm)
OS_ARCH := $(GOOS)-$(GOARCH)$(GOARM)
endif

GO_TEST_TARGET ?= ./...
GO_TEST_OUTPUT ?= _output/bin-test/$(OS_ARCH)/
GO_TEST_FLAGS ?= -timeout 300s -cover -covermode=atomic
GO_TEST_TAGS ?=
GO_TEST_EXTRAS ?=
GO_TEST_COVERAGE ?= _output/coverage.txt

QEMU_TARGET := $(shell go list -f '{{.Dir}}' $(GO_TEST_TARGET))
QEMU_CMD_386 := qemu-i386
QEMU_CMD_amd64 := qemu-x86_64
QEMU_CMD_arm := qemu-arm
QEMU_CMD_arm64 := qemu-aarch64
QEMU_CMD_loong64 := qemu-loong64
QEMU_CMD_mips := qemu-mips
QEMU_CMD_mips64 := qemu-mips64
QEMU_CMD_mips64le := qemu-mips64el
QEMU_CMD_mipsle := qemu-mipsel
QEMU_CMD_ppc64 := qemu-ppc64
QEMU_CMD_ppc64le := qemu-ppc64le
QEMU_CMD_riscv64 := qemu-riscv64
QEMU_CMD_s390x := qemu-s390x
QEMU_CMD_sparc64 := qemu-sparc64


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-test-help`                                                  #
#  Description: Show help message.                                             #
#                                                                              #
.PHONY: go-test-help go-test.mk
go-test-help go-test.mk:
	$(info $(go-test.mk))
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-test`                                                       #
#  Description: Test go packages.                                              #
#  Examples:                                                                   #
#    make go-test                                                              #
#    make go-test GO_TEST_TARGET=./foo/...                                     #
#                                                                              #
GO_TEST_CMD := go test
GO_TEST_CMD += $(GO_TEST_FLAGS)
GO_TEST_CMD += -tags="$(GO_TEST_TAGS)"
GO_TEST_CMD += $(GO_TEST_EXTRAS)
GO_TEST_CMD += -coverprofile=$(GO_TEST_COVERAGE)

.PHONY: go-test
go-test:
	$(info INFO: GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED))
	@mkdir -p $(dir $(GO_TEST_COVERAGE))
	@for target in $(GO_TEST_TARGET); do \
	echo ""; \
	echo "INFO: Testing $$target"; \
	$(GO_TEST_CMD) $$target; \
	done
ifneq ($(GO_TEST_COVERAGE),)
	@go tool cover -html=$(GO_TEST_COVERAGE) -o $(basename $(GO_TEST_COVERAGE)).html
	@go tool cover -func=$(GO_TEST_COVERAGE) -o $(basename $(GO_TEST_COVERAGE)).func.txt
endif
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-test-bin`                                                   #
#  Description: Build test binaries.                                           #
#  Examples:                                                                   #
#    make go-test-bin                                                          #
#    make go-test-bin GO_TEST_TARGET=./foo/...                                 #
#                                                                              #
.PHONY: go-test-bin
go-test-bin:
	$(info INFO: GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED))
	$(info INFO: Output binaries in $(GO_TEST_OUTPUT))
	@for target in $(GO_TEST_TARGET); do \
	echo ""; \
	echo "INFO: Building $$target"; \
	$(GO_TEST_CMD) $(EXTRA_ARGS) -c -o $(GO_TEST_OUTPUT) $$target; \
	done
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-test-qemu`                                                  #
#  Description: Test go packages using QEMU emulator.                          #
#  Examples:                                                                   #
#    make go-test-qemu                                                         #
#    make go-test-qemu GO_TEST_TARGET=./foo/...                                #
#    make go-test-qemu GOARCH=arm64                                            #
#                                                                              #
.PHONY: go-test-qemu
go-test-qemu: 
	$(info INFO: GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED))
	@for target in $(QEMU_TARGET); do \
	echo ""; \
	echo "INFO: Testing $$target"; \
	$(GO_TEST_CMD) $(EXTRA_ARGS) -c -o $$target $$target; \
	find $$target -name "*.test" | xargs -i bash -c "cd $$target && $(QEMU_CMD_$(GOARCH)) {}; rm -f {}"; \
	done
#______________________________________________________________________________#
