#!/usr/bin/env bash
#
# analyze-commits.sh - Analyze conventional commits for a module
#
# Usage:
#   analyze-commits.sh <module-path> [--from-tag <tag>] [--output json|text]
#
# Description:
#   Parses conventional commits for a specific module directory and determines
#   what type of version bump is needed based on commit types and breaking changes.
#   Only analyzes commits that modified files within the module's directory.
#
# Examples:
#   ./analyze-commits.sh cache/inmemory
#   ./analyze-commits.sh core/errors --from-tag core/errors/v0.1.0 --output json
#
# Exit Codes:
#   0 - Success
#   1 - Invalid arguments or execution error
#   2 - Git operations failed
#   3 - No commits found for module

set -euo pipefail

# Source utility libraries
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=./lib/git-utils.sh
source "${SCRIPT_DIR}/lib/git-utils.sh"

# Configuration
OUTPUT_FORMAT="${OUTPUT_FORMAT:-text}"

# Usage information
usage() {
    cat <<EOF
Usage: $(basename "$0") <module-path> [OPTIONS]

Analyze conventional commits for a specific module to determine version bump type.

Arguments:
  <module-path>           Path to the module directory (e.g., cache/inmemory)

Options:
  --from-tag <tag>        Compare from this tag (default: latest tag for module)
  --output <format>       Output format: json or text (default: text)
  -h, --help             Show this help message

Examples:
  $(basename "$0") cache/inmemory
  $(basename "$0") core/errors --from-tag core/errors/v0.1.0 --output json

Exit Codes:
  0 - Success
  1 - Invalid arguments or execution error
  2 - Git operations failed
  3 - No commits found for module
EOF
}

# Parse commit message to extract type and check for breaking changes
parse_commit_message() {
    local commit_msg="$1"
    local commit_type=""
    local is_breaking=false
    
    # Extract type from conventional commit format
    # Match patterns like: feat:, fix:, feat(scope):, etc.
    if [[ "$commit_msg" =~ ^([a-z]+).*: ]]; then
        commit_type="${BASH_REMATCH[1]}"
    fi
    
    # Check for ! indicating breaking change (before the colon)
    if [[ "$commit_msg" =~ ^[a-z]+.*!: ]]; then
        is_breaking=true
    fi
    
    # Check for BREAKING CHANGE footer
    if echo "$commit_msg" | grep -q "^BREAKING CHANGE:"; then
        is_breaking=true
    fi
    
    echo "${commit_type}|${is_breaking}"
}
    
# Classify commit type into version bump category
classify_commit_type() {
    local commit_type="$1"
    local is_breaking="$2"
    
    if [[ "$is_breaking" == "true" ]]; then
        echo "breaking"
        return
    fi
    
    case "$commit_type" in
        feat|feature)
            echo "minor"
            ;;
        fix)
            echo "patch"
            ;;
        chore|docs|style|refactor|perf|test|build|ci)
            echo "patch"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Analyze commits for a module
analyze_commits() {
    local module_path="$1"
    local from_tag="${2:-}"
    
    # Validate module path exists
    if [[ ! -d "$module_path" ]]; then
        echo "Error: Module directory not found: $module_path" >&2
        return 1
    fi
    
    # Get the comparison point (tag or initial commit)
    local compare_from=""
    if [[ -n "$from_tag" ]]; then
        compare_from="$from_tag"
    else
        # Get latest tag for this module (allow failure)
        compare_from=$(get_latest_tag "$module_path" 2>/dev/null || true)
    fi
    
    # Get commits affecting this module
    # Use custom delimiters to avoid conflicts with commit content
    local commits
    if [[ -n "$compare_from" ]]; then
        commits=$(git log "${compare_from}..HEAD" --format="%H<FIELD_SEP>%s<FIELD_SEP>%b<COMMIT_SEP>" -- "$module_path" 2>/dev/null || true)
    else
        # No tag exists, get all commits for this module
        commits=$(git log --format="%H<FIELD_SEP>%s<FIELD_SEP>%b<COMMIT_SEP>" -- "$module_path" 2>/dev/null || true)
    fi
    
    if [[ -z "$commits" ]]; then
        return 3
    fi
    
    # Analyze each commit
    local has_breaking=false
    local has_feat=false
    local has_fix=false
    local has_unknown=false
    local commit_count=0
    declare -a commit_details=()
    
    # Split commits by separator
    IFS='<COMMIT_SEP>' read -ra commit_array <<< "$commits"
    
    for commit_line in "${commit_array[@]}"; do
        # Skip empty lines
        if [[ -z "$commit_line" ]] || [[ "$commit_line" == $'\n'* ]]; then
            continue
        fi
        
        # Parse fields using field separator
        IFS='<FIELD_SEP>' read -r hash subject body <<< "$commit_line"
        
        # Trim trailing whitespace/newlines from subject
        subject="${subject%$'\n'}"
        subject="${subject%$'\r'}"
        
        # Combine subject and body for full message analysis
        local full_msg="${subject}"$'\n'"${body}"
        
        # Parse commit message
        local parse_result
        parse_result=$(parse_commit_message "$full_msg")
        local commit_type="${parse_result%%|*}"
        local is_breaking="${parse_result#*|}"
        
        # Classify commit
        local bump_type
        bump_type=$(classify_commit_type "$commit_type" "$is_breaking")
        
        # Track what types we've seen
        case "$bump_type" in
            breaking)
                has_breaking=true
                ;;
            minor)
                has_feat=true
                ;;
            patch)
                has_fix=true
                ;;
            unknown)
                has_unknown=true
                ;;
        esac
        
        ((commit_count++))
        commit_details+=("${hash}<FIELD_SEP>${commit_type}<FIELD_SEP>${bump_type}<FIELD_SEP>${subject}")
    done
    
    # Determine overall bump type (highest precedence wins)
    local recommended_bump="patch"
    if [[ "$has_breaking" == true ]]; then
        recommended_bump="major"
    elif [[ "$has_feat" == true ]]; then
        recommended_bump="minor"
    elif [[ "$has_fix" == true ]]; then
        recommended_bump="patch"
    fi
    
    # Build warnings
    local warnings=()
    if [[ "$has_unknown" == true ]]; then
        warnings+=("Some commits do not follow conventional commit format - defaulting to patch bump")
    fi
    
    # Output results
    if [[ "$OUTPUT_FORMAT" == "json" ]]; then
        output_json "$module_path" "$recommended_bump" "$commit_count" \
                    "$has_breaking" "$has_feat" "$has_fix" \
                    commit_details warnings
    else
        output_text "$module_path" "$recommended_bump" "$commit_count" \
                    "$has_breaking" "$has_feat" "$has_fix" \
                    commit_details warnings
    fi
}

# Output analysis results as JSON
output_json() {
    local module_path="$1"
    local recommended_bump="$2"
    local commit_count="$3"
    local has_breaking="$4"
    local has_feat="$5"
    local has_fix="$6"
    local -n commits_ref="$7"
    local -n warnings_ref="$8"
    
    # Build commits array
    local commits_json="["
    local first=true
    for commit in "${commits_ref[@]}"; do
        local hash="${commit%%<FIELD_SEP>*}"
        local rest="${commit#*<FIELD_SEP>}"
        local type="${rest%%<FIELD_SEP>*}"
        rest="${rest#*<FIELD_SEP>}"
        local bump="${rest%%<FIELD_SEP>*}"
        local subject="${rest#*<FIELD_SEP>}"
        
        if [[ "$first" != true ]]; then
            commits_json+=","
        fi
        first=false
        commits_json+="{\"hash\":\"$hash\",\"type\":\"$type\",\"bump\":\"$bump\",\"subject\":$(printf '%s' "$subject" | jq -Rs .)}"
    done
    commits_json+="]"
    
    # Build warnings array
    local warnings_json="["
    first=true
    for warning in "${warnings_ref[@]}"; do
        if [[ "$first" != true ]]; then
            warnings_json+=","
        fi
        first=false
        warnings_json+="$(printf '%s' "$warning" | jq -Rs .)"
    done
    warnings_json+="]"
    
    cat <<EOF
{
  "module": "$module_path",
  "recommended_bump": "$recommended_bump",
  "commit_count": $commit_count,
  "has_breaking": $has_breaking,
  "has_feat": $has_feat,
  "has_fix": $has_fix,
  "commits": $commits_json,
  "warnings": $warnings_json
}
EOF
}

# Output analysis results as text
output_text() {
    local module_path="$1"
    local recommended_bump="$2"
    local commit_count="$3"
    local has_breaking="$4"
    local has_feat="$5"
    local has_fix="$6"
    local -n commits_ref="$7"
    local -n warnings_ref="$8"
    
    echo "Module: $module_path"
    echo "Recommended bump: $recommended_bump"
    echo "Commits analyzed: $commit_count"
    echo "Has breaking changes: $has_breaking"
    echo "Has features: $has_feat"
    echo "Has fixes: $has_fix"
    
    if [[ ${#warnings_ref[@]} -gt 0 ]]; then
        echo ""
        echo "Warnings:"
        for warning in "${warnings_ref[@]}"; do
            echo "  - $warning"
        done
    fi
    
    if [[ ${#commits_ref[@]} -gt 0 ]]; then
        echo ""
        echo "Commits:"
        for commit in "${commits_ref[@]}"; do
            local hash="${commit%%<FIELD_SEP>*}"
            local rest="${commit#*<FIELD_SEP>}"
            local type="${rest%%<FIELD_SEP>*}"
            rest="${rest#*<FIELD_SEP>}"
            local bump="${rest%%<FIELD_SEP>*}"
            local subject="${rest#*<FIELD_SEP>}"
            echo "  [$bump] $hash - $type: $subject"
        done
    fi
}

# Main execution
main() {
    local module_path=""
    local from_tag=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                usage
                exit 0
                ;;
            --from-tag)
                from_tag="$2"
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
    
    # Analyze commits
    analyze_commits "$module_path" "$from_tag"
}

# Run main if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
