################################################################################

GO ?= /usr/bin/go
PROTOC ?= protoc
PROTOC_OUT ?= $(shell pwd)/internal/proto

################################################################################

$(PROTOC_OUT):
	mkdir -p $@

################################################################################

.PHONY: test
test: generate
	$(GO) test ./... -race -cover

.PHONY: generate
generate: generate-proto
	$(GO) generate ./...

.PHONY: generate-proto
generate-proto: $(PROTOC_OUT)
	$(PROTOC) --proto_path=proto \
		--go_out=$(PROTOC_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTOC_OUT) \
		--go-grpc_opt=paths=source_relative \
		proto/novacloud/v1/*.proto

.PHONY: clean
clean:
	rm -rf $(PROTOC_OUT) $(LOCALBIN)
