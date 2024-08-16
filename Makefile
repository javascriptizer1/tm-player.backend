include .env

LOCAL_BIN := $(CURDIR)/bin

PROTOC := protoc
PROTOC_GEN_GO := $(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(LOCAL_BIN)/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(LOCAL_BIN)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(LOCAL_BIN)/protoc-gen-openapiv2
GOOSE := $(LOCAL_BIN)/goose

PROTO_SRC_DIR := api/proto
PROTO_VENDOR_DIR := vendor.protogen
MIGRATIONS_SRC_DIR := database/postgres/migrations

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

vendor:
	@if [ ! -d vendor.protogen/buf/validate ]; then \
		mkdir -p vendor.protogen/buf/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/buf/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

generate-api:
	make generate-api-city
	make generate-api-position
	make generate-api-player

generate-api-city:
	$(PROTOC) --proto_path=$(PROTO_SRC_DIR) --proto_path=$(PROTO_VENDOR_DIR) \
		--go_out=./pkg/gengrpc --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=./pkg/gengrpc --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--grpc-gateway_out=./pkg/gengrpc --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GRPC_GATEWAY) \
		$(PROTO_SRC_DIR)/city.proto

generate-api-position:
	$(PROTOC) --proto_path=$(PROTO_SRC_DIR) --proto_path=$(PROTO_VENDOR_DIR) \
		--go_out=./pkg/gengrpc --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=./pkg/gengrpc --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--grpc-gateway_out=./pkg/gengrpc --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GRPC_GATEWAY) \
		$(PROTO_SRC_DIR)/position.proto

generate-api-player:
	$(PROTOC) --proto_path=$(PROTO_SRC_DIR) --proto_path=$(PROTO_VENDOR_DIR) \
		--go_out=./pkg/gengrpc --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--go-grpc_out=./pkg/gengrpc --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--grpc-gateway_out=./pkg/gengrpc --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(PROTOC_GEN_GRPC_GATEWAY) \
		$(PROTO_SRC_DIR)/player.proto

run:
	go run ./cmd/tm-player/main.go

lint:
	docker run --rm -v .:/src -w /src golangci/golangci-lint:v1.59.1 golangci-lint --config .golangci.pipeline.yaml run

gensqlc:
	docker run --rm -v .:/src -w /src sqlc/sqlc:1.26.0 -f database/sqlc.yaml generate

gengoose:
	$(GOOSE) -dir $(MIGRATIONS_SRC_DIR) create ${name} sql

upgoose:
	GOOSE_DBSTRING=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable $(GOOSE) -dir $(MIGRATIONS_SRC_DIR) postgres up
