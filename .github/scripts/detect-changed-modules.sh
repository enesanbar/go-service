#!/usr/bin/env bash
# detect-changed-modules.sh - Detect which Go modules have changed
# Purpose: Compare current state against last tagged version per module
# Usage: ./detect-changed-modules.sh [--base-ref REF] [--head-ref REF] [--format json|text] [--include-all]

set -euo pipefail

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source utility libraries
# shellcheck source=lib/git-utils.sh
source "${SCRIPT_DIR}/lib/git-utils.sh"
# shellcheck source=lib/module-utils.sh
source "${SCRIPT_DIR}/lib/module-utils.sh"

# Default values
BASE_REF=""
HEAD_REF="HEAD"
OUTPUT_FORMAT="json"
INCLUDE_ALL="false"

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --base-ref)
            BASE_REF="$2"
            shift 2
            ;;
        --head-ref)
            HEAD_REF="$2"
            shift 2
            ;;
        --format)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        --include-all)
            INCLUDE_ALL="true"
            shift
            ;;
        -h|--help)
            cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Detect which Go modules have changed by comparing against their last tagged version.

OPTIONS:
    --base-ref REF      Base git reference for comparison (default: auto-detect from tags)
    --head-ref REF      Head git reference (default: HEAD)
    --format FORMAT     Output format: json or text (default: json)
    --include-all       Include all modules even if unchanged (default: false)
    -h, --help          Show this help message

EXAMPLES:
    # Auto-detect changes using module tags
    $(basename "$0")
    
    # Compare against specific ref
    $(basename "$0") --base-ref main --head-ref feature-branch
    
    # Text output
    $(basename "$0") --format text
    
EXIT CODES:
    0 - Success
    1 - Invalid arguments
    2 - Git operation error
EOF
            exit 0
            ;;
        *)
            echo "Error: Unknown option: $1" >&2
            echo "Run with --help for usage information" >&2
            exit 1
            ;;
    esac
done

# Validate output format
if [[ "$OUTPUT_FORMAT" != "json" ]] && [[ "$OUTPUT_FORMAT" != "text" ]]; then
    echo "Error: Invalid format '$OUTPUT_FORMAT'. Use 'json' or 'text'" >&2
    exit 1
fi

# Array to store module change information
declare -a CHANGED_MODULES=()
declare -a ALL_MODULES=()

# Discover all modules in repository
while IFS= read -r module_path; do
    if [[ -z "$module_path" ]]; then
        continue
    fi
    
    # Validate module
    if ! validate_module "$module_path" 2>/dev/null; then
        continue
    fi
    
    # Get module name
    module_name=$(get_module_name "${module_path}/go.mod" 2>/dev/null || echo "unknown")
    
    # Get last tag for this module
    last_tag=$(get_latest_tag "$module_path" 2>/dev/null || echo "")
    
    # Determine comparison base
    local_base_ref="$BASE_REF"
    if [[ -z "$local_base_ref" ]] && [[ -n "$last_tag" ]]; then
        local_base_ref="$last_tag"
    fi
    
    # Check for changes
    has_changes="false"
    changed_files=""
    
    if [[ -z "$local_base_ref" ]]; then
        # No tag exists - treat as new module (all files changed)
        has_changes="true"
        changed_files=$(find "$module_path" -type f -not -path "*/vendor/*" | head -n 100 | tr '\n' ',' | sed 's/,$//')
    else
        # Get changed files for this module
        changed_files_list=$(get_changed_files "$local_base_ref" "$HEAD_REF" "$module_path" 2>/dev/null || echo "")
        
        if [[ -n "$changed_files_list" ]]; then
            has_changes="true"
            changed_files=$(echo "$changed_files_list" | tr '\n' ',' | sed 's/,$//')
        fi
    fi
    
    # Build module info JSON
    module_info=$(cat <<EOF
{
  "path": "$module_path",
  "name": "$module_name",
  "last_tag": "$last_tag",
  "has_changes": $has_changes,
  "changed_files": "$changed_files"
}
EOF
)
    
    # Add to appropriate array
    ALL_MODULES+=("$module_info")
    if [[ "$has_changes" == "true" ]]; then
        CHANGED_MODULES+=("$module_info")
    fi
    
done < <(list_modules)

# Check for non-module file changes (e.g., root README, CI config)
# If any non-module files changed, return all modules
if [[ -z "$BASE_REF" ]]; then
    # Use the oldest tag across all modules as baseline for non-module check
    BASE_REF=$(git tag -l "*/v*" | sort -V | head -n 1 || echo "")
fi

if [[ -n "$BASE_REF" ]]; then
    # Get all changed files
    all_changed=$(get_changed_files "$BASE_REF" "$HEAD_REF" "" 2>/dev/null || echo "")
    
    # Check if any changed files are outside module directories
    non_module_changes="false"
    while IFS= read -r file; do
        if [[ -z "$file" ]]; then
            continue
        fi
        
        # Check if file is in any module directory
        in_module="false"
        while IFS= read -r module_path; do
            if [[ "$file" == "$module_path"/* ]]; then
                in_module="true"
                break
            fi
        done < <(list_modules)
        
        if [[ "$in_module" == "false" ]]; then
            non_module_changes="true"
            break
        fi
    done <<< "$all_changed"
    
    # If non-module files changed, mark all modules as changed
    if [[ "$non_module_changes" == "true" ]]; then
        CHANGED_MODULES=("${ALL_MODULES[@]}")
    fi
fi

# Output results
if [[ "$OUTPUT_FORMAT" == "json" ]]; then
    # Determine which array to output
    if [[ "$INCLUDE_ALL" == "true" ]]; then
        output_array=("${ALL_MODULES[@]}")
    else
        output_array=("${CHANGED_MODULES[@]}")
    fi
    
    # Build compact JSON output on a single line
    if [[ ${#output_array[@]} -eq 0 ]]; then
        echo "[]"
    else
        # Use jq to create proper compact JSON array
        printf '%s\n' "${output_array[@]}" | jq -cs '.'
    fi
else
    # Text output
    if [[ "$INCLUDE_ALL" == "true" ]]; then
        for module_info in "${ALL_MODULES[@]}"; do
            path=$(echo "$module_info" | grep -o '"path": "[^"]*"' | cut -d'"' -f4)
            echo "$path"
        done
    else
        for module_info in "${CHANGED_MODULES[@]}"; do
            path=$(echo "$module_info" | grep -o '"path": "[^"]*"' | cut -d'"' -f4)
            echo "$path"
        done
    fi
fi

# Exit successfully
exit 0
