# DX 개선 & Guardian Mode 기능 설계

**Project**: gzh-cli-gitflow
**Date**: 2026-01-01
**Status**: Approved
**Author**: Claude (brainstorming session)

---

## 1. 배경 및 문제점

### 1.1 DX (Developer Experience) 문제

| 문제 | 설명 |
|-----|------|
| 브랜치 이름 기억 | "저 feature 이름이 뭐였지?" 매번 `git branch` 확인 |
| 워크플로우 순서 실수 | develop에서 시작해야 하는데 main에서 시작 |
| 사전 체크 누락 | uncommitted changes, 충돌 가능성 확인 안 함 |
| 동시 작업 파악 어려움 | 여러 feature 브랜치 전환 시 현재 상태 파악 힘듦 |

### 1.2 팀 협업 문제

| 문제 | 설명 |
|-----|------|
| 네이밍 불일치 | `feature/login` vs `feature/user-login` vs `feat/login` |
| 동시 작업 충돌 | 같은 파일을 여러 명이 수정 중인지 모름 |
| 릴리즈 타이밍 혼란 | 누가 release 브랜치를 만들었는지, 언제 finish할지 |
| 방치 브랜치 | merge 안 된 feature가 쌓임 |

---

## 2. 설계 접근법

**Smart Defaults + Guardian Mode** 조합 채택:

- **Smart Defaults**: 기존 명령어에 똑똑한 기능 추가 (DX 개선)
- **Guardian Mode**: 정책 강제 레이어 추가 (팀 일관성)

---

## 3. Smart Defaults 기능

### 3.1 자동 브랜치 감지 (`--auto`)

```bash
# 현재: 브랜치 이름 직접 입력 필수
gz-flow feature finish user-auth

# 개선: 현재 브랜치에서 자동 감지
gz-flow feature finish          # feature/* 브랜치면 자동 인식
gz-flow feature finish --auto   # 명시적 자동 모드
```

**구현 요구사항:**
- 현재 브랜치가 `feature/*` 패턴인지 확인
- 패턴 불일치 시 명확한 에러 메시지
- `--auto` 플래그로 명시적 활성화 가능

### 3.2 인터랙티브 브랜치 선택 (`--pick`)

```bash
gz-flow feature finish --pick
# ? Select feature branch to finish:
#   ❯ feature/user-auth (3 days ago, 5 commits ahead)
#     feature/payment (1 week ago, 12 commits ahead)
#     feature/dashboard (2 weeks ago, stale)
```

**구현 요구사항:**
- 해당 타입의 모든 브랜치 목록 표시
- 브랜치별 메타데이터 표시 (age, commits ahead)
- 화살표 키로 선택 가능한 인터랙티브 UI

### 3.3 Pre-flight 체크 (자동)

모든 `finish` 명령 전에 자동 실행:

| 체크 항목 | 설명 |
|----------|------|
| Clean working directory | uncommitted changes 확인 |
| Base branch up-to-date | develop/main 최신 상태 확인 |
| Merge conflict detection | `git merge --no-commit --no-ff` dry-run |
| Branch exists | 대상 브랜치 존재 여부 |

**실패 시 동작:**
```bash
gz-flow feature finish user-auth
# ❌ Pre-flight check failed:
#    • Working directory not clean (2 uncommitted files)
#    • develop is 3 commits behind origin/develop
#
# 💡 Fix suggestions:
#    1. git stash or git commit
#    2. git checkout develop && git pull
```

### 3.4 컨텍스트 인식 도움말

잘못된 상황에서 실행 시 가이드:

```bash
gz-flow feature start login   # main 브랜치에서 실행
# ⚠️  You're on 'main', not 'develop'
# 💡 Hint: Switch to develop first, or use --from=main if intentional
```

```bash
gz-flow release finish 1.2.0  # release 브랜치가 아닌 곳에서 실행
# ⚠️  You're on 'develop', not on a release branch
# 💡 Hint: Switch to release/1.2.0 first
```

---

## 4. Guardian Mode 기능

### 4.1 설정 스키마

`.gzflow.yaml`:

```yaml
guardian:
  enabled: true

  naming:
    feature:
      pattern: "^[a-z]+(-[a-z0-9]+)*$"  # kebab-case only
      max_length: 50
      forbidden: ["test", "temp", "wip"]
    release:
      pattern: "^\\d+\\.\\d+\\.\\d+$"   # semver only
    hotfix:
      pattern: "^\\d+\\.\\d+\\.\\d+$"   # semver only

  workflow:
    require_clean_tree: true       # finish 전 clean 필수
    require_up_to_date: true       # base 브랜치 최신 필수
    block_direct_main_commit: true # main 직접 커밋 차단
    max_feature_age_days: 30       # 30일 초과 브랜치 경고

  enforcement:
    mode: "warn"  # "warn" | "block"
```

### 4.2 네이밍 규칙 강제

```bash
gz-flow feature start MyFeature
# ❌ Branch name 'MyFeature' violates naming rule
# 📋 Rule: kebab-case only (e.g., my-feature)
# 💡 Suggested: my-feature

gz-flow feature start temp-fix
# ❌ Branch name 'temp-fix' contains forbidden word: 'temp'
# 📋 Forbidden words: test, temp, wip
```

### 4.3 브랜치 Audit 명령

```bash
gz-flow audit [--type=feature|release|hotfix] [--format=text|json]
```

출력 예시:
```bash
gz-flow audit
# 🔍 Branch Audit Report
# ─────────────────────────────────────
# ⚠️  Stale branches (>30 days):
#    feature/old-login (45 days, @kim)
#    feature/unused-api (60 days, @lee)
#
# ❌ Naming violations:
#    Feature_Test → should be: feature-test
#
# 📊 Summary: 2 stale, 1 violation, 5 healthy
```

### 4.4 팀 설정 초기화

```bash
gz-flow guardian init [--team] [--strict]
```

- `--team`: 팀 기본값으로 설정
- `--strict`: 엄격한 규칙 적용 (block mode)

---

## 5. 충돌 예방 & 팀 인식 기능

### 5.1 충돌 위험 감지

```bash
gz-flow status --conflicts
```

출력:
```
🔍 Conflict Risk Analysis
─────────────────────────────────────
feature/user-auth:
  ⚠️  src/auth/login.go - also modified in:
     └─ feature/payment (@park, 2 days ago)
     └─ develop (merged yesterday)

  💡 Recommend: rebase from develop before finish
```

### 5.2 파일 수정 현황 (`who` 명령)

```bash
gz-flow who <file>
```

예시:
```bash
gz-flow who src/auth/login.go
#
# 📁 src/auth/login.go modified in:
#   feature/user-auth    (+45, -12)  ← current
#   feature/payment      (+23, -5)
#   hotfix/1.2.1         (+3, -1)
```

### 5.3 릴리즈 상태 조회

```bash
gz-flow release status
```

출력:
```
🚀 Release Status
─────────────────────────────────────
Active release: release/1.3.0 (started 2 days ago)

Ready to merge (in develop, not in release):
  ✅ feature/user-auth (finished yesterday)
  ✅ feature/dashboard (finished 3 days ago)

Still in progress:
  🔄 feature/payment (5 commits ahead of develop)

💡 Use 'gz-flow release include user-auth' to cherry-pick
```

### 5.4 브랜치 정리 (`cleanup` 명령)

```bash
gz-flow cleanup [--dry-run] [--force] [--include-remote]
```

예시:
```bash
gz-flow cleanup --dry-run
# 🧹 Cleanup Preview
# ─────────────────────────────────────
# Will delete (already merged):
#   feature/old-login     → merged to develop
#   feature/legacy-api    → merged to develop
#
# Will warn (stale, not merged):
#   feature/abandoned     → 60 days, no activity
#
# Run 'gz-flow cleanup' to execute
```

---

## 6. 우선순위 및 로드맵

### 6.1 기능 우선순위

| 기능 | 카테고리 | 우선순위 | 제안 버전 | 복잡도 |
|-----|---------|---------|----------|-------|
| 자동 브랜치 감지 (`--auto`) | Smart | P1 | v0.1.0 | Low |
| Pre-flight 체크 | Smart | P1 | v0.1.0 | Medium |
| 컨텍스트 인식 도움말 | Smart | P1 | v0.1.0 | Low |
| 네이밍 규칙 강제 | Guardian | P1 | v0.1.5 | Medium |
| 인터랙티브 선택 (`--pick`) | Smart | P2 | v0.2.0 | Medium |
| 워크플로우 정책 | Guardian | P2 | v0.2.0 | Medium |
| `audit` 명령 | Guardian | P2 | v0.2.0 | Medium |
| 충돌 위험 감지 | Team | P2 | v0.2.0 | High |
| `release status` | Team | P2 | v0.2.0 | Medium |
| `cleanup` 명령 | Team | P2 | v0.2.0 | Medium |
| `who` 명령 | Team | P3 | v0.3.0 | Medium |

### 6.2 로드맵 통합

```
v0.1.0 (기존 + 추가)
├── 기존: init, feature, release, hotfix, status, list, config
└── 추가: --auto, pre-flight, 컨텍스트 도움말

v0.1.5 (신규 마일스톤)
└── Guardian 기본: 네이밍 규칙 검증

v0.2.0 (기존 + 추가)
├── 기존: publish, pull, track, delete
└── 추가: --pick, audit, cleanup, release status, 충돌 감지

v0.3.0 (기존 + 추가)
├── 기존: gzh-cli 통합
└── 추가: who 명령, 고급 Guardian 정책
```

---

## 7. Non-Goals (유지)

다음 기능은 이 프로젝트의 범위에 포함되지 않음:

| Non-Goal | 대안 |
|----------|-----|
| GitHub/GitLab API 연동 | gzh-cli-gitforge 사용 |
| 원격 브랜치 실시간 추적 | 로컬 정보 기반만 사용 |
| Git hooks 자동 설치 | CLI 실행 시점에만 정책 적용 |
| 자동 충돌 해결 | 감지 및 가이드만 제공 |

---

## 8. 기술적 고려사항

### 8.1 의존성

- **인터랙티브 UI**: `github.com/charmbracelet/bubbletea` 또는 `github.com/AlecAivazis/survey`
- **출력 포맷팅**: `github.com/fatih/color` (이미 사용 중일 가능성)
- **설정 관리**: Viper (이미 사용 중)

### 8.2 테스트 요구사항

| 기능 | 테스트 유형 |
|-----|-----------|
| Pre-flight 체크 | Unit + Integration (실제 git repo) |
| Guardian 규칙 | Unit (regex 검증) |
| Conflict 감지 | Integration (multi-branch scenarios) |
| Cleanup | Integration + E2E |

### 8.3 보안 고려

- 모든 브랜치 이름은 injection 패턴 검증 필수
- Guardian 규칙 pattern은 ReDoS 방지 검증
- `cleanup --force`는 확인 프롬프트 필수

---

## 9. 승인

- [x] 설계 방향 승인 (2026-01-01)
- [ ] 구현 착수 승인
- [ ] 코드 리뷰 완료
- [ ] 릴리즈 승인

---

**End of Document**
