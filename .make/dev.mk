# dev.mk - Development Workflow targets

# ==============================================================================
# Development Workflow
# ==============================================================================

.PHONY: dev dev-fast pr-check verify

dev: fmt lint test ## standard development workflow
	@echo -e "$(GREEN)Development workflow complete$(RESET)"

dev-fast: fmt test-unit ## quick development cycle
	@echo -e "$(GREEN)Quick dev cycle complete$(RESET)"

pr-check: quality ## pre-PR verification
	@echo -e "$(GREEN)PR check passed$(RESET)"

verify: quality deps-verify ## complete verification
	@echo -e "$(GREEN)Verification complete$(RESET)"
