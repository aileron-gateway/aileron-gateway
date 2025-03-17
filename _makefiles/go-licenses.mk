INCLUDED += go-licenses # Basename of this makefile.
.DEFAULT_GOAL := go-licenses-help # Basename + "-help"
################################################################################
define go-licenses.mk
REQUIREMENTS:
  - go-licenses: `go-licenses` command must be available.
  - go         : `go` command must be available for `go-licenses-install`.

TARGETS:
  - go-licenses-help    : show help message.
  - go-licenses-install : install go-licenses.
  - go-licenses         : run go-licenses command with given args.
  - go-licenses-run     : check licenses and generate licenses file.

VARIABLES:
  - GO_LICENSES_CMD           : go-licenses binary path. (Default "$(GOBIN)go-licenses")
  - GO_LICENSES_VERSION       : go-licenses version to install. (Default "latest")
  - GO_LICENSES_TARGET        : target for check and report. (Default "./...")
  - GO_LICENSES_OPTION_CHECK  : command line option for check. (Default "--allowed_licenses=MIT,Apache-2.0,BSD-2-Clause,BSD-3-Clause,BSD-4-Clause")
  - GO_LICENSES_OPTION_REPORT : command line option for report. (Default "")

REFERENCES:
  - https://github.com/google/go-licenses
  - https://github.com/google/go-licenses/blob/master/licenses/types.go
  - https://pkg.go.dev/github.com/google/go-licenses
  - https://pkg.go.dev/github.com/google/go-licenses/v2
endef
################################################################################


GO_LICENSES_CMD ?= $(GOBIN)go-licenses
GO_LICENSES_VERSION ?= latest
GO_LICENSES_TARGET ?= ./...
GO_LICENSES_OPTION_CHECK ?= --allowed_licenses=MIT,Apache-2.0,BSD-2-Clause,BSD-3-Clause,BSD-4-Clause
GO_LICENSES_OPTION_REPORT ?= 


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-licenses-help`                                              #
#  Description: Show help message.                                             #
#                                                                              #
.PHONY: go-licenses-help
go-licenses-help:
	$(info $(go-licenses.mk))
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-licenses-install`                                           #
#  Description: Install go-licenses using `go install`.                        #
#                                                                              #
.PHONY: go-licenses-install
go-licenses-install:
ifeq ("go-licenses-install","$(MAKECMDGOALS)")
	go install "github.com/google/go-licenses/v2@$(GO_LICENSES_VERSION)"
else
ifeq (,$(shell which $(GO_LICENSES_CMD) 2>/dev/null))
	go install "github.com/google/go-licenses/v2@$(GO_LICENSES_VERSION)"
endif
endif
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-licenses ARGS=""`                                           #
#  Description: Run go-licenses with given arguments.                          #
#  Examples:                                                                   #
#    make go-licenses ARGS="--help"                                            #
#    make go-licenses ARGS="check --help"                                      #
#    make go-licenses ARGS="report --help"                                     #
#                                                                              #
.PHONY: go-licenses
go-licenses: go-licenses-install
	$(GO_LICENSES_CMD) $(ARGS)
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make go-licenses-run`                                               #
#  Description: Check licenses.                                                #
#                                                                              #
.PHONY: go-licenses-run
go-licenses-run: go-licenses-install
	$(GO_LICENSES_CMD) check $(GO_LICENSES_OPTION_CHECK) $(GO_LICENSES_TARGET)
	$(GO_LICENSES_CMD) report $(GO_LICENSES_OPTION_REPORT) $(GO_LICENSES_TARGET)
#______________________________________________________________________________#
