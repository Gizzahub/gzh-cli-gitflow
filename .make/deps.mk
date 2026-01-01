# deps.mk - Dependency Management targets

# ==============================================================================
# Dependency Management
# ==============================================================================

.PHONY: deps-check deps-tidy deps-update deps-update-minor deps-verify deps-security

deps-check: ## check for outdated dependencies
	@echo -e "$(CYAN)Checking for outdated dependencies...$(RESET)"
	@go list -u -m all | grep '\[' || echo -e "$(GREEN)All dependencies up to date$(RESET)"

deps-tidy: ## run go mod tidy
	@echo -e "$(CYAN)Tidying Go modules...$(RESET)"
	@go mod tidy
	@echo -e "$(GREEN)Go modules tidied$(RESET)"

deps-update: ## update dependencies (patch versions)
	@echo -e "$(CYAN)Updating dependencies...$(RESET)"
	@go get -u=patch ./...
	@go mod tidy
	@echo -e "$(GREEN)Dependencies updated$(RESET)"

deps-update-minor: ## update to latest minor versions
	@echo -e "$(CYAN)Updating to minor versions...$(RESET)"
	@go get -u ./...
	@go mod tidy
	@echo -e "$(GREEN)Dependencies updated$(RESET)"

deps-verify: ## verify dependency checksums
	@echo -e "$(CYAN)Verifying dependencies...$(RESET)"
	@go mod verify
	@echo -e "$(GREEN)Dependencies verified$(RESET)"

deps-security: ## run security audit
	@echo -e "$(CYAN)Running security audit...$(RESET)"
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
