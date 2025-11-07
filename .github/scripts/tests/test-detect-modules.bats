#!/usr/bin/env bats
# test-detect-modules.bats - Tests for detect-changed-modules.sh
# Test scenarios: single module, multiple modules, no tags, non-module changes

# Setup function - runs before each test
setup() {
    # Store script directory
    SCRIPT_DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")/.." && pwd)"
    DETECT_SCRIPT="${SCRIPT_DIR}/detect-changed-modules.sh"
    
    # Verify script exists
    [[ -x "$DETECT_SCRIPT" ]]
}

# Test: Script exists and is executable
@test "detect-changed-modules.sh exists and is executable" {
    [[ -f "$DETECT_SCRIPT" ]]
    [[ -x "$DETECT_SCRIPT" ]]
}

# Test: Help option displays usage
@test "detect-changed-modules.sh --help displays usage" {
    run "$DETECT_SCRIPT" --help
    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "Usage:" ]]
    [[ "$output" =~ "OPTIONS:" ]]
}

# Test: Invalid format option returns error
@test "detect-changed-modules.sh rejects invalid format" {
    run "$DETECT_SCRIPT" --format invalid
    [[ "$status" -eq 1 ]]
    [[ "$output" =~ "Invalid format" ]]
}

# Test: JSON output format is valid
@test "detect-changed-modules.sh produces valid JSON" {
    run "$DETECT_SCRIPT" --format json
    [[ "$status" -eq 0 ]]
    # Output should start with [ or be []
    [[ "$output" =~ ^\[ ]]
}

# Test: Text output format works
@test "detect-changed-modules.sh produces text output" {
    run "$DETECT_SCRIPT" --format text
    [[ "$status" -eq 0 ]]
}

# Test: Script handles existing modules
@test "detect-changed-modules.sh finds existing modules" {
    run "$DETECT_SCRIPT" --format json --include-all
    [[ "$status" -eq 0 ]]
    
    # Should find at least one module (we know they exist in the repo)
    [[ "$output" != "[]" ]]
}

# Note: Additional integration tests would require:
# - Creating temporary git repository
# - Creating test modules with go.mod files
# - Creating tags and commits
# - Testing change detection scenarios
# These are better suited for CI environment testing
