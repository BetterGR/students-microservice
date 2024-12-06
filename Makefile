# Makefile for generating Go code from .proto files

# Directories
PROTO_DIR = students_protobuf
OUTPUT_DIR = ./students_protobuf

# Proto file
PROTO_FILE = $(PROTO_DIR)/students_service.proto

# Go-related flags
GO_OUT_FLAGS = --go_out=paths=source_relative:$(OUTPUT_DIR)
GO_GRPC_FLAGS = --go-grpc_out=paths=source_relative:$(OUTPUT_DIR)

# Generated Go files (adjust as needed)
GENERATED_GO_FILES = $(OUTPUT_DIR)/students_service.pb.go $(OUTPUT_DIR)/students_service_grpc.pb.go

# Default target
all: generate

# Generate Go code from proto files
generate: $(GENERATED_GO_FILES)
	@echo - Done.

# Check if any of the generated files are older than the .proto file
$(GENERATED_GO_FILES): $(PROTO_FILE)
	@echo - Generating Go code from proto file...
	protoc -I $(PROTO_DIR) $(PROTO_FILE) $(GO_OUT_FLAGS) $(GO_GRPC_FLAGS)
	@echo - Go code generated.

# If the generated files are up to date, print "No change needed"
.PHONY: all generate
