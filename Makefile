PROTO_DIR=proto
THIRD_PARTY_DIR=third_party/googleapis

PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

generate:
	protoc -I $(PROTO_DIR) \
		-I $(THIRD_PARTY_DIR) \
		--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PROTO_DIR) --grpc-gateway_opt=paths=source_relative \
		$(PROTO_FILES)

.PHONY: generate