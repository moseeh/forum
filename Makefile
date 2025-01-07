project = Forum

# Default target
run:
	@go run ./cmd/web/

test:
	@go test ./...
