# CLAUDE.md

LLM-optimized guidance for Claude Code when working with gzh-cli-gitflow.

______________________________________________________________________

## Quick Start (30s scan)

**Binary**: `gz-flow`
**Module**: `github.com/gizzahub/gzh-cli-gitflow`
**Go Version**: 1.23+
**Purpose**: Git-flow workflow automation CLI

______________________________________________________________________

## Top 10 Commands

| Command              | Purpose             | When to Use           |
| -------------------- | ------------------- | --------------------- |
| `make quality`       | fmt + lint + test   | Pre-commit (CRITICAL) |
| `make dev-fast`      | format + unit tests | Quick dev cycle       |
| `make build`         | Build binary        | After changes         |
| `make test`          | All tests           | Validation            |
| `make cover-html`    | Coverage report     | Check coverage        |
| `make fmt`           | Format code         | Fix formatting        |
| `make lint`          | Run linters         | Fix lint issues       |
| `make pr-check`      | Pre-PR verification | Before PR             |
| `make install`       | Install binary      | Local testing         |
| `make clean`         | Clean artifacts     | Fresh start           |

______________________________________________________________________

## Absolute Rules (DO/DON'T)

### DO

- Use `gzh-cli-core` for common utilities
- Run `make quality` before every commit
- **ALWAYS sanitize git inputs** (prevent command injection)
- Test coverage: 80%+ for core logic
- Follow git-flow command structure for familiarity

### DON'T

- Use shell execution (`sh -c`) - command injection risk
- Concatenate user input into commands
- Skip input validation
- Log credentials or sensitive data
- Break git-flow compatibility without approval

______________________________________________________________________

## Directory Structure

```
.
├── cmd/gz-flow/            # CLI commands
│   ├── main.go             # Entry point
│   └── cmd/                # Subcommands
│       ├── root.go
│       ├── init.go
│       ├── feature.go
│       ├── release.go
│       ├── hotfix.go
│       ├── status.go
│       ├── list.go
│       └── config.go
├── pkg/                    # Public library
│   ├── flow/               # Core flow operations
│   ├── branch/             # Branch management
│   └── config/             # Configuration
├── internal/               # Private packages
│   ├── gitcmd/             # Safe git execution
│   └── validator/          # Input validation
├── docs/                   # Documentation
│   └── 00-product/         # Product strategy
└── tests/                  # Test suites
    ├── integration/
    └── e2e/
```

______________________________________________________________________

## Configuration

**Global**: `~/.gz/gitflow` (YAML)
**Local**: `.gzflow.yaml` (project root, overrides global)

```yaml
branches:
  master: master      # or "main"
  develop: develop

prefixes:
  feature: feature/
  release: release/
  hotfix: hotfix/

options:
  delete_branch_after_finish: true
  push_after_finish: false
  tag_format: "v%s"
```

______________________________________________________________________

## Command Structure

```bash
gz-flow init                     # Initialize git-flow
gz-flow feature start <name>     # Start feature branch
gz-flow feature finish <name>    # Finish feature (merge to develop)
gz-flow release start <version>  # Start release branch
gz-flow release finish <version> # Finish release (merge, tag)
gz-flow hotfix start <version>   # Start hotfix branch
gz-flow hotfix finish <version>  # Finish hotfix (merge to main+develop)
gz-flow status                   # Show current state
gz-flow list [type]              # List flow branches
gz-flow config [key] [value]     # Manage configuration
```

______________________________________________________________________

## Security (CRITICAL)

### Safe Command Execution

```go
// SAFE - Arguments passed separately
cmd := exec.Command("git", "checkout", "-b", branchName)

// DANGEROUS - Shell execution
cmd := exec.Command("sh", "-c", "git checkout -b " + branchName)
```

### Input Validation

```go
// Always validate branch names
if !isValidBranchName(name) {
    return errors.New("invalid branch name")
}
```

______________________________________________________________________

## Shared Library (gzh-cli-core)

```go
import (
    "github.com/gizzahub/gzh-cli-core/logger"
    "github.com/gizzahub/gzh-cli-core/errors"
)
```

______________________________________________________________________

## Git Commit Format

```
{type}({scope}): {description}

Model: claude-{model}
Co-Authored-By: Claude <noreply@anthropic.com>
```

**Types**: feat, fix, docs, refactor, test, chore
**Scope**: cmd, internal, pkg/flow, pkg/branch, pkg/config

______________________________________________________________________

## Context Documentation

| Document                               | Purpose                |
| -------------------------------------- | ---------------------- |
| [PRODUCT.md](PRODUCT.md)               | Goals and constraints  |
| [REQUIREMENTS.md](REQUIREMENTS.md)     | Technical requirements |
| [docs/00-product/](docs/00-product/)   | Product strategy       |

______________________________________________________________________

**Last Updated**: 2026-01-01
