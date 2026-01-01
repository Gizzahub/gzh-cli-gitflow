# build.mk - Build and Installation targets

# ==============================================================================
# Build Configuration
# ==============================================================================

BINEXT := $(shell go env GOEXE)
BINARY := $(executablename)$(BINEXT)
GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)

ifeq ($(strip $(GOBIN)),)
  BINDIR := $(GOPATH)/bin
else
  BINDIR := $(GOBIN)
endif

# ==============================================================================
# Build Targets
# ==============================================================================

.PHONY: build install run clean

build: ## build golang binary
	@printf "$(CYAN)Building %s...$(RESET)\n" "$(BINARY)"
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/gz-flow
	@printf "$(GREEN)Built %s successfully$(RESET)\n" "$(BINARY)"

install: build ## install golang binary
	@printf "$(CYAN)Installing $(BINARY) $(VERSION) to %s$(RESET)\n" "$(BINDIR)/$(BINARY)"
	@mkdir -p "$(BINDIR)"
	@mv $(BINARY) "$(BINDIR)"/
	@printf "$(GREEN)Installed $(BINARY) $(VERSION)$(RESET)\n"

run: ## run the application
	@go run -ldflags "-X main.version=$(VERSION)" ./cmd/gz-flow $(ARGS)

clean: ## clean up environment
	@echo -e "$(CYAN)Cleaning up...$(RESET)"
	@rm -rf coverage.out coverage.html dist/ $(executablename) $(BINARY)
	@echo -e "$(GREEN)Cleanup completed$(RESET)"
