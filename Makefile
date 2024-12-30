################################################################################

PROTOC ?= protoc
PROTO_DIR ?= ./shared/proto
PROTO_GO_OUT ?= ./proto/go

################################################################################

.PHONY: generate
generate:
	mkdir -p $(PROTO_GO_OUT) || true
	$(PROTOC) \
		-I=$(PROTO_DIR) \
		--go_out=$(PROTO_GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/capacity.proto

.PHONY: clean
clean:
	rm -rf $(PROTO_GO_OUT)
