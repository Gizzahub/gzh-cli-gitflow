# Principles

## Core Principles

### P1. Safety First

**Statement**: Never corrupt repository state or lose work.

**Implications**:

- Validate merge state before finish operations
- Require clean working directory by default
- Provide clear error messages with recovery steps
- Use `--dry-run` for destructive operations

### P2. Git-flow Compatibility

**Statement**: Familiar to existing git-flow users.

**Implications**:

- Command structure mirrors git-flow
- Branch naming conventions match git-flow defaults
- Support migration from existing git-flow setup
- Don't break existing workflows

### P3. Library First

**Statement**: Core logic is reusable outside CLI.

**Implications**:

- `pkg/*` contains all business logic
- CLI is a thin wrapper
- No CLI framework in library packages
- Full test coverage on library layer

### P4. Minimal Dependencies

**Statement**: Small binary, fast startup.

**Implications**:

- Use standard library where possible
- Only essential external dependencies
- No network calls for local operations
- No background daemons

### P5. Explicit Over Implicit

**Statement**: No magic, predictable behavior.

**Implications**:

- Show what will happen before doing it
- Require explicit flags for non-default behavior
- Log operations for debugging
- Configuration is transparent

______________________________________________________________________

## Trade-off Priority

When principles conflict, use this priority:

1. **Safety** - Never lose data or corrupt state
2. **Compatibility** - Don't break git-flow muscle memory
3. **Explicitness** - Users should understand what's happening
4. **Library First** - Maintainability matters
5. **Minimal Dependencies** - Nice to have, not critical

______________________________________________________________________

## Decision Framework

When making design decisions:

1. Does it preserve safety? → Must be yes
2. Is it git-flow compatible? → Strongly preferred
3. Is it explicit to users? → Should be clear
4. Is it in the library layer? → Prefer library
5. Does it add dependencies? → Justify if needed

______________________________________________________________________

## Examples

### Example 1: Merge Conflict Detection

**Situation**: User runs `gz-flow feature finish` with conflicts

**P1 (Safety)** says: Don't attempt merge, show clear error

**P5 (Explicit)** says: Tell user exactly which files conflict

**Decision**: Detect conflicts before merge, show file list, exit with error

### Example 2: Branch Name Validation

**Situation**: User tries to create feature with invalid characters

**P1 (Safety)** says: Reject to prevent git issues

**P2 (Compatibility)** says: Match git-flow validation rules

**Decision**: Validate against git-flow rules, reject with clear message

### Example 3: Remote Push After Finish

**Situation**: Should we auto-push after finish?

**P5 (Explicit)** says: Don't do things implicitly

**P2 (Compatibility)** says: git-flow doesn't auto-push

**Decision**: Don't push by default, require `--push` flag
