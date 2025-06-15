################################################################################

GO ?= /usr/bin/go

LOCALBIN ?= $(shell pwd)/bin
PROTOC_OUT ?= $(shell pwd)/internal/proto

PROTOC ?= $(LOCALBIN)/protoc
PROTOC_VERSION ?= v30.2
PROTOC_GEN_GO ?= $(LOCALBIN)/protoc-gen-go
PROTOC_GEN_GO_VERSION ?= v1.36.6
PROTOC_GEN_GO_GRPC ?= $(LOCALBIN)/protoc-gen-go-grpc
PROTOC_GEN_GO_GRPC_VERSION ?= v1.5.1
MOCKGEN ?= $(LOCALBIN)/mockgen
MOCKGEN_VERSION ?= v0.5.2

################################################################################

$(LOCALBIN):
	mkdir -p $@

$(PROTOC_OUT):
	mkdir -p $@

################################################################################

.PHONY: test
test: generate
	$(GO) test ./... -race -cover

.PHONY: generate
generate: mockgen generate-proto
	PATH=$(LOCALBIN):$(PATH) $(GO) generate ./...

.PHONY: generate-proto
generate-proto: $(PROTOC_OUT) protoc protoc-gen-go protoc-gen-go-grpc
	$(PROTOC) --proto_path=proto \
		--go_out=$(PROTOC_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTOC_OUT) \
		--go-grpc_opt=paths=source_relative \
		proto/novacloud/v1/*.proto

.PHONY: protoc
protoc: $(PROTOC)
$(PROTOC): PROTOC_ARCHIVE = protoc-$(subst v,,$(PROTOC_VERSION))-linux-x86_64.zip
$(PROTOC): PROTOC_URL = https://github.com/protocolbuffers/protobuf/releases/download/$(PROTOC_VERSION)/$(PROTOC_ARCHIVE)
$(PROTOC): | $(LOCALBIN)
	curl -L $(PROTOC_URL) -o $(LOCALBIN)/$(PROTOC_ARCHIVE)
	unzip $(LOCALBIN)/$(PROTOC_ARCHIVE) -d $(LOCALBIN)/protoc_$(PROTOC_VERSION)
	mv $(LOCALBIN)/protoc_$(PROTOC_VERSION)/bin/protoc $@
	rm -f $(LOCALBIN)/$(PROTOC_ARCHIVE)
	rm -rf $(LOCALBIN)/protoc_$(PROTOC_VERSION)

.PHONY: protoc-gen-go
protoc-gen-go: $(PROTOC_GEN_GO)
$(PROTOC_GEN_GO): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(GO) install \
		  google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC)
$(PROTOC_GEN_GO_GRPC): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(GO) install \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

.PHONY: mockgen
mockgen: $(MOCKGEN)
$(MOCKGEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) $(GO) install \
		go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

clean:
	rm -rf $(PROTOC_OUT) $(LOCALBIN)
