# Metrics

## Success Criteria

### Primary KPIs

| KPI                    | Target       | Measurement                        |
| ---------------------- | ------------ | ---------------------------------- |
| Operation time         | p95 < 200ms  | Benchmark tests                    |
| Command success rate   | >= 99%       | Error tracking in tests            |
| git-flow compatibility | 100%         | E2E tests against git-flow repos   |
| Test coverage          | >= 80%       | Go coverage reports                |

### Secondary KPIs

| KPI                    | Target       | Measurement                        |
| ---------------------- | ------------ | ---------------------------------- |
| Binary size            | < 20MB       | Build output                       |
| Startup time           | < 50ms       | Benchmark tests                    |
| Memory usage           | < 128MB      | Runtime profiling                  |
| Documentation coverage | 100%         | GoDoc and CLI help audit           |

______________________________________________________________________

## Quality Metrics

### Code Quality

| Metric              | Target       | Tool                    |
| ------------------- | ------------ | ----------------------- |
| golangci-lint       | 0 issues     | make lint               |
| gofumpt compliance  | 100%         | make fmt                |
| Cyclomatic complexity| < 15/func   | gocyclo                 |
| Package coverage    | pkg>=85%, internal>=80% | go test -cover |

### Security

| Metric              | Target       | Tool                    |
| ------------------- | ------------ | ----------------------- |
| gosec issues        | 0 critical   | gosec                   |
| Input validation    | 100%         | Code review + tests     |
| Dependency vulns    | 0 known      | govulncheck             |

______________________________________________________________________

## User Experience Metrics

### Learnability

| Metric              | Target       | Measurement             |
| ------------------- | ------------ | ----------------------- |
| Time to first flow  | < 5 min      | User testing            |
| git-flow migration  | < 10 min     | User testing            |
| Command discoverability | High     | --help coverage         |

### Error Experience

| Metric              | Target       | Measurement             |
| ------------------- | ------------ | ----------------------- |
| Error message clarity| High        | User feedback           |
| Recovery guidance   | Always       | Error message audit     |
| Safe failure        | 100%         | E2E error tests         |

______________________________________________________________________

## Adoption Metrics (Post-Release)

### Usage

| Metric              | Target       | When                    |
| ------------------- | ------------ | ----------------------- |
| GitHub stars        | 50           | 3 months post-release   |
| GitHub stars        | 200          | 6 months post-release   |
| GitHub stars        | 500          | 12 months post-release  |

### Community

| Metric              | Target       | When                    |
| ------------------- | ------------ | ----------------------- |
| Contributors        | 3+           | 6 months post-release   |
| Issues resolved     | 80% in 30d   | Ongoing                 |
| Documentation PRs   | Welcome      | Ongoing                 |

______________________________________________________________________

## Validation Methods

### Automated

- Unit tests with coverage reporting
- Integration tests with real git repos
- E2E tests for complete workflows
- Benchmark tests for performance
- Linting and security scanning

### Manual

- User testing for UX metrics
- Documentation review
- git-flow compatibility verification
- Cross-platform testing

______________________________________________________________________

## Reporting

### Pre-Release

- Coverage report in CI
- Performance benchmarks in CI
- Security scan results

### Post-Release

- GitHub Insights for adoption
- Issue triage for quality feedback
- User surveys (if applicable)
