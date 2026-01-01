# Roadmap

## Phases

### Phase 1: Foundation (PLANNED)

**Status**: Planned (v0.1.0)
**Goal**: Core Git-flow functionality

| Milestone | Deliverables                           | Target  |
| --------- | -------------------------------------- | ------- |
| M1.1      | Project structure, CI/CD, basic CLI    | Q1 2026 |
| M1.2      | init, feature start/finish             | Q1 2026 |
| M1.3      | release start/finish with tagging      | Q1 2026 |
| M1.4      | hotfix start/finish                    | Q1 2026 |
| M1.5      | status, list, config commands          | Q1 2026 |
| M1.6      | Documentation, testing, v0.1.0 release | Q1 2026 |

**Exit criteria**:

- All core commands working
- Unit test coverage >= 80%
- CLI reference documentation complete
- v0.1.0 released

### Phase 2: Remote Integration (PLANNED)

**Status**: Planned (v0.2.0)
**Goal**: Remote repository support

| Feature    | Description                        | Priority |
| ---------- | ---------------------------------- | -------- |
| publish    | Push flow branch to remote         | P1       |
| pull       | Pull updates from remote           | P1       |
| track      | Set up tracking for flow branches  | P2       |
| delete     | Delete local and remote branches   | P2       |

**Entry criteria**: Phase 1 complete

| Milestone | Deliverables                | Target  |
| --------- | --------------------------- | ------- |
| M2.1      | publish command             | Q2 2026 |
| M2.2      | pull and track commands     | Q2 2026 |
| M2.3      | delete with remote support  | Q2 2026 |
| M2.4      | v0.2.0 release              | Q2 2026 |

### Phase 3: gzh-cli Integration (PLANNED)

**Status**: Planned (v0.3.0)
**Goal**: Integration with main CLI

| Feature          | Description                        | Priority |
| ---------------- | ---------------------------------- | -------- |
| Library API      | Stable pkg/* API for gzh-cli       | P1       |
| CLI integration  | `gz flow` command in gzh-cli       | P1       |
| Shared config    | Integration with gz config system  | P2       |

**Entry criteria**: Phase 2 complete

| Milestone | Deliverables                | Target  |
| --------- | --------------------------- | ------- |
| M3.1      | Stable library API          | Q3 2026 |
| M3.2      | gzh-cli integration         | Q3 2026 |
| M3.3      | v0.3.0 release              | Q3 2026 |

### Phase 4: Alternative Workflows (PLANNED)

**Status**: Planned (v0.4.0)
**Goal**: Support for GitHub Flow and custom workflows

| Feature              | Description                         | Priority |
| -------------------- | ----------------------------------- | -------- |
| GitHub Flow          | Simplified main + feature workflow  | P1       |
| Trunk-Based          | Main + short-lived branches         | P2       |
| Custom workflows     | User-defined workflow definitions   | P3       |

**Entry criteria**: Phase 3 complete, user feedback collected

### Phase 5: Advanced Features (FUTURE)

**Status**: Future Planning
**Goal**: Enhanced user experience

| Feature          | Description                             | Priority |
| ---------------- | --------------------------------------- | -------- |
| Interactive mode | Guided workflows with prompts           | P2       |
| TUI              | Rich terminal UI for complex operations | P3       |
| Hooks            | Pre/post flow operation hooks           | P3       |
| Templates        | Branch description templates            | P3       |

______________________________________________________________________

## Version Summary

| Version | Focus                    | Key Features                              | Target  |
| ------- | ------------------------ | ----------------------------------------- | ------- |
| v0.1.0  | Core Git-flow            | init, feature, release, hotfix, status    | Q1 2026 |
| v0.2.0  | Remote Support           | publish, pull, track, delete              | Q2 2026 |
| v0.3.0  | gzh-cli Integration      | Library API, CLI integration              | Q3 2026 |
| v0.4.0  | Alternative Workflows    | GitHub Flow, Trunk-Based                  | Q4 2026 |
| v1.0.0  | Stable Release           | API stability guarantee                   | 2027    |

______________________________________________________________________

## Milestones

### Near-term (v0.1.0)

| Milestone        | Description             | Target  | Status |
| ---------------- | ----------------------- | ------- | ------ |
| Project setup    | Structure, CI/CD        | Q1 2026 | 📋     |
| Core commands    | feature, release, hotfix| Q1 2026 | 📋     |
| Utility commands | status, list, config    | Q1 2026 | 📋     |
| Documentation    | README, CLI reference   | Q1 2026 | 📋     |
| v0.1.0 release   | First public release    | Q1 2026 | 📋     |

### Medium-term (v0.2.0 - v0.3.0)

| Milestone           | Description              | Target  | Status |
| ------------------- | ------------------------ | ------- | ------ |
| Remote integration  | publish, pull commands   | Q2 2026 | 📋     |
| v0.2.0 release      | Remote support release   | Q2 2026 | 📋     |
| gzh-cli integration | Full library integration | Q3 2026 | 📋     |
| v0.3.0 release      | Integration release      | Q3 2026 | 📋     |

### Long-term (v0.4.0+)

| Milestone            | Description              | Target | Status |
| -------------------- | ------------------------ | ------ | ------ |
| v0.4.0               | Alternative workflows    | Q4 2026| 📋     |
| v1.0.0               | Stable API guarantee     | 2027   | 📋     |

______________________________________________________________________

## Decision Points

| Decision              | When           | Options                             |
| --------------------- | -------------- | ----------------------------------- |
| GitHub Flow design    | Phase 4 start  | Subset of git-flow, separate cmd    |
| Custom workflow format| Phase 4 mid    | YAML definition, Go plugins         |
| TUI framework         | Phase 5 start  | Bubble Tea, tview, none             |
| API versioning        | Pre-v1.0.0     | Semantic versioning commitment      |

______________________________________________________________________

## Command Roadmap

### v0.1.0 Commands

```bash
gz-flow init                     # Initialize git-flow
gz-flow feature start <name>     # Start feature branch
gz-flow feature finish <name>    # Finish feature branch
gz-flow release start <version>  # Start release branch
gz-flow release finish <version> # Finish release branch
gz-flow hotfix start <version>   # Start hotfix branch
gz-flow hotfix finish <version>  # Finish hotfix branch
gz-flow status                   # Show current state
gz-flow list [type]              # List flow branches
gz-flow config [key] [value]     # Manage configuration
```

### v0.2.0 Commands

```bash
gz-flow feature publish <name>   # Push feature to remote
gz-flow feature pull <name>      # Pull feature from remote
gz-flow feature track <name>     # Track remote feature
gz-flow feature delete <name>    # Delete feature branch
gz-flow release publish <version># Push release to remote
gz-flow hotfix publish <version> # Push hotfix to remote
```

### v0.3.0+ Commands

```bash
gz-flow support start <version>  # Long-term support branch
gz-flow support finish <version> # Finish support branch
```

______________________________________________________________________

## Legend

| Symbol | Meaning     |
| ------ | ----------- |
| ✅     | Complete    |
| 🔄     | In progress |
| 📋     | Planned     |
| ⏸️     | On hold     |
| ❌     | Cancelled   |
