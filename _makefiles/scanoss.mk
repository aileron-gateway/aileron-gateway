SHELL := /bin/bash

INCLUDED += scanoss # Basename of this makefile.
.DEFAULT_GOAL := scanoss-help # Basename + "-help"
################################################################################
define scanoss.mk
REQUIREMENTS:
  - scanoss-py : `scanoss-py` command must be available.
  - pip        : `pip` command must be available for `scanoss-install`.

TARGETS:
  - scanoss-help    : show help message.
  - scanoss-install : install scanoss using `pip install`.
  - scanoss         : run scanoss command with given args.
  - scanoss-run     : run scanning.

VARIABLES:
  - SCANOSS_CMD            : scanoss-py binary path. (Default "scanoss-py")
  - SCANOSS_TARGET         : target of scanning. (Default "./")
  - SCANOSS_OUTPUT         : scan result ouput file. (Default "_output/scanoss-output.json")
  - SCANOSS_OPTION_SCAN    : scanoss-py scan command line option. (Default "--no-wfp-output")
  - SCANOSS_OPTION_INSPECT : scanoss-py inspect command line option. (Default "copyleft")

REFERENCES:
  - https://github.com/scanoss/scanoss.py
  - https://github.com/scanoss/scanoss.py/blob/main/PACKAGE.md
endef
################################################################################


SCANOSS_CMD ?= scanoss-py
SCANOSS_TARGET ?= ./
SCANOSS_OUTPUT ?= _output/scanoss-output.json
SCANOSS_OPTION_SCAN ?= --no-wfp-output
SCANOSS_OPTION_INSPECT ?= copyleft 


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make scanoss-help`                                                  #
#  Description: Show help message.                                             #
#                                                                              #
.PHONY: scanoss-help
scanoss-help:
	$(info $(scanoss.mk))
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make scanoss-install`                                               #
#  Description: Install scanoss using `pip install`                            #
#                                                                              #
.PHONY: scanoss-install
scanoss-install:
ifeq ("scanoss-install","$(MAKECMDGOALS)")
	pip install scanoss[fast_winnowing]
  # pip install scancode-toolkit
else
ifeq (,$(shell which $(SCANOSS_CMD) 2>/dev/null))
	pip install scanoss[fast_winnowing]
  # pip install scancode-toolkit
endif
endif
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make scanoss ARGS=""`                                               #
#  Description: Run scanoss with given arguments.                              #
#  Examples:                                                                   #
#    make scanoss ARGS="--version"                                             #
#    make scanoss ARGS="--help"                                                #
#                                                                              #
.PHONY: scanoss
scanoss: scanoss-install
	$(SCANOSS_CMD) $(ARGS)
#______________________________________________________________________________#


##### ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ #####
##                                                                            ##
#  Usage: `make scanoss-run`                                                   #
#  Description: Run scanning.                                                  #
#                                                                              #
.PHONY: scanoss-run
scanoss-run: scanoss-install
	mkdir -p $(dir $(SCANOSS_OUTPUT))
	$(SCANOSS_CMD) scan $(SCANOSS_OPTION_SCAN) -o $(SCANOSS_OUTPUT) $(SCANOSS_TARGET)
	$(SCANOSS_CMD) inspect $(SCANOSS_OPTION_INSPECT)  -i $(SCANOSS_OUTPUT) -q | grep "{}"
#______________________________________________________________________________#
