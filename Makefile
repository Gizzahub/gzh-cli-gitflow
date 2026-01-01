# Makefile - gzh-cli-gitflow CLI Tool
# Git-flow workflow automation

# ==============================================================================
# Project Configuration
# ==============================================================================

projectname := gzh-cli-gitflow
executablename := gz-flow
VERSION ?= $(shell cat VERSION 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "0.1.0")

# Go configuration
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org

# Colors
export CYAN := \033[36m
export GREEN := \033[32m
export YELLOW := \033[33m
export RED := \033[31m
export BLUE := \033[34m
export MAGENTA := \033[35m
export RESET := \033[0m

# ==============================================================================
# Include Modular Makefiles
# ==============================================================================

include .make/build.mk
include .make/test.mk
include .make/quality.mk
include .make/deps.mk
include .make/dev.mk
include .make/tools.mk
include .make/docker.mk

# ==============================================================================
# Help System
# ==============================================================================

.DEFAULT_GOAL := help

.PHONY: help help-build help-test help-quality help-deps help-dev

help: ## show main help menu
	@echo -e "$(CYAN)"
	@echo "╔══════════════════════════════════════════════════════════════════════════════╗"
	@echo -e "║                        $(MAGENTA)gzh-cli-gitflow Makefile Help$(CYAN)                       ║"
	@echo -e "║                    $(YELLOW)Git-flow Workflow Automation CLI$(CYAN)                        ║"
	@echo "╚══════════════════════════════════════════════════════════════════════════════╝"
	@echo -e "$(RESET)"
	@echo -e "$(GREEN)Quick Commands:$(RESET)"
	@echo -e "  $(CYAN)make build$(RESET)         Build binary ($(executablename))"
	@echo -e "  $(CYAN)make test$(RESET)          Run tests with coverage"
	@echo -e "  $(CYAN)make quality$(RESET)       Format + lint + test"
	@echo -e "  $(CYAN)make dev-fast$(RESET)      Quick dev cycle (format + unit tests)"
	@echo -e "  $(CYAN)make install$(RESET)       Install to GOPATH/bin"
	@echo ""
	@echo -e "$(GREEN)Categories:$(RESET)"
	@echo -e "  $(YELLOW)make help-build$(RESET)    Build and installation"
	@echo -e "  $(YELLOW)make help-test$(RESET)     Testing and coverage"
	@echo -e "  $(YELLOW)make help-quality$(RESET)  Code quality"
	@echo -e "  $(YELLOW)make help-deps$(RESET)     Dependency management"
	@echo -e "  $(YELLOW)make help-dev$(RESET)      Development workflow"

help-build: ## show build help
	@echo -e "$(GREEN)Build Commands:$(RESET)"
	@echo -e "  $(CYAN)build$(RESET)              Build binary"
	@echo -e "  $(CYAN)install$(RESET)            Install to GOPATH/bin"
	@echo -e "  $(CYAN)run$(RESET)                Run the application"
	@echo -e "  $(CYAN)clean$(RESET)              Clean build artifacts"

help-test: ## show test help
	@echo -e "$(GREEN)Test Commands:$(RESET)"
	@echo -e "  $(CYAN)test$(RESET)               Run all tests with coverage"
	@echo -e "  $(CYAN)test-unit$(RESET)          Run unit tests only"
	@echo -e "  $(CYAN)cover-html$(RESET)         Generate HTML coverage report"
	@echo -e "  $(CYAN)bench$(RESET)              Run benchmarks"

help-quality: ## show quality help
	@echo -e "$(GREEN)Quality Commands:$(RESET)"
	@echo -e "  $(CYAN)fmt$(RESET)                Format code"
	@echo -e "  $(CYAN)lint$(RESET)               Run linter"
	@echo -e "  $(CYAN)quality$(RESET)            Format + lint + test"
	@echo -e "  $(CYAN)security$(RESET)           Run security checks"

help-deps: ## show deps help
	@echo -e "$(GREEN)Dependency Commands:$(RESET)"
	@echo -e "  $(CYAN)deps-check$(RESET)         Check for updates"
	@echo -e "  $(CYAN)deps-tidy$(RESET)          Run go mod tidy"
	@echo -e "  $(CYAN)deps-update$(RESET)        Update dependencies"

help-dev: ## show dev help
	@echo -e "$(GREEN)Development Commands:$(RESET)"
	@echo -e "  $(CYAN)dev$(RESET)                Standard dev workflow"
	@echo -e "  $(CYAN)dev-fast$(RESET)           Quick dev cycle"
	@echo -e "  $(CYAN)pr-check$(RESET)           Pre-PR verification"

# ==============================================================================
# Project Information
# ==============================================================================

.PHONY: info

info: ## show project information
	@echo -e "$(CYAN)"
	@echo "╔══════════════════════════════════════════════════════════════════════════════╗"
	@echo -e "║                         $(MAGENTA)gzh-cli-gitflow Info$(CYAN)                              ║"
	@echo "╚══════════════════════════════════════════════════════════════════════════════╝"
	@echo -e "$(RESET)"
	@echo -e "$(GREEN)Project:$(RESET)"
	@echo -e "  Name:           $(YELLOW)$(projectname)$(RESET)"
	@echo -e "  Executable:     $(YELLOW)$(executablename)$(RESET)"
	@echo -e "  Version:        $(YELLOW)$(VERSION)$(RESET)"
	@echo ""
	@echo -e "$(GREEN)Environment:$(RESET)"
	@echo "  Go Version:     $$(go version | cut -d' ' -f3)"
	@echo "  GOPATH:         $$(go env GOPATH)"
