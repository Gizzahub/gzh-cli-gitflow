# Scope

## In Scope

### v0.1.0 (Core)

| Feature          | Description                           | Why Included                    |
| ---------------- | ------------------------------------- | ------------------------------- |
| init             | Initialize git-flow in repository     | Required for any git-flow usage |
| feature start    | Create feature branch from develop    | Core git-flow operation         |
| feature finish   | Merge feature to develop              | Core git-flow operation         |
| release start    | Create release branch from develop    | Core git-flow operation         |
| release finish   | Merge release, tag, cleanup           | Core git-flow operation         |
| hotfix start     | Create hotfix branch from master      | Core git-flow operation         |
| hotfix finish    | Merge hotfix to master and develop    | Core git-flow operation         |
| status           | Show current workflow state           | Essential for users             |
| list             | List active flow branches             | Essential for users             |
| config           | Manage configuration                  | Required for customization      |

### v0.2.0 (Remote)

| Feature          | Description                           | Why Included                    |
| ---------------- | ------------------------------------- | ------------------------------- |
| publish          | Push flow branch to remote            | Team collaboration              |
| pull             | Pull updates from remote              | Team collaboration              |
| track            | Set up remote tracking                | Team collaboration              |
| delete           | Delete local and remote branches      | Cleanup operations              |

### v0.3.0 (Integration)

| Feature          | Description                           | Why Included                    |
| ---------------- | ------------------------------------- | ------------------------------- |
| Library API      | Stable pkg/* for gzh-cli              | Ecosystem integration           |
| gzh-cli cmd      | `gz flow` in main CLI                 | Unified experience              |

______________________________________________________________________

## Out of Scope

### Explicitly Excluded

| Feature              | Why Excluded                                  |
| -------------------- | --------------------------------------------- |
| GUI/Web interface    | CLI-focused tool, GUIs are separate products  |
| Forge API calls      | Use gzh-cli-gitforge for GitHub/GitLab API    |
| CI/CD integration    | Scriptable CLI is sufficient                  |
| Git hooks manager    | Separate concern, use pre-commit or similar   |
| Automatic conflicts  | Human judgment required                       |
| IDE plugins          | Separate products per IDE                     |
| Repository hosting   | Out of scope entirely                         |

### Deferred (Maybe Later)

| Feature              | Why Deferred                                  | When                |
| -------------------- | --------------------------------------------- | ------------------- |
| GitHub Flow          | Different workflow, need design               | v0.4.0+             |
| Custom workflows     | Complexity, need user feedback                | v0.4.0+             |
| TUI mode             | Nice-to-have, not essential                   | v0.5.0+             |
| Interactive mode     | Complexity, evaluate need                     | v0.4.0+             |
| Support branches     | Less common, evaluate demand                  | v0.2.0 or v0.3.0    |

______________________________________________________________________

## Scope Boundaries

### We Handle

- Local Git operations via git CLI
- Branch creation, switching, merging
- Tag creation
- Configuration management (local/global)
- Branch name validation

### We Don't Handle

- Remote authentication (use git credential helpers)
- Pull request creation (use gzh-cli-gitforge)
- Issue tracking integration
- Code review integration
- Merge conflict resolution

______________________________________________________________________

## Feature Request Criteria

New features must meet at least one:

1. **Core Git-flow** - Part of standard git-flow workflow
2. **Safety** - Prevents common mistakes
3. **Usability** - Significantly improves user experience
4. **Integration** - Required for gzh-cli ecosystem

Features that don't meet these criteria should be deferred or rejected.
