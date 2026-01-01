# gzh-cli-gitflow

Git-flow workflow automation CLI tool.

## Features

- **Classic Git-flow** - 5 branch types (master, develop, feature, release, hotfix)
- **Safe Operations** - Input validation, merge conflict detection
- **Flexible Config** - Global + per-project configuration
- **Cross-platform** - Linux, macOS, Windows (amd64, arm64)

## Installation

```bash
# From source
go install github.com/gizzahub/gzh-cli-gitflow/cmd/gz-flow@latest

# Or build locally
make build && make install
```

## Quick Start

```bash
# Initialize git-flow in a repository
cd your-project
gz-flow init

# Start a feature
gz-flow feature start user-authentication

# Finish the feature (merges to develop)
gz-flow feature finish user-authentication

# Start a release
gz-flow release start 1.0.0

# Finish the release (merges to main + develop, creates tag)
gz-flow release finish 1.0.0
```

## Commands

| Command | Description |
|---------|-------------|
| `gz-flow init` | Initialize git-flow in repository |
| `gz-flow feature start <name>` | Create feature branch from develop |
| `gz-flow feature finish <name>` | Merge feature to develop |
| `gz-flow release start <version>` | Create release branch from develop |
| `gz-flow release finish <version>` | Merge release, create tag |
| `gz-flow hotfix start <version>` | Create hotfix from master |
| `gz-flow hotfix finish <version>` | Merge hotfix to main + develop |
| `gz-flow status` | Show current workflow state |
| `gz-flow list [type]` | List active flow branches |
| `gz-flow config [key] [value]` | Manage configuration |

## Configuration

### Global Config (`~/.gz/gitflow`)

```yaml
branches:
  master: master
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

### Project Config (`.gzflow.yaml`)

Project-level config overrides global settings:

```yaml
branches:
  master: main  # This project uses 'main'
```

## Development

```bash
# Setup
make install-tools

# Build
make build

# Test
make test

# Quality check (required before commit)
make quality

# Quick development cycle
make dev-fast
```

## Documentation

- [PRODUCT.md](PRODUCT.md) - Product goals and constraints
- [REQUIREMENTS.md](REQUIREMENTS.md) - Technical requirements
- [docs/00-product/](docs/00-product/) - Product strategy documents

## Related Projects

- [gzh-cli](https://github.com/gizzahub/gzh-cli) - Main CLI tool
- [gzh-cli-gitforge](https://github.com/gizzahub/gzh-cli-gitforge) - Git forge operations
- [gzh-cli-core](https://github.com/gizzahub/gzh-cli-core) - Shared utilities

## License

MIT License - see [LICENSE](LICENSE) for details.
