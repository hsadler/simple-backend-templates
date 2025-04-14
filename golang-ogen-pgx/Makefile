.PHONY: generate build run test clean docker-build docker-run docker-compose-up docker-compose-down up down

BINARY_NAME=./build/app

# Main package path
MAIN_PACKAGE=./cmd/serverd

# Generate API code from OpenAPI spec
openapi-generate:
	docker run --rm --volume=".:/workspace" ghcr.io/ogen-go/ogen:latest \
		--target=/workspace/internal/openapi/ogen \
		--package=ogen \
		--clean \
		/workspace/openapi-schema.yaml \
		--generate-wrappers=on

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Start Docker Compose services
up:
	docker compose up -d

# Stop Docker Compose services
down:
	docker compose down
