PROTO_DIR = ../api/auth-service
OUT_DIR = ../auth-service/generated/
PROTOC = protoc
PROTOC_GEN_GO = $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC = $(shell which protoc-gen-go-grpc)

PROTO_FILES = auth-service.proto

MIGRATION_DIR = ./db/migrations
DATABASE_URL = postgres://postgres:postgres@localhost:5432/public?sslmode=disable

.PHONY: proto
proto:
	@if [ -z "$(PROTOC_GEN_GO)" ] || [ -z "$(PROTOC_GEN_GO_GRPC)" ]; then \
		echo "Error: protoc-gen-go or protoc-gen-go-grpc is not installed"; \
		exit 1; \
	fi
	$(PROTOC) -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go-grpc_out=$(OUT_DIR) \
		$(PROTO_DIR)/$(PROTO_FILES)

.PHONY: migrate-up migrate-down create-migration status

create-migration:
	goose -dir $(MIGRATION_DIR) create $(NAME) sql

migrate-up:
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" down

status:
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" status