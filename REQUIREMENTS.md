# Technical Requirements Document

**Project**: gzh-cli-gitflow
**Version**: 1.0
**Last Updated**: 2026-01-01
**Status**: Draft

______________________________________________________________________

## 1. System Requirements

### 1.1 Runtime Requirements

**Minimum System Specifications:**

| Component        | Requirement           | Notes                               |
| ---------------- | --------------------- | ----------------------------------- |
| Operating System | Linux, macOS, Windows | Kernel 4.x+, macOS 11+, Windows 10+ |
| Architecture     | amd64, arm64          | Native binaries for each platform   |
| Memory           | 128MB RAM             | For typical operations              |
| Disk Space       | 20MB                  | Binary + config                     |
| Git              | 2.30+                 | System Git CLI required             |

**Recommended Specifications:**

| Component | Recommendation | Benefit                      |
| --------- | -------------- | ---------------------------- |
| Memory    | 256MB+ RAM     | Better performance           |
| Git       | 2.40+          | Latest features              |

### 1.2 Development Requirements

**Build Environment:**

```yaml
Go: 1.23.0+
Make: GNU Make 4.0+
Git: 2.30+
golangci-lint: 1.55+
```

**Development Tools:**

```yaml
Required:
  - go: Go compiler and toolchain
  - git: Version control
  - make: Build automation

Optional:
  - gomock: Mock generation
  - gotestfmt: Test output formatting
  - pre-commit: Git hooks framework
```

______________________________________________________________________

## 2. Functional Requirements

### 2.1 Initialization (F1)

#### F1.1 Git-flow Init

**REQ-F1.1.1**: Repository Initialization

- MUST detect if repository is already initialized for git-flow
- MUST create develop branch if not exists
- MUST set default branch names in configuration
- MUST support re-initialization with `--force`

**REQ-F1.1.2**: Branch Detection

- MUST auto-detect main branch (master/main)
- MUST prompt for branch names if interactive mode
- MUST use defaults in non-interactive mode

**REQ-F1.1.3**: Configuration Storage

- MUST save configuration to `.gzflow.yaml`
- MUST merge with global config `~/.gz/gitflow`
- MUST validate configuration on save

### 2.2 Feature Branch Management (F2)

#### F2.1 Feature Start

**REQ-F2.1.1**: Branch Creation

- MUST create feature branch from develop
- MUST apply configured prefix (default: `feature/`)
- MUST validate branch name (no special characters)
- MUST check if branch already exists

**REQ-F2.1.2**: Branch Naming

- MUST support custom name: `gz-flow feature start <name>`
- MUST sanitize name (lowercase, hyphens)
- MUST reject reserved names (master, main, develop)

#### F2.2 Feature Finish

**REQ-F2.2.1**: Merge to Develop

- MUST merge feature branch into develop
- MUST use `--no-ff` merge by default
- MUST support `--squash` option
- MUST check for uncommitted changes before merge

**REQ-F2.2.2**: Cleanup

- MUST delete local feature branch after merge (configurable)
- MUST support `--keep` to preserve branch
- MUST NOT delete remote branch by default

**REQ-F2.2.3**: Conflict Handling

- MUST detect merge conflicts before attempting
- MUST provide clear error message on conflict
- MUST NOT auto-resolve conflicts

### 2.3 Release Branch Management (F3)

#### F3.1 Release Start

**REQ-F3.1.1**: Branch Creation

- MUST create release branch from develop
- MUST apply configured prefix (default: `release/`)
- MUST accept version number as argument
- MUST validate version format (semver recommended)

**REQ-F3.1.2**: Version Validation

- MUST check version is greater than last tag
- SHOULD support version bump helpers (major/minor/patch)
- MUST reject duplicate version numbers

#### F3.2 Release Finish

**REQ-F3.2.1**: Dual Merge

- MUST merge release branch into master/main
- MUST merge release branch into develop
- MUST handle merge order correctly

**REQ-F3.2.2**: Tagging

- MUST create annotated tag on master/main
- MUST use configured tag format (default: `v{version}`)
- MUST support custom tag message
- MUST support `--no-tag` option

**REQ-F3.2.3**: Cleanup

- MUST delete release branch after successful merge
- MUST NOT delete if either merge fails

### 2.4 Hotfix Branch Management (F4)

#### F4.1 Hotfix Start

**REQ-F4.1.1**: Branch Creation

- MUST create hotfix branch from master/main
- MUST apply configured prefix (default: `hotfix/`)
- MUST accept version number as argument

**REQ-F4.1.2**: Emergency Context

- SHOULD allow creation even with uncommitted changes
- MUST warn about uncommitted changes

#### F4.2 Hotfix Finish

**REQ-F4.2.1**: Dual Merge

- MUST merge hotfix branch into master/main
- MUST merge hotfix branch into develop
- MUST handle active release branch (merge into release instead)

**REQ-F4.2.2**: Tagging

- MUST create annotated tag on master/main
- MUST use same format as release tags

### 2.5 Status and Information (F5)

#### F5.1 Status Command

**REQ-F5.1.1**: Current State

- MUST show current branch type (feature/release/hotfix/other)
- MUST show base branch for current flow branch
- MUST indicate if working directory is clean

**REQ-F5.1.2**: Configuration Display

- MUST show active configuration
- MUST indicate source (global/local)

#### F5.2 List Command

**REQ-F5.2.1**: Branch Listing

- MUST list all active flow branches by type
- MUST support filtering by type: `gz-flow list feature`
- MUST show branch age and last commit

#### F5.3 Config Command

**REQ-F5.3.1**: Configuration Management

- MUST show current config: `gz-flow config`
- MUST get single value: `gz-flow config branches.master`
- MUST set value: `gz-flow config branches.master main`
- MUST support `--global` flag for global config

______________________________________________________________________

## 3. Non-Functional Requirements

### 3.1 Performance (NFR-P)

**REQ-NFR-P1**: Operation Speed

- Flow operations MUST complete in < 200ms (p95)
- All operations MUST complete in < 1000ms (p99)
- Git command overhead MUST be < 50ms

**REQ-NFR-P2**: Resource Usage

- Memory usage MUST stay under 128MB for typical operations
- No background processes or daemons

### 3.2 Reliability (NFR-R)

**REQ-NFR-R1**: Atomic Operations

- Finish operations MUST be atomic (all or nothing)
- Failed merges MUST leave repository in original state
- Configuration writes MUST be atomic

**REQ-NFR-R2**: Error Recovery

- MUST provide clear recovery instructions on failure
- MUST NOT corrupt repository state

### 3.3 Usability (NFR-U)

**REQ-NFR-U1**: Command Line Interface

- MUST provide `--help` for all commands
- MUST provide meaningful error messages
- MUST support `--dry-run` for destructive operations

**REQ-NFR-U2**: Migration Support

- MUST detect existing git-flow configuration
- SHOULD offer migration path from git-flow
- MUST be compatible with git-flow branch naming

### 3.4 Security (NFR-S)

**REQ-NFR-S1**: Input Validation

- ALL user inputs MUST be sanitized
- Branch names MUST be validated against injection patterns
- MUST reject paths with `..` or absolute paths

**REQ-NFR-S2**: Safe Execution

- MUST use `exec.Command` with separate arguments (no shell)
- MUST NOT pass unsanitized input to git commands
- MUST NOT log sensitive information

______________________________________________________________________

## 4. Technical Constraints

### 4.1 Technology Stack

| Component  | Choice       | Rationale                  |
| ---------- | ------------ | -------------------------- |
| Language   | Go 1.23+     | Consistency with gzh-cli   |
| CLI        | Cobra        | Standard for Go CLIs       |
| Config     | Viper + YAML | Flexible configuration     |
| Git        | System Git   | No library dependency      |
| Core Utils | gzh-cli-core | Shared utilities           |

### 4.2 Architecture Constraints

**Library-First Design:**

```
cmd/gz-flow/     # Thin CLI layer
├── main.go
└── cmd/
    ├── root.go
    ├── init.go
    ├── feature.go
    └── ...

pkg/             # Reusable library
├── flow/        # Core flow operations
├── branch/      # Branch utilities
└── config/      # Configuration management

internal/        # Private implementation
├── gitcmd/      # Safe git execution
└── validator/   # Input validation
```

**Dependency Rules:**

- `pkg/` MUST NOT import from `cmd/`
- `pkg/` MUST NOT import CLI frameworks
- `internal/` is private to this module
- Use `gzh-cli-core` for common utilities

### 4.3 Configuration Schema

```yaml
# ~/.gz/gitflow (global) or .gzflow.yaml (local)
branches:
  master: master      # or "main"
  develop: develop

prefixes:
  feature: feature/
  release: release/
  hotfix: hotfix/
  support: support/   # optional

options:
  delete_branch_after_finish: true
  push_after_finish: false
  tag_format: "v%s"   # %s = version
  require_clean_tree: true
```

______________________________________________________________________

## 5. Testing Requirements

### 5.1 Unit Tests

**Coverage Targets:**

| Package    | Target | Notes                    |
| ---------- | ------ | ------------------------ |
| pkg/flow   | >= 85% | Core logic               |
| pkg/branch | >= 85% | Branch operations        |
| pkg/config | >= 85% | Configuration management |
| internal/* | >= 80% | Internal utilities       |
| cmd/*      | >= 70% | CLI layer                |

**Test Patterns:**

- Table-driven tests for all functions
- Mock git commands for unit tests
- Test both success and error paths

### 5.2 Integration Tests

**Requirements:**

- Use real Git repositories (temp directories)
- Test complete workflows (init → start → finish)
- Test error scenarios (conflicts, dirty tree)
- Test configuration loading/merging

### 5.3 E2E Tests

**Scenarios:**

| Scenario           | Description                        |
| ------------------ | ---------------------------------- |
| Full Feature Flow  | init → feature start → finish      |
| Full Release Flow  | init → release start → finish      |
| Hotfix Flow        | init → hotfix start → finish       |
| Config Migration   | git-flow config → gzflow migration |
| Error Recovery     | Conflict detection and messaging   |

______________________________________________________________________

## 6. Deployment Requirements

### 6.1 Build

- Cross-compile for linux/darwin/windows (amd64/arm64)
- Use goreleaser for release builds
- Include version info in binary

### 6.2 Distribution

- GitHub Releases with checksums
- Homebrew formula (future)
- go install support

______________________________________________________________________

## 7. Documentation Requirements

### 7.1 User Documentation

- README with quick start
- Command reference (all commands and flags)
- Configuration reference
- Migration guide from git-flow

### 7.2 Developer Documentation

- CLAUDE.md for AI assistance
- CONTRIBUTING.md for contributors
- ARCHITECTURE.md for design overview
- GoDoc for public APIs

______________________________________________________________________

## 8. Acceptance Criteria

### 8.1 v0.1.0 Release Criteria

| Criteria             | Requirement                       |
| -------------------- | --------------------------------- |
| Core Commands        | init, feature, release, hotfix    |
| Utility Commands     | status, list, config              |
| Configuration        | Global + local config working     |
| Testing              | All coverage targets met          |
| Documentation        | README, CLI help complete         |
| CI/CD                | Build and test pipeline working   |

### 8.2 Quality Gates

- All tests passing
- No critical/high security issues
- Documentation complete
- Performance targets met

______________________________________________________________________

## 9. Traceability Matrix

| Requirement | Feature              | Test             | Status |
| ----------- | -------------------- | ---------------- | ------ |
| REQ-F1.1.1  | gz-flow init         | test_init.go     | -      |
| REQ-F2.1.1  | gz-flow feature      | test_feature.go  | -      |
| REQ-F2.2.1  | gz-flow feature      | test_feature.go  | -      |
| REQ-F3.1.1  | gz-flow release      | test_release.go  | -      |
| REQ-F3.2.1  | gz-flow release      | test_release.go  | -      |
| REQ-F4.1.1  | gz-flow hotfix       | test_hotfix.go   | -      |
| REQ-F4.2.1  | gz-flow hotfix       | test_hotfix.go   | -      |
| REQ-F5.1.1  | gz-flow status       | test_status.go   | -      |
| REQ-F5.2.1  | gz-flow list         | test_list.go     | -      |
| REQ-F5.3.1  | gz-flow config       | test_config.go   | -      |
| REQ-NFR-P1  | All commands         | benchmark_test   | -      |
| REQ-NFR-S1  | Input handling       | security_test    | -      |

______________________________________________________________________

## 10. Revision History

| Version | Date       | Author | Changes         |
| ------- | ---------- | ------ | --------------- |
| 1.0     | 2026-01-01 | Claude | Initial draft   |

______________________________________________________________________

**End of Document**
