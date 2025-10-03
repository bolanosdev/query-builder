.PHONY: test test-unit test-integration test-verbose test-coverage clean

test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	@go test -v ./tests/unit/

test-integration:
	@echo "Running integration tests..."
	@go test -v ./tests/integration/

test-verbose:
	@echo "Running all tests with verbose output..."
	@go test -v ./tests/...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./tests/...

clean:
	@echo "Cleaning test cache..."
	@go clean -testcache
