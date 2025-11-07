#!/usr/bin/env bash
#
# calculate-version.sh - Calculate next semantic version for a module
#
# Usage:
#   calculate-version.sh <module-path> [--override-version <version>] [--output json|text]
#
# Description:
#   Determines the next semantic version for a module based on conventional commit
#   analysis. Supports manual version override for exceptional cases.
#
# Examples:
#   ./calculate-version.sh cache/inmemory
#   ./calculate-version.sh core/errors --override-version 1.2.0 --output json
#
# Exit Codes:
#   0 - Success
#   1 - Invalid arguments or execution error
#   2 - Version calculation failed
#   3 - Invalid version format

set -euo pipefail

# Source utility libraries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/git-utils.sh
source "${SCRIPT_DIR}/lib/git-utils.sh"
# shellcheck source=./lib/semver-utils.sh
source "${SCRIPT_DIR}/lib/semver-utils.sh"

# Configuration
OUTPUT_FORMAT="${OUTPUT_FORMAT:-text}"

# Usage information
usage() {
    cat <<EOF
Usage: $(basename "$0") <module-path> [OPTIONS]

Calculate the next semantic version for a module based on conventional commits.

Arguments:
  <module-path>           Path to the module directory (e.g., cache/inmemory)

Options:
  --override-version <v>  Manually specify the next version (must be valid semver)
  --output <format>       Output format: json or text (default: text)
  -h, --help             Show this help message

Examples:
  $(basename "$0") cache/inmemory
  $(basename "$0") core/errors --override-version 1.2.0 --output json

Exit Codes:
  0 - Success
  1 - Invalid arguments or execution error
  2 - Version calculation failed
  3 - Invalid version format
EOF
}

# Calculate next version based on bump type
calculate_next_version() {
    local module_path="$1"
    local override_version="${2:-}"
    
    # Validate module path exists
    if [[ ! -d "$module_path" ]]; then
        echo "Error: Module directory not found: $module_path" >&2
        return 1
    fi
    
    # If override version is provided, validate and use it
    if [[ -n "$override_version" ]]; then
        # Remove 'v' prefix if present
        override_version="${override_version#v}"
        
        if ! validate_version "$override_version"; then
            echo "Error: Invalid version format: $override_version (must be MAJOR.MINOR.PATCH)" >&2
            return 3
        fi
        
        # Get current version for comparison
        local current_tag
        current_tag=$(get_latest_tag "$module_path" 2>/dev/null || echo "")
        
        local current_version="0.0.0"
        if [[ -n "$current_tag" ]]; then
            current_version=$(get_version_from_tag "$current_tag")
        fi
        
        # Warn if override version is not greater than current
        if [[ -n "$current_tag" ]]; then
            local comparison
            comparison=$(compare_versions "$override_version" "$current_version")
            if [[ "$comparison" != "greater" ]]; then
                echo "Warning: Override version $override_version is not greater than current version $current_version" >&2
            fi
        fi
        
        output_version_result "$module_path" "$current_version" "$override_version" "manual" "Version manually overridden"
        return 0
    fi
    
    # Get current version from latest tag
    local current_tag
    current_tag=$(get_latest_tag "$module_path" 2>/dev/null || echo "")
    
    local current_version="0.0.0"
    if [[ -n "$current_tag" ]]; then
        current_version=$(get_version_from_tag "$current_tag")
    fi
    
    # Analyze commits to determine bump type
    local analyze_output
    local analyze_exit_code
    
    analyze_output=$("${SCRIPT_DIR}/analyze-commits.sh" "$module_path" --output json 2>&1 || true)
    analyze_exit_code=$?
    
    if [[ $analyze_exit_code -ne 0 ]]; then
        if [[ $analyze_exit_code -eq 3 ]]; then
            # No commits found - no version bump needed
            output_version_result "$module_path" "$current_version" "$current_version" "none" "No commits found for module"
            return 0
        else
            echo "Error: Failed to analyze commits (exit code: $analyze_exit_code)" >&2
            if [[ -n "$analyze_output" ]]; then
                echo "Output: $analyze_output" >&2
            fi
            return 2
        fi
    fi
    
    # Extract recommended bump type from JSON
    local recommended_bump
    recommended_bump=$(echo "$analyze_output" | jq -r '.recommended_bump' 2>/dev/null || echo "")
    
    if [[ -z "$recommended_bump" ]] || [[ "$recommended_bump" == "null" ]]; then
        echo "Error: Failed to extract recommended_bump from analyze output" >&2
        echo "Analyze output was: $analyze_output" >&2
        return 2
    fi
    
    # Check for warnings
    local warnings
    warnings=$(echo "$analyze_output" | jq -r '.warnings[]' 2>/dev/null || echo "")
    
    # Calculate next version
    local next_version
    if [[ -z "$current_tag" ]]; then
        # No existing tag - start at v0.0.1
        next_version="0.0.1"
        recommended_bump="initial"
    else
        next_version=$(bump_version "$current_version" "$recommended_bump")
    fi
    
    # Output result
    local message="Calculated based on conventional commits"
    if [[ -n "$warnings" ]]; then
        message="$message (with warnings)"
    fi
    
    output_version_result "$module_path" "$current_version" "$next_version" "$recommended_bump" "$message"
}

# Output version calculation results
output_version_result() {
    local module_path="$1"
    local current_version="$2"
    local next_version="$3"
    local bump_type="$4"
    local message="$5"
    
    if [[ "$OUTPUT_FORMAT" == "json" ]]; then
        cat <<EOF
{
  "module": "$module_path",
  "current_version": "$current_version",
  "next_version": "$next_version",
  "bump_type": "$bump_type",
  "message": "$message"
}
EOF
    else
        echo "Module: $module_path"
        echo "Current version: $current_version"
        echo "Next version: $next_version"
        echo "Bump type: $bump_type"
        echo "Message: $message"
    fi
}

# Main execution
main() {
    local module_path=""
    local override_version=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                usage
                exit 0
                ;;
            --override-version)
                override_version="$2"
                shift 2
                ;;
            --output)
                OUTPUT_FORMAT="$2"
                shift 2
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
    
    # Validate output format
    if [[ "$OUTPUT_FORMAT" != "json" && "$OUTPUT_FORMAT" != "text" ]]; then
        echo "Error: Invalid output format: $OUTPUT_FORMAT (must be 'json' or 'text')" >&2
        exit 1
    fi
    
    # Calculate version
    calculate_next_version "$module_path" "$override_version"
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
