# CI/CD Scripts

This directory contains Bash scripts for the multi-module CI/CD pipeline.

## Structure

```text
.github/scripts/
├── detect-changed-modules.sh    # Main: Module change detection
├── analyze-commits.sh           # Conventional commit analysis
├── calculate-version.sh         # Semantic version calculation
├── create-tags.sh               # Git tag creation and pushing
├── lib/                          # Utility libraries
│   ├── git-utils.sh             # Git operations
│   ├── semver-utils.sh          # Semantic versioning
│   └── module-utils.sh          # Module discovery and validation
└── tests/                        # BATS integration tests
    ├── test-detect-modules.bats # Tests for change detection
    ├── test-analyze-commits.bats # Tests for commit analysis
    ├── test-calculate-version.bats # Tests for version calculation
    └── fixtures/                 # Test data
```

## Scripts

### `detect-changed-modules.sh`

**Purpose**: Detect which Go modules have changed by comparing against their last tagged version.

**Usage**:
```bash
./detect-changed-modules.sh [OPTIONS]

OPTIONS:
    --base-ref REF      Base git reference for comparison (default: auto-detect from tags)
    --head-ref REF      Head git reference (default: HEAD)
    --format FORMAT     Output format: json or text (default: json)
    --include-all       Include all modules even if unchanged (default: false)
    -h, --help          Show help message
```

**Examples**:
```bash
# Auto-detect changes using module tags
./detect-changed-modules.sh

# Compare against specific ref
./detect-changed-modules.sh --base-ref main --head-ref feature-branch

# Text output
./detect-changed-modules.sh --format text

# Include all modules
./detect-changed-modules.sh --include-all
```

**Output (JSON)**:
```json
[
  {
    "path": "cache/inmemory",
    "name": "github.com/example/go-service/cache/inmemory",
    "last_tag": "cache/inmemory/v0.1.0",
    "has_changes": true,
    "changed_files": "cache/inmemory/inmemory.go,cache/inmemory/config.go"
  }
]
```

**Exit Codes**:
- `0`: Success
- `1`: Invalid arguments
- `2`: Git operation error

**Algorithm**:

1. Discover all Go modules (find go.mod files)
2. For each module, get its latest tag (e.g., `cache/inmemory/v0.1.0`)
3. If no tag exists, treat as new module (all files changed)
4. Run `git diff` between tag and HEAD for module directory
5. Check for non-module file changes (README, CI config)
6. If non-module files changed, return all modules
7. Output JSON array of changed modules

---

### `analyze-commits.sh`

**Purpose**: Analyze conventional commits for a module to determine version bump type.

**Usage**:

```bash
./analyze-commits.sh <module-path> [OPTIONS]

OPTIONS:
    --from-tag TAG      Compare from this tag (default: latest tag for module)
    --output FORMAT     Output format: json or text (default: text)
    -h, --help          Show help message
```

**Examples**:

```bash
# Analyze commits since last tag
./analyze-commits.sh cache/inmemory --output json

# Analyze from specific tag
./analyze-commits.sh core/errors --from-tag core/errors/v0.1.0
```

**Output (JSON)**:

```json
{
  "module": "cache/inmemory",
  "recommended_bump": "minor",
  "commit_count": 3,
  "has_breaking": false,
  "has_feat": true,
  "has_fix": true,
  "commits": [
    {
      "hash": "abc123",
      "type": "feat",
      "bump": "minor",
      "subject": "add caching feature"
    }
  ],
  "warnings": []
}
```

**Exit Codes**:

- `0`: Success
- `1`: Invalid arguments
- `2`: Git operation error
- `3`: No commits found for module

**Commit Type Classification**:

- `feat:` → minor bump
- `fix:` → patch bump
- `feat!:` or `BREAKING CHANGE:` → major bump
- Other types → patch bump (with warning if non-conventional)

---

### `calculate-version.sh`

**Purpose**: Calculate the next semantic version for a module based on commit analysis.

**Usage**:

```bash
./calculate-version.sh <module-path> [OPTIONS]

OPTIONS:
    --override-version V  Manually specify next version (must be valid semver)
    --output FORMAT       Output format: json or text (default: text)
    -h, --help            Show help message
```

**Examples**:

```bash
# Auto-calculate next version
./calculate-version.sh cache/inmemory --output json

# Manual version override
./calculate-version.sh core/errors --override-version 2.0.0
```

**Output (JSON)**:

```json
{
  "module": "cache/inmemory",
  "current_version": "0.1.0",
  "next_version": "0.2.0",
  "bump_type": "minor",
  "message": "Calculated based on conventional commits"
}
```

**Exit Codes**:

- `0`: Success
- `1`: Invalid arguments
- `2`: Version calculation failed
- `3`: Invalid version format

**Version Calculation Logic**:

1. Get current version from latest module tag
2. Analyze commits using `analyze-commits.sh`
3. Apply bump type: major (X.0.0), minor (x.Y.0), patch (x.y.Z)
4. If no tag exists, start at `0.0.1`
5. If no commits found, return current version (no bump)

---

### `create-tags.sh`

**Purpose**: Create and push directory-prefixed git tags for modules.

**Usage**:

```bash
./create-tags.sh <module-path> <version> [OPTIONS]

OPTIONS:
    --dry-run           Show what would be done without making changes
    --push              Push tags to remote repository
    -h, --help          Show help message
```

**Examples**:

```bash
# Create tag locally (test)
./create-tags.sh cache/inmemory 0.1.0 --dry-run

# Create and push tag
./create-tags.sh core/errors 1.2.3 --push
```

**Exit Codes**:

- `0`: Success
- `1`: Invalid arguments
- `2`: Tag creation failed
- `3`: Tag push failed after retries

**Tag Format**: `<module-path>/v<version>`

Examples: `cache/inmemory/v0.1.0`, `core/errors/v1.2.3`

**Retry Logic**:

- Maximum retries: 3 (configurable via `MAX_RETRIES`)
- Exponential backoff: 2s, 4s, 8s (configurable via `INITIAL_BACKOFF`)
- Automatic cleanup of local tag on push failure

---


## Library Functions

### `lib/git-utils.sh`

Git utility functions for tag operations and change detection.

**Functions**:

- `get_latest_tag MODULE_PATH`: Get latest tag for a module
  ```bash
  tag=$(get_latest_tag "cache/inmemory")
  # Returns: cache/inmemory/v0.1.0
  ```

- `list_modules`: List all Go modules in repository
  ```bash
  modules=$(list_modules)
  # Returns: newline-separated module paths
  ```

- `get_changed_files BASE_REF HEAD_REF [PATH_FILTER]`: Get changed files between refs
  ```bash
  files=$(get_changed_files "v0.1.0" "HEAD" "cache/inmemory")
  ```

- `get_module_path_from_tag TAG`: Extract module path from tag
  ```bash
  path=$(get_module_path_from_tag "cache/inmemory/v0.1.0")
  # Returns: cache/inmemory
  ```

- `get_commits_for_module MODULE_PATH [BASE_REF] [HEAD_REF]`: Get commits for module
  ```bash
  commits=$(get_commits_for_module "cache/inmemory" "v0.1.0" "HEAD")
  ```

### `lib/semver-utils.sh`

Semantic versioning operations (parse, compare, bump).

**Functions**:

- `parse_version VERSION`: Parse version into major/minor/patch
  ```bash
  read -r major minor patch < <(parse_version "v1.2.3")
  ```

- `compare_versions V1 V2`: Compare two versions (-1, 0, 1)
  ```bash
  result=$(compare_versions "v1.2.3" "v1.3.0")
  # Returns: -1 (v1 < v2)
  ```

- `bump_version VERSION TYPE`: Bump version by type
  ```bash
  new_ver=$(bump_version "v1.2.3" "minor")
  # Returns: v1.3.0
  ```

- `validate_version VERSION`: Check if version is valid
  ```bash
  if validate_version "v1.2.3"; then echo "valid"; fi
  ```

- `get_version_from_tag TAG`: Extract version from module tag
  ```bash
  version=$(get_version_from_tag "cache/inmemory/v1.2.3")
  # Returns: v1.2.3
  ```

### `lib/module-utils.sh`

Module discovery and validation utilities.

**Functions**:

- `discover_modules`: Find all modules and return JSON array
  ```bash
  modules=$(discover_modules)
  # Returns: [{"path":"cache/inmemory","name":"..."}]
  ```

- `validate_module MODULE_PATH`: Check if path contains valid Go module
  ```bash
  if validate_module "cache/inmemory"; then echo "valid"; fi
  ```

- `get_module_name GO_MOD_PATH`: Extract module name from go.mod
  ```bash
  name=$(get_module_name "cache/inmemory/go.mod")
  ```

- `get_go_version MODULE_PATH`: Get Go version from go.mod
  ```bash
  go_version=$(get_go_version "cache/inmemory")
  # Returns: 1.21
  ```

- `has_tests MODULE_PATH`: Check if module has test files
  ```bash
  if has_tests "cache/inmemory"; then echo "has tests"; fi
  ```

- `get_module_display_name MODULE_PATH`: Get display name
  ```bash
  display=$(get_module_display_name "cache/inmemory")
  # Returns: cache-inmemory
  ```

## Testing

Run BATS tests:

```bash
# Install BATS (if not installed)
npm install -g bats

# Run all tests
bats .github/scripts/tests/

# Run specific test file
bats .github/scripts/tests/test-detect-modules.bats
```

## Development Guidelines

### Error Handling

All scripts use:
```bash
set -euo pipefail
```

- `-e`: Exit on error
- `-u`: Error on undefined variables
- `-o pipefail`: Catch errors in pipes

### Input Validation

Always validate inputs:
```bash
if [[ -z "$module_path" ]]; then
    echo "Error: module_path is required" >&2
    return 1
fi
```

### Exit Codes

Use consistent exit codes:
- `0`: Success
- `1`: Invalid arguments
- `2`: Git/filesystem error

### Logging

Use structured logging:
```bash
echo "✅ Success message"
echo "❌ Error message" >&2
echo "::warning title=Title::Warning message"  # GitHub Actions annotation
```

## Dependencies

- Bash 5.x
- Git 2.x
- Go 1.21+ (for module operations)
- Standard Unix utilities (find, sed, awk, grep, bc)

## Troubleshooting

### Permission denied

```bash
chmod +x .github/scripts/*.sh
chmod +x .github/scripts/lib/*.sh
```

### Shellcheck warnings

Run shellcheck to validate scripts:
```bash
shellcheck .github/scripts/*.sh
shellcheck .github/scripts/lib/*.sh
```

### Module not found

Ensure go.mod exists and contains `module` directive.

### Tag not found

Modules without tags are treated as new (all files changed).
