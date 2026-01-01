# tools.mk - Tool Installation targets

# ==============================================================================
# Tool Installation
# ==============================================================================

.PHONY: install-tools tools-status

install-tools: ## install development tools
	@echo -e "$(CYAN)Installing development tools...$(RESET)"
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	@which gofumpt > /dev/null || go install mvdan.cc/gofumpt@latest
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo -e "$(GREEN)Tools installed$(RESET)"

tools-status: ## check installed tools
	@echo -e "$(CYAN)Checking installed tools...$(RESET)"
	@echo -n "goimports: " && which goimports || echo "not installed"
	@echo -n "gofumpt: " && which gofumpt || echo "not installed"
	@echo -n "golangci-lint: " && which golangci-lint || echo "not installed"
