# Vision

## Why gzh-cli-gitflow Exists

Git-flow is a proven branching model, but existing tools have limitations:

1. **git-flow (AVH Edition)** - Bash-based, no cross-platform consistency
2. **git-flow-next** - Separate binary, not integrated with existing toolchain
3. **Manual workflow** - Error-prone, inconsistent branch naming

gzh-cli-gitflow provides a **Go-native, cross-platform Git-flow implementation** that integrates seamlessly with the gzh-cli ecosystem.

______________________________________________________________________

## Vision Statement

> Make Git-flow branching effortless, consistent, and safe across all platforms and team sizes.

______________________________________________________________________

## Why Now

**Technical Readiness:**

- Go ecosystem mature for CLI development
- gzh-cli-core provides shared utilities
- gzh-cli-gitforge establishes Git safety patterns

**Market Need:**

- Remote/hybrid teams need consistent workflows
- CI/CD pipelines require predictable branch patterns
- Version management complexity increasing

**Strategic Fit:**

- Completes gzh-cli Git tooling story
- Reuses gitforge safety patterns
- Shared codebase reduces maintenance

______________________________________________________________________

## Anti-Goals

Things we explicitly will NOT do:

| Anti-Goal                  | Reason                                    |
| -------------------------- | ----------------------------------------- |
| Replace git commands       | Git CLI is the source of truth            |
| Auto-merge conflicts       | Human judgment required for conflicts     |
| IDE integration            | Focus on CLI, IDE plugins are separate    |
| Custom workflow DSL        | YAML config is sufficient (v0.x)          |
| Repository hosting         | Out of scope (use gitforge for forge ops) |

______________________________________________________________________

## Success Vision

**For Individual Developers:**

- Branch creation in one command, not four
- Consistent naming without memorizing patterns
- Safe merges with automatic validation

**For Teams:**

- Everyone follows the same workflow
- Branch naming disputes eliminated
- Release process standardized

**For the gzh-cli Ecosystem:**

- Complete Git workflow solution
- Reusable library for other tools
- Consistent UX across all gz-* commands
