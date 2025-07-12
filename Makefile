PROJECT_NAME := aclapi
PROTO_DIR     := internal/grpcserver/protos
GO_PACKAGES   := ./cmd/... ./internal/...
BUILD_DIR     := bin
TAR_BUILD_DIR := build

API_BIN       := $(BUILD_DIR)/$(PROJECT_NAME)

PROTOC         := protoc
PROTOC_GEN_GO  := protoc-gen-go
PROTOC_GEN_GRPC:= protoc-gen-go-grpc

TARGETS := \
	linux_amd64 \
	linux_arm64

.PHONY: all
all: fmt proto vet build

.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)... (online)"
	@mkdir -p $(BUILD_DIR)
	go build -o $(API_BIN) ./cmd/$(PROJECT_NAME)

.PHONY: build-offline
build-offline:
	@echo "Building $(PROJECT_NAME)... (offline using vendor)"
	@mkdir -p $(BUILD_DIR)
	go build -mod=vendor -o $(API_BIN) ./cmd/$(PROJECT_NAME)

.PHONY: build-cross
build-cross:
	@echo "Cross building for: $(TARGETS) (online)"
	@mkdir -p $(BUILD_DIR)
	@for target in $(TARGETS); do \
		OS=$${target%_*}; \
		ARCH=$${target#*_}; \
		OUT=$(BUILD_DIR)/$(PROJECT_NAME)-$$OS-$$ARCH; \
		echo "Building $$OUT..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $$OUT ./cmd/$(PROJECT_NAME); \
	done

.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	$(PROTOC) \
	  --proto_path=. \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  $(PROTO_DIR)/*.proto

.PHONY: vendor
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

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
	@echo "Installing binary to /usr/local/bin..."
	install -m 755 $(API_BIN) /usr/local/bin/$(PROJECT_NAME)

.PHONY: docker-api
docker-api:
	docker build -f Dockerfile.api -t $(PROJECT_NAME):latest .

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	find . -type f -name '*.pb.go' -delete
	rm -rf vendor

.PHONY: run-api
run-api: $(API_BIN)
	@echo "Running $(PROJECT_NAME) as user 'aclapi'..."
	su -s /bin/bash aclapi -c "$(API_BIN)"

## Package project with vendor in tar.gz for offline install
.PHONY: package
package: vendor
	@echo "Packaging project source with vendor..."
	@mkdir -p $(TAR_BUILD_DIR)
	@TMP_DIR=$$(mktemp -d); \
	NAME=$(PROJECT_NAME)-source; \
	echo "Copying files to $$TMP_DIR/$$NAME..."; \
	mkdir -p $$TMP_DIR/$$NAME; \
	rsync -a --exclude '$(TAR_BUILD_DIR)' --exclude '$(API_BIN)' --exclude '*.tar.gz' ./ $$TMP_DIR/$$NAME; \
	TARBALL=$(PROJECT_NAME)-source.tar.gz; \
	tar -czf $(TAR_BUILD_DIR)/$$TARBALL -C $$TMP_DIR $$NAME; \
	echo "Created $(TAR_BUILD_DIR)/$$TARBALL"; \
	rm -rf $$TMP_DIR

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  all           : fmt, vet, proto, build"
	@echo "  build         : build for local OS/arch (online)"
	@echo "  build-offline : build for local OS/arch using vendor (offline)"
	@echo "  build-cross   : cross build for $(TARGETS) (online)"
	@echo "  proto         : generate protobuf code"
	@echo "  vendor        : vendor dependencies (must be online)"
	@echo "  package       : create tarball of project + vendor"
	@echo "  fmt           : format Go code"
	@echo "  vet           : run go vet"
	@echo "  lint          : run golangci-lint"
	@echo "  test          : run tests"
	@echo "  install       : install binary to /usr/local/bin"
	@echo "  docker-api    : build Docker image for API"
	@echo "  clean         : clean build artifacts and vendor"
	@echo "  run-api       : run API daemon as 'aclapi' user"
	@echo "  help          : show this help"
