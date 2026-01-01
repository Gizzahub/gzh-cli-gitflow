# Release Command Implementation Design

**Project**: gzh-cli-gitflow
**Date**: 2026-01-02
**Status**: Approved
**Author**: Claude (brainstorming session)

---

## 1. Overview

Implement `gz-flow release start` and `gz-flow release finish` commands to complete the core git-flow lifecycle. These commands manage release branches for preparing production releases.

### Goals

- Enable semantic versioning workflow (X.Y.Z format)
- Automate release branch creation and merging
- Create git tags for releases
- Follow established patterns from feature command
- Maintain production release integrity (master first strategy)

### Non-Goals

- Calendar versioning (CalVer)
- Pre-release version support (1.0.0-beta) - documented for future
- Automated changelog generation
- Release notes management

---

## 2. Architecture & Component Integration

### Component Reuse

The release commands follow the same architecture pattern as feature command:

| Component | Usage |
|-----------|-------|
| `internal/gitcmd` | Safe git operations (checkout, merge, tag) |
| `internal/validator` | Add `ValidateVersion()` for semver |
| `pkg/config` | Load release config (branches, prefixes, tag format) |
| `internal/preflight` | Pre-flight checks before finish |

### New Components

**Version Validator:**
```go
// internal/validator/version.go
func ValidateVersion(version string) error {
    // Strict semver: X.Y.Z where X, Y, Z are integers
    pattern := `^\d+\.\d+\.\d+$`
    if !regexp.MustCompile(pattern).MatchString(version) {
        return fmt.Errorf("invalid version format (expected: X.Y.Z)")
    }
    return nil
}
```

**Git Tag Functions:**
```go
// internal/gitcmd/git.go
func (e *Executor) CreateTag(ctx context.Context, tag, message string) error
func (e *Executor) TagExists(ctx context.Context, tag string) (bool, error)
```

### Integration Flow

**Release Start:**
```
User Input → ValidateVersion → LoadConfig → BranchExists →
Checkout(develop) → CreateBranch(release/X.Y.Z)
```

**Release Finish:**
```
User Input → ValidateVersion → LoadConfig → PreflightChecks →
Merge(release→master) → CreateTag → Merge(release→develop) → DeleteBranch
```

---

## 3. Release Start Implementation

### Command Behavior

```go
func runReleaseStart(cmd *cobra.Command, args []string) error {
    version := args[0]

    // 1. Validate version format (strict semver)
    if err := validator.ValidateVersion(version); err != nil {
        return fmt.Errorf("invalid version: %v\n💡 Use semver format: 1.0.0", err)
    }

    // 2. Load config
    cfg, err := config.LoadFromDir(".")
    if err != nil {
        fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
        cfg = config.Default()
    }

    // 3. Check if release branch already exists
    releaseBranch := cfg.Prefixes.Release + version
    exists, _ := git.BranchExists(ctx, releaseBranch)
    if exists {
        return fmt.Errorf("release branch '%s' already exists", releaseBranch)
    }

    // 4. Context hint: warn if not on develop
    currentBranch, _ := git.CurrentBranch(ctx)
    if currentBranch != cfg.Branches.Develop {
        fmt.Printf("⚠️  You're on '%s', not '%s'\n", currentBranch, cfg.Branches.Develop)
        fmt.Printf("💡 Will checkout '%s' first\n\n", cfg.Branches.Develop)
    }

    // 5. Create release branch from develop
    git.Checkout(ctx, cfg.Branches.Develop)
    git.CreateBranch(ctx, releaseBranch)

    fmt.Printf("✅ Started release branch '%s'\n", releaseBranch)
    fmt.Printf("📍 Switched to branch '%s'\n", releaseBranch)

    return nil
}
```

### Key Features

- **Strict semver validation**: Only `X.Y.Z` format accepted
- **Context-aware hints**: Warns if not on develop branch
- **No Guardian rules**: Release versions don't need naming policies (by default)
- **Simple, focused**: Just create the branch from develop

### Examples

```bash
# Success
gz-flow release start 1.0.0
✅ Started release branch 'release/1.0.0'
📍 Switched to branch 'release/1.0.0'

# Validation error
gz-flow release start v1.0.0
❌ invalid version: invalid version format (expected: X.Y.Z)
💡 Use semver format: 1.0.0

# Already exists
gz-flow release start 1.0.0
❌ release branch 'release/1.0.0' already exists
```

---

## 4. Release Finish Implementation

### Master First Merge Strategy

The implementation uses a "master first" strategy to protect production releases:

**Merge Sequence:**
1. ✅ Merge `release/X.Y.Z` → `master` (--no-ff)
2. 🏷️ Create tag `vX.Y.Z` on master
3. ✅ Merge `release/X.Y.Z` → `develop` (--no-ff)
4. 🗑️ Delete `release/X.Y.Z` branch

**Rationale:**
- Production release (master + tag) completes first
- If develop merge fails, production is already released
- Tag is created on stable master branch
- Develop merge failure doesn't break production

### Command Behavior

```go
func runReleaseFinish(cmd *cobra.Command, args []string) error {
    version := args[0]

    // 1. Validate version
    if err := validator.ValidateVersion(version); err != nil {
        return fmt.Errorf("invalid version: %v", err)
    }

    // 2. Load config
    cfg, err := config.LoadFromDir(".")
    if err != nil {
        fmt.Printf("⚠️  Failed to load config, using defaults: %v\n", err)
        cfg = config.Default()
    }

    releaseBranch := cfg.Prefixes.Release + version
    masterBranch := cfg.Branches.Master
    developBranch := cfg.Branches.Develop

    // 3. Pre-flight checks
    checker := preflight.NewChecker(git, masterBranch)
    results := checker.RunAll(ctx)
    fmt.Println("🔍 Pre-flight checks:")
    fmt.Print(results.String())
    if results.HasErrors() {
        return fmt.Errorf("pre-flight checks failed")
    }

    // 4. Verify release branch exists
    exists, _ := git.BranchExists(ctx, releaseBranch)
    if !exists {
        return fmt.Errorf("release branch '%s' does not exist", releaseBranch)
    }

    // 5. STEP 1: Merge to master (--no-ff)
    git.Checkout(ctx, masterBranch)
    if err := git.Merge(ctx, releaseBranch, true); err != nil {
        return fmt.Errorf("merge to %s failed: %v", masterBranch, err)
    }
    fmt.Printf("✅ Merged '%s' into '%s'\n", releaseBranch, masterBranch)

    // 6. STEP 2: Create tag on master
    if !noTag {
        tagName := fmt.Sprintf(cfg.Options.TagFormat, version)

        // Check if tag already exists
        exists, _ := git.TagExists(ctx, tagName)
        if exists {
            return fmt.Errorf("tag '%s' already exists\n💡 Use different version or delete existing tag", tagName)
        }

        message := tagMessage
        if message == "" {
            message = fmt.Sprintf("Release version %s", version)
        }
        if err := git.CreateTag(ctx, tagName, message); err != nil {
            return fmt.Errorf("failed to create tag: %v", err)
        }
        fmt.Printf("🏷️  Created tag '%s'\n", tagName)
    }

    // 7. STEP 3: Merge to develop
    git.Checkout(ctx, developBranch)
    if err := git.Merge(ctx, releaseBranch, true); err != nil {
        fmt.Printf("⚠️  PARTIAL SUCCESS:\n")
        fmt.Printf("  ✅ Merged to %s and tagged '%s'\n", masterBranch, tagName)
        fmt.Printf("  ❌ Merge to %s failed: %v\n", developBranch, err)
        fmt.Printf("\n💡 To complete:\n")
        fmt.Printf("  1. git checkout %s\n", developBranch)
        fmt.Printf("  2. git merge --no-ff %s\n", releaseBranch)
        fmt.Printf("  3. Resolve conflicts and commit\n")
        return err
    }
    fmt.Printf("✅ Merged '%s' into '%s'\n", releaseBranch, developBranch)

    // 8. STEP 4: Delete release branch
    if cfg.Options.DeleteBranchAfterFinish && !keepBranch {
        if err := git.DeleteBranch(ctx, releaseBranch); err != nil {
            fmt.Printf("⚠️  Failed to delete branch: %v\n", err)
        } else {
            fmt.Printf("🗑️  Deleted branch '%s'\n", releaseBranch)
        }
    }

    return nil
}
```

### Flags

| Flag | Short | Default | Purpose |
|------|-------|---------|---------|
| `--message` | `-m` | "Release version X.Y.Z" | Custom tag message |
| `--no-tag` | - | false | Skip tag creation |
| `--keep` | `-k` | false | Keep release branch after finish |

### Examples

```bash
# Success
gz-flow release finish 1.0.0
🔍 Pre-flight checks:
  ✅ Working tree is clean
  ✅ Branch 'master' is up-to-date

✅ Merged 'release/1.0.0' into 'master'
🏷️  Created tag 'v1.0.0'
✅ Merged 'release/1.0.0' into 'develop'
🗑️  Deleted branch 'release/1.0.0'

# With custom tag message
gz-flow release finish 1.0.0 -m "Initial production release"
🏷️  Created tag 'v1.0.0'

# Skip tag
gz-flow release finish 1.0.0 --no-tag
✅ Merged 'release/1.0.0' into 'master'
✅ Merged 'release/1.0.0' into 'develop'

# Partial success (develop merge fails)
gz-flow release finish 1.0.0
⚠️  PARTIAL SUCCESS:
  ✅ Merged to master and tagged 'v1.0.0'
  ❌ Merge to develop failed: merge conflict

💡 To complete:
  1. git checkout develop
  2. git merge --no-ff release/1.0.0
  3. Resolve conflicts and commit
```

---

## 5. Error Handling & Edge Cases

### Critical Edge Cases

**1. Tag Already Exists**
```go
tagName := fmt.Sprintf(cfg.Options.TagFormat, version)
exists, _ := git.TagExists(ctx, tagName)
if exists {
    return fmt.Errorf("tag '%s' already exists\n💡 Use a different version or delete the existing tag", tagName)
}
```

**2. Merge Conflicts on Master**
```go
if err := git.Merge(ctx, releaseBranch, true); err != nil {
    return fmt.Errorf("merge to %s failed: %v\n💡 Resolve conflicts:\n  1. git checkout %s\n  2. git merge --no-ff %s\n  3. Resolve conflicts\n  4. git merge --continue\n  5. Retry: gz-flow release finish %s",
        masterBranch, err, masterBranch, releaseBranch, version)
}
```

**3. Merge Conflicts on Develop (Partial Success)**
- Master is already tagged (production release successful)
- Develop merge failure is recoverable
- Clear instructions provided to user
- Can retry develop merge without affecting master

**4. Version Mismatch**
```go
// Auto-detect current release branch
currentBranch, _ := git.CurrentBranch(ctx)
expectedBranch := cfg.Prefixes.Release + version
if currentBranch == expectedBranch {
    fmt.Printf("📍 Detected release version: %s\n", version)
}
```

**5. No Develop Branch**
```go
exists, _ := git.BranchExists(ctx, developBranch)
if !exists {
    fmt.Printf("⚠️  Develop branch '%s' does not exist\n", developBranch)
    fmt.Printf("💡 Skipping merge to develop\n")
    // Continue without develop merge (master-only workflow)
}
```

### Error Priority

| Severity | Issue | Handling |
|----------|-------|----------|
| Critical | Pre-flight check fails | Block operation, clear guidance |
| Critical | Master merge conflict | Block, provide resolution steps |
| Warning | Develop merge conflict | Partial success, guide completion |
| Warning | Tag creation fails | Log warning, continue if --no-tag |
| Info | Branch delete fails | Log warning, non-blocking |

---

## 6. Testing Strategy

### Unit Tests

**1. Version Validator (`internal/validator/version_test.go`)**
```go
func TestValidateVersion(t *testing.T) {
    tests := []struct {
        name    string
        version string
        wantErr bool
    }{
        {"valid semver", "1.0.0", false},
        {"valid with large numbers", "12.34.56", false},
        {"invalid with v prefix", "v1.0.0", true},
        {"invalid two parts", "1.0", true},
        {"invalid four parts", "1.0.0.0", true},
        {"invalid with dash", "1.0.0-beta", true},
        {"invalid with letters", "1.0.x", true},
        {"empty", "", true},
    }
}
```

**2. GitCmd Tag Functions (`internal/gitcmd/git_test.go`)**
```go
func TestExecutor_CreateTag(t *testing.T) {
    // Test tag creation with message
    // Test tag validation (name format)
    // Test error when tag exists
}

func TestExecutor_TagExists(t *testing.T) {
    // Test tag existence check
    // Test non-existent tag
}
```

**3. Release Command Integration (`tests/integration/release_test.go`)**
```go
func TestReleaseStart(t *testing.T) {
    dir := setupTestRepo(t)
    binary := buildBinary(t)

    // Test successful release start
    run(t, dir, binary, "release", "start", "1.0.0")

    // Verify release/1.0.0 branch created
    // Verify on release/1.0.0 branch
}

func TestReleaseFinish(t *testing.T) {
    // Setup: create release/1.0.0 branch with commits
    // Run: release finish 1.0.0
    // Verify: master has merge, tag exists, develop has merge, branch deleted
}
```

### Coverage Goals

| Component | Target | Rationale |
|-----------|--------|-----------|
| `internal/validator/version.go` | 100% | Security-critical validation |
| `internal/gitcmd` (tag functions) | 85%+ | Core git operations |
| `cmd/release.go` | Integration tests | Command orchestration |

### Manual Testing Scenarios

1. **Happy path**: start → commits → finish
2. **Conflict on master merge**: Verify clear error message
3. **Conflict on develop merge**: Verify partial success handling
4. **Tag already exists**: Verify blocked with guidance
5. **Release branch doesn't exist**: Verify error message
6. **No develop branch**: Verify graceful handling (master-only)

---

## 7. Alternative Merge Strategies (Documentation Only)

### Current: Master First (Strategy A) ✅

**Implementation:**
1. Merge release → master (--no-ff)
2. Create tag on master
3. Merge release → develop (--no-ff)
4. Delete release branch

**Pros:**
- Production release (master + tag) completes first
- If develop merge fails, production is already released
- Tag is created on stable master branch

**Cons:**
- Develop merge failure leaves inconsistent state (requires manual fix)

---

### Alternative: Develop First (Strategy B) ⚠️

**Not Implemented - Documented for Reference**

**Sequence:**
1. Merge release → develop
2. Merge release → master
3. Create tag on master
4. Delete release branch

**Use Case:** When you want develop to get changes first, safer rollback

**Risk:** Master merge failure after develop has changes creates inconsistency

**When to Consider:**
- Development-heavy workflow
- Master releases are rare
- You want develop to be source of truth

---

### Alternative: Git-Flow Classic (Strategy C) ⚠️

**Not Implemented - Documented for Reference**

**Sequence:**
1. Merge release → master (--no-ff)
2. Create tag on master
3. Merge **master** → develop (brings tagged version to develop)
4. Delete release branch

**Use Case:** Original git-flow behavior

**Difference:** Develop gets master's tagged commit, not release branch directly

**When to Consider:**
- Strict adherence to original git-flow
- Want develop to have exact tagged state
- Master is always canonical

---

## 8. Guardian Mode (Future Enhancement)

### Current Implementation

**Strict Semver Validation:**
```go
pattern := `^\d+\.\d+\.\d+$`
```

**Accepts:**
- `1.0.0`
- `2.5.3`
- `10.0.0`

**Rejects:**
- `v1.0.0` (prefix)
- `1.0` (incomplete)
- `1.0.0-beta` (pre-release)

### Future: Configuration-Based Validation

**Status:** Design approved, implementation deferred to v0.2.0

**Configuration:**
```yaml
# .gzflow.yaml
guardian:
  enabled: true
  naming:
    release:
      pattern: "^\\d+\\.\\d+\\.\\d+(-[a-z0-9]+)?$"
      examples:
        - "1.0.0"
        - "2.1.0-beta.1"
        - "3.0.0-rc.2"
```

**Implementation:**
```go
// In runReleaseStart/Finish
if cfg.Guardian.Enabled && cfg.Guardian.Naming.Release.Pattern != "" {
    if err := cfg.Guardian.Naming.Release.Validate(version); err != nil {
        return fmt.Errorf("guardian: %v", err)
    }
} else {
    // Fallback to strict semver
    if err := validator.ValidateVersion(version); err != nil {
        return err
    }
}
```

**Use Cases:**
- **CalVer**: `2024.01.15`, `2026.1.2`
- **Pre-release**: `1.0.0-beta.1`, `2.0.0-rc.2`
- **Build metadata**: `1.0.0+20240115`
- **Custom schemes**: Company-specific versioning

**Benefits:**
- Flexible per-project configuration
- Backward compatible (defaults to strict semver)
- Consistent with Guardian philosophy

---

## 9. Implementation Checklist

### Phase 1: Core Functionality
- [ ] `internal/validator/version.go` - Semver validation
- [ ] `internal/gitcmd/git.go` - Tag functions (CreateTag, TagExists)
- [ ] `cmd/release.go` - Implement runReleaseStart
- [ ] `cmd/release.go` - Implement runReleaseFinish

### Phase 2: Testing
- [ ] `internal/validator/version_test.go` - 8+ test cases
- [ ] `internal/gitcmd/git_test.go` - Tag function tests
- [ ] `tests/integration/release_test.go` - E2E tests

### Phase 3: Documentation
- [ ] Update CLAUDE.md with release command examples
- [ ] Add troubleshooting guide for merge conflicts
- [ ] Document tag format configuration

### Phase 4: Quality
- [ ] Linter passes (0 issues)
- [ ] Tests pass (100%)
- [ ] Coverage: validator 100%, gitcmd 85%+
- [ ] Manual testing: all 6 scenarios

---

## 10. Approval

- [x] Design approved (2026-01-02)
- [ ] Implementation plan created
- [ ] Code review completed
- [ ] Integration tests passing
- [ ] Release approved

---

**End of Document**
