################################################################################

PROTOC ?= protoc

################################################################################

PROTO_GO_OUT ?= ./proto/go

################################################################################

PROTO_DIR ?= ./shared/proto
PROTO_FILES ?= $(wildcard $(PROTO_DIR)/*.proto)

################################################################################

proto-setup: proto-clean
	mkdir -p $(PROTO_GO_OUT)

.PHONY: generate
proto-go: proto-setup
	$(PROTOC) \
		-I=$(PROTO_DIR) \
		--go_out=$(PROTO_GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

.PHONY: clean
proto-clean:
	rm -rf $(PROTO_GO_OUT)
