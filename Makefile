.PHONY: test tidy fmt build

# Run all tests
test:
	go test -v ./...

# Run go mod tidy
tidy:
	go mod tidy

# Format code with goimports, gofumpt, and run go vet
fmt:
	goimports -w .
	gofumpt -w .
	go vet ./...

# Build the project
build:
	go build ./...

# Clean build cache
clean:
	go clean -cache

