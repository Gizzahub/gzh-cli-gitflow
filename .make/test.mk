# test.mk - Testing targets

# ==============================================================================
# Testing Targets
# ==============================================================================

.PHONY: test test-unit test-integration test-e2e cover cover-html bench

test: ## run all tests with coverage
	@echo -e "$(CYAN)Running all tests...$(RESET)"
	go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	@echo -e "$(GREEN)Tests completed$(RESET)"

test-unit: ## run only unit tests
	@echo -e "$(CYAN)Running unit tests...$(RESET)"
	go test -short --cover -v ./...
	@echo -e "$(GREEN)Unit tests completed$(RESET)"

test-integration: build ## run integration tests
	@echo -e "$(CYAN)Running integration tests...$(RESET)"
	@if [ -d "./tests/integration" ]; then \
		go test -v ./tests/integration/...; \
	else \
		echo -e "$(YELLOW)No integration tests found$(RESET)"; \
	fi

test-e2e: build ## run e2e tests
	@echo -e "$(CYAN)Running E2E tests...$(RESET)"
	@if [ -d "./tests/e2e" ]; then \
		go test -v ./tests/e2e/...; \
	else \
		echo -e "$(YELLOW)No E2E tests found$(RESET)"; \
	fi

# ==============================================================================
# Coverage Targets
# ==============================================================================

cover: ## display test coverage
	@echo -e "$(CYAN)Generating coverage report...$(RESET)"
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

cover-html: ## generate HTML coverage report
	@echo -e "$(CYAN)Generating HTML coverage report...$(RESET)"
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo -e "$(GREEN)Coverage report: coverage.html$(RESET)"

# ==============================================================================
# Benchmark Targets
# ==============================================================================

bench: ## run benchmarks
	@echo -e "$(CYAN)Running benchmarks...$(RESET)"
	go test -bench=. -benchmem ./...
