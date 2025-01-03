################################################################################

PROTOC ?= protoc

################################################################################

PROTO_GO_OUT ?= ./proto/go

################################################################################

PROTO_DIR ?= ./shared/proto
PROTO_FILES ?= $(wildcard $(PROTO_DIR)/*.proto)

################################################################################

.PHONY: setup
setup:
	cp hack/github/pre-push .git/hooks/pre-push

.PHONY: clean
clean: proto-clean
	rm -f .git/hooks/pre-push

.PHONY: proto-setup
proto-setup: proto-clean
	mkdir -p $(PROTO_GO_OUT)

.PHONY: proto-go
proto-go: proto-setup
	$(PROTOC) \
		-I=$(PROTO_DIR) \
		--go_out=$(PROTO_GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

.PHONY: proto-clean
proto-clean:
	rm -rf $(PROTO_GO_OUT)
