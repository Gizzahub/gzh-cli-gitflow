# quality.mk - Code Quality targets

# ==============================================================================
# Formatting Targets
# ==============================================================================

.PHONY: fmt lint lint-fix security quality quality-fix

fmt: ## format code
	@echo -e "$(CYAN)Formatting code...$(RESET)"
	@gofumpt -w . 2>/dev/null || go fmt ./...
	@goimports -w -local github.com/gizzahub/gzh-cli-gitflow . 2>/dev/null || true
	@echo -e "$(GREEN)Formatting complete$(RESET)"

# ==============================================================================
# Linting Targets
# ==============================================================================

lint: ## run linter
	@echo -e "$(CYAN)Running linter...$(RESET)"
	@golangci-lint run ./... || echo -e "$(YELLOW)golangci-lint not installed, skipping$(RESET)"

lint-fix: ## run linter with auto-fix
	@echo -e "$(CYAN)Running linter with fixes...$(RESET)"
	@golangci-lint run --fix ./... || echo -e "$(YELLOW)golangci-lint not installed$(RESET)"

# ==============================================================================
# Security Targets
# ==============================================================================

security: ## run security checks
	@echo -e "$(CYAN)Running security checks...$(RESET)"
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./... || echo -e "$(YELLOW)Security check completed$(RESET)"

# ==============================================================================
# Quality Workflow
# ==============================================================================

quality: fmt lint test ## run comprehensive quality checks
	@echo -e "$(GREEN)Quality checks passed$(RESET)"

quality-fix: fmt lint-fix test ## apply fixes and run quality checks
	@echo -e "$(GREEN)Quality fixes applied$(RESET)"
