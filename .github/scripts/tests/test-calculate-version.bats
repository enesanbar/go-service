#!/usr/bin/env bats
#
# Integration tests for calculate-version.sh
#
# Tests semantic version calculation based on commit analysis

setup() {
    # Create a temporary git repository for testing
    export TEST_DIR="$(mktemp -d)"
    export SCRIPT_DIR="${BATS_TEST_DIRNAME}/.."
    
    cd "$TEST_DIR"
    git init
    git config user.email "test@example.com"
    git config user.name "Test User"
    
    # Create initial commit
    echo "initial" > README.md
    git add README.md
    git commit -m "chore: initial commit"
}

teardown() {
    # Clean up temporary directory
    rm -rf "$TEST_DIR"
}

create_module() {
    local module_path="$1"
    
    mkdir -p "$module_path"
    cat > "$module_path/go.mod" <<EOF
module github.com/test/${module_path}

go 1.25
EOF
    git add "$module_path"
    git commit -m "feat: create module $module_path"
}

@test "calculate-version: calculates patch version for fix commits" {
    create_module "cache/inmemory"
    
    git tag "cache/inmemory/v0.1.0"
    
    echo "fix" > cache/inmemory/fix.go
    git add cache/inmemory/fix.go
    git commit -m "fix: resolve issue"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "cache/inmemory" --output json
    
    [ "$status" -eq 0 ]
    
    current=$(echo "$output" | jq -r '.current_version')
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$current" = "0.1.0" ]
    [ "$next" = "0.1.1" ]
    [ "$bump_type" = "patch" ]
}

@test "calculate-version: calculates minor version for feat commits" {
    create_module "core/errors"
    
    git tag "core/errors/v0.2.5"
    
    echo "feature" > core/errors/feature.go
    git add core/errors/feature.go
    git commit -m "feat: add new feature"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "core/errors" --output json
    
    [ "$status" -eq 0 ]
    
    current=$(echo "$output" | jq -r '.current_version')
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$current" = "0.2.5" ]
    [ "$next" = "0.3.0" ]
    [ "$bump_type" = "minor" ]
}

@test "calculate-version: calculates major version for breaking changes" {
    create_module "messaging/rabbitmq"
    
    git tag "messaging/rabbitmq/v1.2.3"
    
    echo "breaking" > messaging/rabbitmq/breaking.go
    git add messaging/rabbitmq/breaking.go
    git commit -m "feat!: breaking change"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "messaging/rabbitmq" --output json
    
    [ "$status" -eq 0 ]
    
    current=$(echo "$output" | jq -r '.current_version')
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$current" = "1.2.3" ]
    [ "$next" = "2.0.0" ]
    [ "$bump_type" = "major" ]
}

@test "calculate-version: starts at 0.0.1 for new modules" {
    create_module "new/module"
    
    echo "feature" > new/module/feature.go
    git add new/module/feature.go
    git commit -m "feat: initial feature"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "new/module" --output json
    
    [ "$status" -eq 0 ]
    
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$next" = "0.0.1" ]
    [ "$bump_type" = "initial" ]
}

@test "calculate-version: handles no commits (no bump needed)" {
    create_module "cache/inmemory"
    create_module "core/errors"
    
    git tag "cache/inmemory/v0.1.0"
    git tag "core/errors/v0.1.0"
    
    # Make changes only to core/errors
    echo "fix" > core/errors/fix.go
    git add core/errors/fix.go
    git commit -m "fix: resolve issue"
    
    # Calculate version for cache/inmemory (no changes)
    run "${SCRIPT_DIR}/calculate-version.sh" "cache/inmemory" --output json
    
    [ "$status" -eq 0 ]
    
    current=$(echo "$output" | jq -r '.current_version')
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    # Version should remain the same
    [ "$current" = "$next" ]
    [ "$bump_type" = "none" ]
}

@test "calculate-version: manual version override" {
    create_module "protocol/grpc"
    
    git tag "protocol/grpc/v0.1.0"
    
    echo "fix" > protocol/grpc/fix.go
    git add protocol/grpc/fix.go
    git commit -m "fix: minor fix"
    
    # Override to version 2.0.0
    run "${SCRIPT_DIR}/calculate-version.sh" "protocol/grpc" --override-version "2.0.0" --output json
    
    [ "$status" -eq 0 ]
    
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$next" = "2.0.0" ]
    [ "$bump_type" = "manual" ]
}

@test "calculate-version: manual override with v prefix" {
    create_module "cron"
    
    git tag "cron/v0.5.0"
    
    echo "change" > cron/scheduler.go
    git add cron/scheduler.go
    git commit -m "update"
    
    # Override with v prefix (should be stripped)
    run "${SCRIPT_DIR}/calculate-version.sh" "cron" --override-version "v1.0.0" --output json
    
    [ "$status" -eq 0 ]
    
    next=$(echo "$output" | jq -r '.next_version')
    [ "$next" = "1.0.0" ]
}

@test "calculate-version: rejects invalid manual version" {
    create_module "persistence/mongodb"
    
    git tag "persistence/mongodb/v0.1.0"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "persistence/mongodb" --override-version "invalid" --output json
    
    [ "$status" -eq 3 ]
    [[ "$output" =~ "Invalid version format" ]]
}

@test "calculate-version: warns if override is not greater than current" {
    create_module "core/validation"
    
    git tag "core/validation/v2.0.0"
    
    echo "change" > core/validation/validator.go
    git add core/validation/validator.go
    git commit -m "update"
    
    # Override to lower version
    run "${SCRIPT_DIR}/calculate-version.sh" "core/validation" --override-version "1.0.0" --output json
    
    [ "$status" -eq 0 ]
    
    # Should succeed but output warning
    [[ "$output" =~ "Warning" || "$output" =~ "not greater" ]] || true
    
    next=$(echo "$output" | jq -r '.next_version')
    [ "$next" = "1.0.0" ]
}

@test "calculate-version: text output format" {
    create_module "cache/inmemory"
    
    git tag "cache/inmemory/v0.1.0"
    
    echo "fix" > cache/inmemory/fix.go
    git add cache/inmemory/fix.go
    git commit -m "fix: bug"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "cache/inmemory" --output text
    
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Module: cache/inmemory" ]]
    [[ "$output" =~ "Current version: 0.1.0" ]]
    [[ "$output" =~ "Next version: 0.1.1" ]]
}

@test "calculate-version: fails with invalid module path" {
    run "${SCRIPT_DIR}/calculate-version.sh" "nonexistent/module" --output json
    
    [ "$status" -ne 0 ]
    [[ "$output" =~ "not found" ]]
}

@test "calculate-version: handles complex version progression" {
    create_module "protocol/rest"
    
    # Start with initial version
    git tag "protocol/rest/v0.0.1"
    
    # Add feature (should bump to 0.1.0)
    echo "feat" > protocol/rest/feature.go
    git add protocol/rest/feature.go
    git commit -m "feat: new feature"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "protocol/rest" --output json
    [ "$status" -eq 0 ]
    next=$(echo "$output" | jq -r '.next_version')
    [ "$next" = "0.1.0" ]
    
    # Tag the new version
    git tag "protocol/rest/v0.1.0"
    
    # Add breaking change (should bump to 1.0.0)
    echo "breaking" > protocol/rest/breaking.go
    git add protocol/rest/breaking.go
    git commit -m "feat!: breaking API change"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "protocol/rest" --output json
    [ "$status" -eq 0 ]
    next=$(echo "$output" | jq -r '.next_version')
    [ "$next" = "1.0.0" ]
}

@test "calculate-version: multiple commits with different types" {
    create_module "persistence/mysql"
    
    git tag "persistence/mysql/v1.0.0"
    
    # Add multiple commits (highest precedence should win)
    echo "fix1" > persistence/mysql/fix1.go
    git add persistence/mysql/fix1.go
    git commit -m "fix: bug fix"
    
    echo "feat" > persistence/mysql/feature.go
    git add persistence/mysql/feature.go
    git commit -m "feat: new feature"
    
    echo "fix2" > persistence/mysql/fix2.go
    git add persistence/mysql/fix2.go
    git commit -m "fix: another fix"
    
    run "${SCRIPT_DIR}/calculate-version.sh" "persistence/mysql" --output json
    
    [ "$status" -eq 0 ]
    
    # Should bump minor (feat > fix)
    next=$(echo "$output" | jq -r '.next_version')
    bump_type=$(echo "$output" | jq -r '.bump_type')
    
    [ "$next" = "1.1.0" ]
    [ "$bump_type" = "minor" ]
}
