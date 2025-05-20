PROJECT_NAME := acl-daemon
PROTO_DIR     := proto
GO_PACKAGES   := ./cmd/... ./internal/...
BUILD_DIR     := bin

CORE_BIN      := $(BUILD_DIR)/aclcore
API_BIN       := $(BUILD_DIR)/aclapi

PROTOC         := protoc
PROTOC_GEN_GO  := protoc-gen-go
PROTOC_GEN_GRPC:= protoc-gen-go-grpc

.PHONY: all
all: fmt vet proto build

.PHONY: build
build: $(CORE_BIN) $(API_BIN)

$(CORE_BIN):
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/aclcore

$(API_BIN):
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/aclapi

.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	$(PROTOC) \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	  --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	  $(PROTO_DIR)/*.proto

.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt $(GO_PACKAGES)

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet $(GO_PACKAGES)

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: install
install: build
	@echo "Installing binaries to /usr/local/bin..."
	install -m 755 $(CORE_BIN) /usr/local/bin/aclcore
	install -m 755 $(API_BIN)  /usr/local/bin/aclapi

.PHONY: docker-core docker-api
docker-core:
	docker build -f Dockerfile.core -t aclcore:latest .

docker-api:
	docker build -f Dockerfile.api  -t aclapi:latest .

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	find . -type f -name '*.pb.go' -delete

.PHONY: run-core run-api
run-core: $(CORE_BIN)
	@echo "Running aclcore (root)..."
	sudo $(CORE_BIN)

run-api: $(API_BIN)
	@echo "Running aclapi (as aclapi user)..."
	su -s /bin/bash aclapi -c "$(API_BIN)"

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  all         : fmt, vet, proto, build"
	@echo "  build       : build both binaries"
	@echo "  proto       : generate protobuf code"
	@echo "  fmt         : format Go code"
	@echo "  vet         : run go vet"
	@echo "  lint        : run golangci-lint"
	@echo "  test        : run tests"
	@echo "  install     : install binaries to /usr/local/bin"
	@echo "  docker-core : build Docker image for core daemon"
	@echo "  docker-api  : build Docker image for API daemon"
	@echo "  clean       : remove build artifacts"
	@echo "  run-core    : run core daemon as root"
	@echo "  run-api     : run API daemon as 'aclapi' user"
	@echo "  help        : this help message"
