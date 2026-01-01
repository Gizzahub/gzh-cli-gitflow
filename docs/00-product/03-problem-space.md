# Problem Space

## Target Users

### Primary Persona: Team Lead / Senior Developer

**Profile:**

- Manages 3-10 developers
- Responsible for release process
- Uses Git-flow or similar workflow
- Values consistency across team

**Pain Points:**

- Team members create branches with inconsistent naming
- Manual merge process is error-prone
- Release tagging is forgotten or inconsistent
- New team members don't follow the workflow

**Needs:**

- Enforce branch naming conventions
- Automate release tagging
- Onboard new developers quickly
- Track workflow state

### Secondary Persona: Solo Developer

**Profile:**

- Works on personal or small team projects
- Wants structure without overhead
- May release infrequently
- Values simplicity

**Pain Points:**

- Remembers git-flow exists but not the commands
- Manual process for versioning
- Inconsistent between projects

**Needs:**

- Simple commands that just work
- Sensible defaults
- Quick setup

### Tertiary Persona: CI/CD Engineer

**Profile:**

- Builds automation pipelines
- Needs predictable branch patterns
- Values reliability

**Pain Points:**

- Unpredictable branch names break automation
- Manual tagging interrupts pipeline
- Different projects use different patterns

**Needs:**

- Consistent branch naming
- Scriptable commands
- Non-interactive mode

______________________________________________________________________

## Current Solutions and Gaps

### git-flow (AVH Edition)

**What it does well:**

- Established, well-known
- Good documentation
- Flexible configuration

**Gaps:**

- Bash-based, inconsistent on Windows
- No library for integration
- Development stalled
- Installation varies by platform

### Manual Git Commands

**What it does well:**

- No additional tools
- Full control
- Works everywhere

**Gaps:**

- Error-prone (wrong base branch, typos)
- Inconsistent naming
- Easy to forget steps (tagging, double merge)
- No workflow state tracking

### git-flow-next

**What it does well:**

- Modern Go implementation
- Active development
- Better conflict handling

**Gaps:**

- Separate tool, not integrated
- Different command structure
- No library API

______________________________________________________________________

## Opportunity

gzh-cli-gitflow can fill these gaps:

| Gap                     | Solution                              |
| ----------------------- | ------------------------------------- |
| Cross-platform          | Go binary, works identically          |
| Integration             | Part of gzh-cli ecosystem             |
| Library API             | `pkg/*` reusable in other tools       |
| Familiar commands       | Compatible with git-flow              |
| Modern development      | Active maintenance, Go best practices |

______________________________________________________________________

## Market Size Indicators

- git-flow GitHub stars: 26k+
- Monthly "git flow" searches: 40k+
- Teams using Git-flow model: Significant portion of enterprise

**Conclusion**: Git-flow is established but tooling has gaps. Modern implementation has opportunity.
