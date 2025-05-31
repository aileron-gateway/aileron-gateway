############################################################
#                      Global Configs                      #
############################################################

.DEFAULT_GOAL := proto

ifneq (,$(wildcard .env))
	include .env
endif

############################################################
#                          Proto                           #
############################################################

PROTO_PATH += --proto_path ./
PROTOC_OPTS += --plugin=protoc-gen-go=$(shell which protoc-gen-go)
PROTOC_OPTS += --go_out=../ --go_opt=module="github.com/aileron-gateway/aileron-gateway"

.PHONY: proto
proto: protoc-gen-go
	protoc $(PROTO_PATH) $(PROTOC_OPTS) $(shell find ./ -type f -name "*.proto" -not -path "./buf/*")
	@git update-index --assume-unchanged $(shell find ../apis/ -type f)

.PHONY: proto-clean
proto-clean:
	rm -rf $(shell find ../apis/ -type f -name "*.pb.go")

.PHONY: proto-lint
proto-lint: protolint
	protolint -config_path .protolint.yaml kernel/
	protolint -config_path .protolint.yaml core/
	protolint -config_path .protolint.yaml app/

############################################################
#                          Tools                           #
############################################################

.PHONY: protoc-gen-go
protoc-gen-go:
ifeq (, $(shell which protoc-gen-go 2>/dev/null))
	go install "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.4"
endif

.PHONY: protolint
protolint:
ifeq (,$(shell which protolint 2>/dev/null))
	go install "github.com/yoheimuta/protolint/cmd/protolint@latest"
endif
