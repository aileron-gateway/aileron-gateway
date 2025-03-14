INCLUDED += nfpm # Basename of this makefile.
.DEFAULT_GOAL := nfpm-help # Basename + "-help"
################################################################################
define nfpm.mk
REQUIREMENTS:
  - nfpm : `nfpm` command must be available.
  - go    : `go` command must be available for `nfpm-install`.

TARGETS:
  - nfpm-help    : show help message.
  - nfpm         : run nfpm command with given args.
  - nfpm-install : install nfpm using `go install`.
  - nfpm-rpm     : create rpm packages.
  - nfpm-deb     : create deb packages.
  - nfpm-apk     : create apk packages.
  - nfpm-arch    : create archlinux packages.

VARIABLES:
  - NFPM_CMD       : nfpm binary path. (Default "$(GOBIN)nfpm")
  - NFPM_VERSION   : nfpm version to install. (Default "latest")
  - NFPM_OPTION    : nfpm command line option. (Default "")
  - NFPM_OUTPUT    : package output directory. (Default "./_output/pkg/")
  - NFPM_RPM_ARCH  : rpm target architecture. (Default "386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x")
  - NFPM_DEB_ARCH  : deb target architecture. (Default "386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x")
  - NFPM_APK_ARCH  : apk target architecture. (Default "386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x")
  - NFPM_ARCH_ARCH : archlinux target architecture. (Default "amd64 arm arm5 arm6 arm7 arm64 riscv64")

REFERENCES:
  - https://github.com/goreleaser/nfpm
  - https://nfpm.goreleaser.com/
endef
################################################################################


NFPM_CMD ?= $(GOBIN)nfpm
NFPM_VERSION ?= latest
NFPM_OPTION ?=
NFPM_OUTPUT ?= ./_output/pkg/
NFPM_RPM_ARCH ?= 386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x
NFPM_DEB_ARCH ?= 386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x
NFPM_APK_ARCH ?= 386 amd64 arm arm5 arm6 arm7 arm64 ppc64le riscv64 s390x
NFPM_ARCH_ARCH ?= amd64 arm arm5 arm6 arm7 arm64 riscv64


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-help`                                                     #
#  Description: Show help message.                                             #
#                                                                              #
.PHONY: nfpm-help
nfpm-help:
	$(info $(nfpm.mk))
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-install`                                                  #
#  Description: Install nfpm using `go install`.                               #
#                                                                              #
.PHONY: nfpm-install
nfpm-install:
ifeq ("nfpm-install","$(MAKECMDGOALS)")
	go install "github.com/goreleaser/nfpm/cmd/nfpm@$(NFPM_VERSION)"
else
ifeq (,$(shell which $(NFPM_CMD) 2>/dev/null))
	go install "github.com/goreleaser/nfpm/cmd/nfpm@$(NFPM_VERSION)"
endif
endif
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm ARGS=""`                                                  #
#  Description: Run nfpm with given arguments.                                 #
#  Examples:                                                                   #
#    make nfpm ARGS="--version"                                                #
#    make nfpm ARGS="--help"                                                   #
#                                                                              #
.PHONY: nfpm
nfpm: nfpm-install
	$(NFPM_CMD) $(ARGS)
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-rpm`                                                      #
#  Description: Create rpm packages.                                           #
#  Examples:                                                                   #
#    make nfpm-rpm NFPM_RPM_ARCH="amd64"                                       #
#    make nfpm-rpm NFPM_RPM_ARCH="amd64 arm64"                                 #
#                                                                              #
#                                                                              #
.PHONY: nfpm-rpm
nfpm-rpm: nfpm-install
	mkdir -p $(NFPM_OUTPUT)
	@for target in $(NFPM_RPM_ARCH); do \
	ARCH=$$target nfpm package -p rpm -t $(NFPM_OUTPUT) $(NFPM_OPTION); \
	done
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-deb`                                                      #
#  Description: Create deb packages.                                           #
#  Examples:                                                                   #
#    make nfpm-deb NFPM_DEB_ARCH="amd64"                                       #
#    make nfpm-deb NFPM_DEB_ARCH="amd64 arm64"                                 #
#                                                                              #
#                                                                              #
.PHONY: nfpm-deb
nfpm-deb: nfpm-install
	mkdir -p $(NFPM_OUTPUT)
	@for target in $(NFPM_DEB_ARCH); do \
	ARCH=$$target nfpm package -p deb -t $(NFPM_OUTPUT) $(NFPM_OPTION); \
	done
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-apk`                                                      #
#  Description: Create apk packages.                                           #
#  Examples:                                                                   #
#    make nfpm-apk NFPM_APK_ARCH="amd64"                                       #
#    make nfpm-apk NFPM_APK_ARCH="amd64 arm64"                                 #
#                                                                              #
#                                                                              #
.PHONY: nfpm-apk
nfpm-apk: nfpm-install
	mkdir -p $(NFPM_OUTPUT)
	@for target in $(NFPM_APK_ARCH); do \
	ARCH=$$target nfpm package -p apk -t $(NFPM_OUTPUT) $(NFPM_OPTION); \
	done
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make nfpm-arch`                                                     #
#  Description: Create archlinux packages.                                     #
#  Examples:                                                                   #
#    make nfpm-arch NFPM_ARCH_ARCH="amd64"                                     #
#    make nfpm-arch NFPM_ARCH_ARCH="amd64 arm64"                               #
#                                                                              #
#                                                                              #
.PHONY: nfpm-arch
nfpm-arch: nfpm-install
	mkdir -p $(NFPM_OUTPUT)
	@for target in $(NFPM_ARCH_ARCH); do \
	ARCH=$$target nfpm package -p archlinux -t $(NFPM_OUTPUT) $(NFPM_OPTION); \
	done
#______________________________________________________________________________#
