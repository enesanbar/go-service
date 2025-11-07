# CI/CD Workflows

This directory contains GitHub Actions workflows for the multi-module Go monorepo.

## Workflows

### `ci.yml` - Continuous Integration

**Purpose**: Validates code quality on all branches and pull requests.

**Triggers**:
- Push to any branch
- Pull request opened, synchronized, or reopened

**Jobs**:

1. **detect-changes**: Identifies which modules have changed
   - Compares current state against module tags
   - Outputs JSON array of changed modules
   - Handles modules without tags (treats as new)
   - Detects non-module file changes (runs all modules)

2. **test**: Runs tests in parallel for each changed module
   - Matrix strategy across changed modules
   - Collects test coverage data
   - Uploads coverage artifacts (30-day retention)
   - Checks coverage against threshold (default: 80%)
   - Warnings for low coverage (non-breaking)
   - Job summary with test and coverage status

3. **lint**: Runs golangci-lint in parallel for each changed module
   - Uses repository-wide `.golangci.yml` config
   - Matrix strategy across changed modules
   - PR annotations for violations
   - Job summary with lint results

**Environment Variables**:
- `COVERAGE_THRESHOLD`: Coverage percentage threshold (default: 80)
- `COVERAGE_BREAKING`: Whether to fail on low coverage (default: false)
- `COVERAGE_RETENTION_DAYS`: Coverage artifact retention (default: 30)

**Status Checks**:
- All jobs must pass for PR merge
- Test failures block merge
- Lint violations block merge
- Low coverage issues warnings (non-breaking)

**Performance**:
- Parallel execution reduces time by 50%+ for 3+ modules
- Single module validation: <2 minutes
- Caching for Go dependencies and build artifacts

## Usage

### Viewing Test Results

Test results appear in:
- PR checks summary
- Workflow run summary page
- Job summaries with coverage percentages

### Viewing Coverage Reports

1. **In PR Comments/Summaries**: Coverage percentage displayed per module
2. **Download Artifacts**: 
   ```bash
   # From GitHub Actions UI
   Actions → Workflow Run → Artifacts → coverage-<module-path>
   ```
3. **Local Viewing**:
   ```bash
   # Download coverage file
   go tool cover -html=coverage/module-coverage.out
   ```

### Bypassing Checks (Admin Only)

Admins can merge PRs with failing checks if necessary, but this is not recommended.

## Troubleshooting

### No modules detected

**Cause**: No changes since last tag  
**Solution**: This is expected behavior - no tests/linting needed

### Coverage warnings

**Cause**: Coverage below threshold (80%)  
**Solution**: Add tests to increase coverage (non-blocking warning)

### Lint failures

**Cause**: Code quality issues detected  
**Solution**: Fix violations reported in PR annotations

### Test failures

**Cause**: Tests failing for changed modules  
**Solution**: Fix failing tests before merge

### Permission errors

**Cause**: Missing workflow permissions  
**Solution**: Ensure repository settings allow Actions to write PR comments

## Configuration

### Adjusting Coverage Threshold

Edit `.github/workflows/ci.yml`:
```yaml
env:
  COVERAGE_THRESHOLD: 75  # Lower to 75%
```

### Changing Linter Rules

Edit `.golangci.yml` at repository root.

### Disabling Specific Jobs

Add `if: false` to job definition (not recommended).

## Dependencies

- GitHub Actions runners (ubuntu-latest)
- Go 1.21+ (version per module's go.mod)
- golangci-lint 1.55+
- Bash 5.x for scripts

## Related Files

- `.github/scripts/detect-changed-modules.sh`: Module change detection script
- `.github/scripts/lib/`: Utility libraries
- `.golangci.yml`: Repository-wide linter configuration
