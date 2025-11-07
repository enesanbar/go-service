#!/usr/bin/env bash
#
# create-tags.sh - Create and push git tags for modules
#
# Usage:
#   create-tags.sh <module-path> <version> [--dry-run] [--push]
#
# Description:
#   Creates directory-prefixed git tags for modules and optionally pushes them
#   to the remote repository. Includes retry logic with exponential backoff for
#   network failures.
#
# Examples:
#   ./create-tags.sh cache/inmemory 0.1.0 --dry-run
#   ./create-tags.sh core/errors 1.2.3 --push
#
# Exit Codes:
#   0 - Success
#   1 - Invalid arguments or execution error
#   2 - Git tag creation failed
#   3 - Git push failed after retries

set -euo pipefail

# Source utility libraries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/semver-utils.sh
source "${SCRIPT_DIR}/lib/semver-utils.sh"

# Configuration
DRY_RUN=false
PUSH_TAGS=false
MAX_RETRIES=3
INITIAL_BACKOFF=2

# Usage information
usage() {
    cat <<EOF
Usage: $(basename "$0") <module-path> <version> [OPTIONS]

Create and push git tags for a module with directory-level prefix.

Arguments:
  <module-path>           Path to the module directory (e.g., cache/inmemory)
  <version>              Semantic version without 'v' prefix (e.g., 1.2.3)

Options:
  --dry-run              Show what would be done without making changes
  --push                 Push tags to remote repository (default: false)
  -h, --help             Show this help message

Examples:
  $(basename "$0") cache/inmemory 0.1.0 --dry-run
  $(basename "$0") core/errors 1.2.3 --push

Exit Codes:
  0 - Success
  1 - Invalid arguments or execution error
  2 - Git tag creation failed
  3 - Git push failed after retries

Environment Variables:
  MAX_RETRIES            Maximum number of push retry attempts (default: 3)
  INITIAL_BACKOFF        Initial backoff delay in seconds (default: 2)
EOF
}

# Create a git tag with annotation
create_tag() {
    local module_path="$1"
    local version="$2"
    
    # Validate version format
    if ! validate_version "$version"; then
        echo "Error: Invalid version format: $version (must be MAJOR.MINOR.PATCH)" >&2
        return 1
    fi
    
    # Construct tag name with module path prefix
    local tag_name="${module_path}/v${version}"
    
    # Check if tag already exists
    if git rev-parse "$tag_name" >/dev/null 2>&1; then
        echo "Error: Tag $tag_name already exists" >&2
        return 2
    fi
    
    # Create annotated tag
    local tag_message="Release ${module_path} version ${version}"
    
    if [[ "$DRY_RUN" == true ]]; then
        echo "[DRY RUN] Would create tag: $tag_name"
        echo "[DRY RUN] Tag message: $tag_message"
        return 0
    fi
    
    if ! git tag -a "$tag_name" -m "$tag_message"; then
        echo "Error: Failed to create tag $tag_name" >&2
        return 2
    fi
    
    echo "✓ Created tag: $tag_name" >&2
    echo "$tag_name"  # Output tag name for use by caller
}

# Push tags to remote with retry logic
push_tag() {
    local tag_name="$1"
    local attempt=1
    local backoff="$INITIAL_BACKOFF"
    
    if [[ "$DRY_RUN" == true ]]; then
        echo "[DRY RUN] Would push tag: $tag_name"
        return 0
    fi
    
    while [[ $attempt -le $MAX_RETRIES ]]; do
        echo "Pushing tag $tag_name (attempt $attempt/$MAX_RETRIES)..."
        
        if git push origin "refs/tags/$tag_name"; then
            echo "✓ Successfully pushed tag: $tag_name"
            return 0
        fi
        
        local exit_code=$?
        
        if [[ $attempt -lt $MAX_RETRIES ]]; then
            echo "⚠ Push failed (exit code: $exit_code). Retrying in ${backoff}s..." >&2
            sleep "$backoff"
            backoff=$((backoff * 2))  # Exponential backoff
            ((attempt++))
        else
            echo "Error: Failed to push tag $tag_name after $MAX_RETRIES attempts" >&2
            return 3
        fi
    done
}

# Delete a local tag (used for cleanup on failure)
delete_local_tag() {
    local tag_name="$1"
    
    if [[ "$DRY_RUN" == true ]]; then
        return 0
    fi
    
    if git rev-parse "$tag_name" >/dev/null 2>&1; then
        echo "Cleaning up local tag: $tag_name" >&2
        git tag -d "$tag_name" >/dev/null 2>&1 || true
    fi
}

# Main execution
main() {
    local module_path=""
    local version=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                usage
                exit 0
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --push)
                PUSH_TAGS=true
                shift
                ;;
            -*)
                echo "Error: Unknown option: $1" >&2
                usage
                exit 1
                ;;
            *)
                if [[ -z "$module_path" ]]; then
                    module_path="$1"
                    shift
                elif [[ -z "$version" ]]; then
                    version="$1"
                    shift
                else
                    echo "Error: Unexpected argument: $1" >&2
                    usage
                    exit 1
                fi
                ;;
        esac
    done
    
    # Validate required arguments
    if [[ -z "$module_path" ]]; then
        echo "Error: Module path is required" >&2
        usage
        exit 1
    fi
    
    if [[ -z "$version" ]]; then
        echo "Error: Version is required" >&2
        usage
        exit 1
    fi
    
    # Remove 'v' prefix if present
    version="${version#v}"
    
    # Validate module directory exists
    if [[ ! -d "$module_path" ]]; then
        echo "Error: Module directory not found: $module_path" >&2
        exit 1
    fi
    
    # Create tag
    local tag_name
    tag_name=$(create_tag "$module_path" "$version") || exit $?
    
    # Push tag if requested
    if [[ "$PUSH_TAGS" == true ]]; then
        if ! push_tag "$tag_name"; then
            # Cleanup local tag on push failure
            delete_local_tag "$tag_name"
            exit 3
        fi
    else
        echo "ℹ Tag created locally. Use --push to push to remote."
    fi
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
