#!/usr/bin/env bats
#
# Integration tests for analyze-commits.sh
#
# Tests conventional commit parsing and version bump recommendations

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

@test "analyze-commits: detects fix commits as patch bump" {
    create_module "cache/inmemory"
    
    # Tag the module
    git tag "cache/inmemory/v0.1.0"
    
    # Make a fix commit
    echo "fix" > cache/inmemory/fix.go
    git add cache/inmemory/fix.go
    git commit -m "fix: resolve cache eviction bug"
    
    # Analyze commits
    run "${SCRIPT_DIR}/analyze-commits.sh" "cache/inmemory" --output json
    
    [ "$status" -eq 0 ]
    
    # Check recommended bump is patch
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "patch" ]
    
    # Check has_fix is true
    has_fix=$(echo "$output" | jq -r '.has_fix')
    [ "$has_fix" = "true" ]
}

@test "analyze-commits: detects feat commits as minor bump" {
    create_module "core/errors"
    
    git tag "core/errors/v0.1.0"
    
    # Make a feature commit
    echo "feature" > core/errors/feature.go
    git add core/errors/feature.go
    git commit -m "feat: add new error type"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "core/errors" --output json
    
    [ "$status" -eq 0 ]
    
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "minor" ]
    
    has_feat=$(echo "$output" | jq -r '.has_feat')
    [ "$has_feat" = "true" ]
}

@test "analyze-commits: detects breaking changes with ! suffix" {
    create_module "messaging/rabbitmq"
    
    git tag "messaging/rabbitmq/v0.5.0"
    
    # Make a breaking change commit with ! suffix
    echo "breaking" > messaging/rabbitmq/breaking.go
    git add messaging/rabbitmq/breaking.go
    git commit -m "feat!: change API signature"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "messaging/rabbitmq" --output json
    
    [ "$status" -eq 0 ]
    
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "major" ]
    
    has_breaking=$(echo "$output" | jq -r '.has_breaking')
    [ "$has_breaking" = "true" ]
}

@test "analyze-commits: detects breaking changes with BREAKING CHANGE footer" {
    create_module "protocol/grpc"
    
    git tag "protocol/grpc/v1.0.0"
    
    # Make a breaking change commit with footer
    echo "breaking" > protocol/grpc/breaking.go
    git add protocol/grpc/breaking.go
    git commit -m "refactor: restructure API

BREAKING CHANGE: removed deprecated methods"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "protocol/grpc" --output json
    
    [ "$status" -eq 0 ]
    
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "major" ]
    
    has_breaking=$(echo "$output" | jq -r '.has_breaking')
    [ "$has_breaking" = "true" ]
}

@test "analyze-commits: handles multiple commit types (breaking takes precedence)" {
    create_module "persistence/mongodb"
    
    git tag "persistence/mongodb/v0.2.0"
    
    # Add multiple commits
    echo "fix" > persistence/mongodb/fix.go
    git add persistence/mongodb/fix.go
    git commit -m "fix: resolve connection leak"
    
    echo "feat" > persistence/mongodb/feature.go
    git add persistence/mongodb/feature.go
    git commit -m "feat: add connection pooling"
    
    echo "breaking" > persistence/mongodb/breaking.go
    git add persistence/mongodb/breaking.go
    git commit -m "feat!: change connection interface"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "persistence/mongodb" --output json
    
    [ "$status" -eq 0 ]
    
    # Breaking should take precedence
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "major" ]
    
    # All flags should be true
    has_breaking=$(echo "$output" | jq -r '.has_breaking')
    has_feat=$(echo "$output" | jq -r '.has_feat')
    has_fix=$(echo "$output" | jq -r '.has_fix')
    [ "$has_breaking" = "true" ]
    [ "$has_feat" = "true" ]
    [ "$has_fix" = "true" ]
}

@test "analyze-commits: only analyzes commits affecting module directory" {
    create_module "cache/inmemory"
    create_module "core/errors"
    
    git tag "cache/inmemory/v0.1.0"
    git tag "core/errors/v0.1.0"
    
    # Make changes to core/errors only
    echo "fix" > core/errors/fix.go
    git add core/errors/fix.go
    git commit -m "fix: resolve error formatting"
    
    # Analyze cache/inmemory (should find no commits)
    run "${SCRIPT_DIR}/analyze-commits.sh" "cache/inmemory" --output json
    
    # Should exit with code 3 (no commits found)
    [ "$status" -eq 3 ]
}

@test "analyze-commits: handles non-conventional commits with warning" {
    create_module "cron"
    
    git tag "cron/v0.1.0"
    
    # Make a non-conventional commit
    echo "change" > cron/scheduler.go
    git add cron/scheduler.go
    git commit -m "update scheduler logic"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "cron" --output json
    
    [ "$status" -eq 0 ]
    
    # Should default to patch with warning
    recommended_bump=$(echo "$output" | jq -r '.recommended_bump')
    [ "$recommended_bump" = "patch" ]
    
    warnings=$(echo "$output" | jq -r '.warnings[]' | wc -l)
    [ "$warnings" -gt 0 ]
}

@test "analyze-commits: works with no existing tags" {
    create_module "new/module"
    
    # Make a commit without any tags
    echo "feature" > new/module/feature.go
    git add new/module/feature.go
    git commit -m "feat: add initial feature"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "new/module" --output json
    
    [ "$status" -eq 0 ]
    
    # Should analyze all commits since module creation
    commit_count=$(echo "$output" | jq -r '.commit_count')
    [ "$commit_count" -ge 1 ]
}

@test "analyze-commits: text output format" {
    create_module "cache/inmemory"
    
    git tag "cache/inmemory/v0.1.0"
    
    echo "fix" > cache/inmemory/fix.go
    git add cache/inmemory/fix.go
    git commit -m "fix: resolve issue"
    
    run "${SCRIPT_DIR}/analyze-commits.sh" "cache/inmemory" --output text
    
    [ "$status" -eq 0 ]
    [[ "$output" =~ "Module: cache/inmemory" ]]
    [[ "$output" =~ "Recommended bump: patch" ]]
}

@test "analyze-commits: fails with invalid module path" {
    run "${SCRIPT_DIR}/analyze-commits.sh" "nonexistent/module" --output json
    
    [ "$status" -ne 0 ]
    [[ "$output" =~ "not found" ]]
}

@test "analyze-commits: supports --from-tag option" {
    create_module "core/validation"
    
    # Create first tag
    git tag "core/validation/v0.1.0"
    
    echo "fix" > core/validation/fix1.go
    git add core/validation/fix1.go
    git commit -m "fix: first fix"
    
    # Create second tag
    git tag "core/validation/v0.1.1"
    
    echo "fix" > core/validation/fix2.go
    git add core/validation/fix2.go
    git commit -m "fix: second fix"
    
    # Analyze from specific tag
    run "${SCRIPT_DIR}/analyze-commits.sh" "core/validation" --from-tag "core/validation/v0.1.1" --output json
    
    [ "$status" -eq 0 ]
    
    # Should only find one commit (after v0.1.1)
    commit_count=$(echo "$output" | jq -r '.commit_count')
    [ "$commit_count" -eq 1 ]
}
