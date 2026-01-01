# Product Goals (No-PRD)

**Project**: gzh-cli-gitflow
**Doc Type**: Goals + Constraints + Quality Gates
**Status**: Active
**Last Updated**: 2026-01-01

______________________________________________________________________

## 1) Product Intent

gzh-cli-gitflow provides a Git-flow workflow CLI tool that:

- automates Git-flow branching operations (feature, release, hotfix),
- enforces consistent branch naming conventions,
- simplifies release management with tagging,
- and offers flexible configuration (global + per-project).

This document replaces a full PRD. It defines goals, non-goals, guardrails,
and release quality gates.

**Detailed Documentation**: See [docs/00-product/](docs/00-product/) for comprehensive product documentation.

| Document                                             | Description                         |
| ---------------------------------------------------- | ----------------------------------- |
| [Vision](docs/00-product/01-vision.md)               | Why this project exists, anti-goals |
| [Principles](docs/00-product/02-principles.md)       | Core values and trade-offs          |
| [Problem Space](docs/00-product/03-problem-space.md) | Target users and pain points        |
| [Scope](docs/00-product/04-scope.md)                 | What's in/out of scope              |
| [Metrics](docs/00-product/05-metrics.md)             | Success criteria and measurement    |
| [Roadmap](docs/00-product/06-roadmap.md)             | Phases and milestones               |

______________________________________________________________________

## 2) Goals (Measurable Targets)

G1. **Reduce Git-flow operation time by 50%**

- Target: branch creation/merge operations p95 < 200ms
- Compared to manual git commands sequence

G2. **Zero merge conflicts from workflow errors**

- Target: 100% of finish operations validate merge state first
- Automatic detection of unmerged changes

G3. **Branch naming consistency**

- Target: >= 95% of branches follow configured prefix patterns
- Validation on branch creation

G4. **Seamless migration from git-flow**

- Target: existing git-flow users productive within 5 minutes
- Command structure familiar to git-flow users

G5. **Configuration flexibility**

- Target: 100% of settings overridable per-project
- Global defaults with local overrides

______________________________________________________________________

## 3) Non-Goals (Explicitly Out of Scope)

- No GUI or web interface
- No GitHub/GitLab/Bitbucket API integration (use gzh-cli-gitforge)
- No custom workflow definitions beyond Git-flow/GitHub-flow (v1.0+)
- No CI/CD integration (only CLI output for automation)
- No git hooks installation/management
- No automatic conflict resolution

______________________________________________________________________

## 4) Guardrails and Technical Constraints

**Architecture**

- Library-first: core logic lives in `pkg/*` and is reusable
- CLI in `cmd/` is thin; it should delegate to library packages
- Git CLI is the source of truth (no go-git as the primary engine)
- All operations accept `context.Context` for cancellation/timeouts

**Dependency Boundaries**

- `pkg/` should avoid CLI framework dependencies
- Use `gzh-cli-core` for common utilities (logging, errors, testutil)
- Minimal external dependencies

**Compatibility**

- Go 1.23+ (align with `go.mod`)
- Git 2.30+ on Linux/macOS/Windows (amd64/arm64)
- Compatible with existing `.gitflow` configurations (migration support)

**Safety**

- Finish operations require clean working directory or explicit flag
- Force merge is blocked unless explicitly overridden
- Inputs must be sanitized before Git CLI execution
- Branch names validated against injection patterns

**Configuration**

- Global config: `~/.gz/gitflow` (YAML)
- Local config: `.gzflow.yaml` (project root)
- Local overrides global
- Environment variables for CI/CD contexts

______________________________________________________________________

## 5) Quality Gates (Release Readiness)

**Build and Lint**

- `make build` and `make quality` pass with no warnings

**Testing**

- Unit + integration + E2E test suites pass
- Coverage targets:
  - `internal/` >= 80%
  - `pkg/` >= 85%
  - `cmd/` >= 70% (or equivalent CLI integration coverage)

**Performance**

- 95% of flow operations < 200ms
- 100% of flow operations < 1000ms

**Docs and Examples**

- CLI reference complete for all commands
- Configuration reference with examples
- Migration guide from git-flow

**Integration**

- gzh-cli integration complete and tested (v0.2.0+)

______________________________________________________________________

## 6) Decision Rules

- New features must map to at least one goal or be explicitly approved
- Anything that violates guardrails requires a documented exception
- Release is blocked if quality gates are not met
- git-flow compatibility takes priority over innovation (v0.x)

______________________________________________________________________

**End of Document**
