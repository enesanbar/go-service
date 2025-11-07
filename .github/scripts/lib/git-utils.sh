#!/usr/bin/env bash
# git-utils.sh - Git utility functions for CI/CD pipeline
# Functions: get_latest_tag, list_modules, get_changed_files

set -euo pipefail

# get_latest_tag: Get the latest tag for a specific module path
# Args:
#   $1 - module_path: Path to the module (e.g., "cache/inmemory")
# Returns:
#   Latest tag name for the module, or empty string if no tags exist
# Example:
#   tag=$(get_latest_tag "cache/inmemory")  # Returns "cache/inmemory/v0.1.0"
get_latest_tag() {
    local module_path="${1:-}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    # List all tags matching the module path pattern, sort by version, get latest
    git tag -l "${module_path}/v*" \
        | sort -V \
        | tail -n 1 \
        || echo ""
}

# list_modules: Discover all Go modules in the repository
# Returns:
#   Newline-separated list of module paths (relative to repo root)
# Example:
#   modules=$(list_modules)
list_modules() {
    # Find all go.mod files, exclude vendor directories, extract directory paths
    find . -name "go.mod" -not -path "*/vendor/*" -not -path "*/.git/*" \
        | sed 's|^\./||' \
        | sed 's|/go.mod$||' \
        | sort
}

# get_changed_files: Get list of files that changed between two refs
# Args:
#   $1 - base_ref: Base git reference (commit, tag, branch)
#   $2 - head_ref: Head git reference (defaults to HEAD if not provided)
#   $3 - path_filter: Optional path filter (e.g., "cache/inmemory")
# Returns:
#   Newline-separated list of changed file paths
# Example:
#   files=$(get_changed_files "cache/inmemory/v0.1.0" "HEAD" "cache/inmemory")
get_changed_files() {
    local base_ref="${1:-}"
    local head_ref="${2:-HEAD}"
    local path_filter="${3:-}"
    
    if [[ -z "$base_ref" ]]; then
        echo "Error: base_ref is required" >&2
        return 1
    fi
    
    # Check if refs exist
    if ! git rev-parse --verify "$base_ref" >/dev/null 2>&1; then
        echo "Error: base_ref '$base_ref' does not exist" >&2
        return 1
    fi
    
    if ! git rev-parse --verify "$head_ref" >/dev/null 2>&1; then
        echo "Error: head_ref '$head_ref' does not exist" >&2
        return 1
    fi
    
    # Get changed files
    if [[ -n "$path_filter" ]]; then
        git diff --name-only "$base_ref".."$head_ref" -- "$path_filter" || echo ""
    else
        git diff --name-only "$base_ref".."$head_ref" || echo ""
    fi
}

# get_module_path_from_tag: Extract module path from a tag name
# Args:
#   $1 - tag: Tag name (e.g., "cache/inmemory/v0.1.0")
# Returns:
#   Module path (e.g., "cache/inmemory")
# Example:
#   path=$(get_module_path_from_tag "cache/inmemory/v0.1.0")
get_module_path_from_tag() {
    local tag="${1:-}"
    
    if [[ -z "$tag" ]]; then
        echo "Error: tag is required" >&2
        return 1
    fi
    
    # Remove /vX.Y.Z suffix
    echo "$tag" | sed 's|/v[0-9].*$||'
}

# get_commits_for_module: Get commits that affected a specific module path
# Args:
#   $1 - module_path: Path to the module
#   $2 - base_ref: Base git reference (optional, defaults to first commit)
#   $3 - head_ref: Head git reference (optional, defaults to HEAD)
# Returns:
#   Newline-separated list of commit hashes
# Example:
#   commits=$(get_commits_for_module "cache/inmemory" "cache/inmemory/v0.1.0" "HEAD")
get_commits_for_module() {
    local module_path="${1:-}"
    local base_ref="${2:-}"
    local head_ref="${3:-HEAD}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    # Build git log command
    local git_cmd="git log --format=%H"
    
    if [[ -n "$base_ref" ]]; then
        if ! git rev-parse --verify "$base_ref" >/dev/null 2>&1; then
            echo "Error: base_ref '$base_ref' does not exist" >&2
            return 1
        fi
        git_cmd="$git_cmd ${base_ref}..${head_ref}"
    else
        git_cmd="$git_cmd ${head_ref}"
    fi
    
    git_cmd="$git_cmd -- ${module_path}"
    
    eval "$git_cmd" || echo ""
}

# Export functions for use in other scripts
export -f get_latest_tag
export -f list_modules
export -f get_changed_files
export -f get_module_path_from_tag
export -f get_commits_for_module
