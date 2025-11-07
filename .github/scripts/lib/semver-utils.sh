#!/usr/bin/env bash
# semver-utils.sh - Semantic versioning utility functions
# Functions: parse_version, compare_versions, bump_version

set -euo pipefail

# parse_version: Parse a semantic version string into components
# Args:
#   $1 - version: Version string (e.g., "v1.2.3" or "1.2.3")
# Returns:
#   Space-separated major minor patch (e.g., "1 2 3")
# Example:
#   read -r major minor patch < <(parse_version "v1.2.3")
parse_version() {
    local version="${1:-}"
    
    if [[ -z "$version" ]]; then
        echo "Error: version is required" >&2
        return 1
    fi
    
    # Remove 'v' prefix if present
    version="${version#v}"
    
    # Validate format
    if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: invalid semantic version format: $version (expected: MAJOR.MINOR.PATCH)" >&2
        return 1
    fi
    
    # Split into components
    local major minor patch
    IFS='.' read -r major minor patch <<< "$version"
    
    echo "$major $minor $patch"
}

# compare_versions: Compare two semantic versions
# Args:
#   $1 - version1: First version (e.g., "v1.2.3")
#   $2 - version2: Second version (e.g., "v1.3.0")
# Returns:
#   -1 if version1 < version2
#    0 if version1 = version2
#    1 if version1 > version2
# Example:
#   result=$(compare_versions "v1.2.3" "v1.3.0")  # Returns -1
compare_versions() {
    local version1="${1:-}"
    local version2="${2:-}"
    
    if [[ -z "$version1" ]] || [[ -z "$version2" ]]; then
        echo "Error: both versions are required" >&2
        return 1
    fi
    
    # Parse versions
    local major1 minor1 patch1 major2 minor2 patch2
    read -r major1 minor1 patch1 < <(parse_version "$version1")
    read -r major2 minor2 patch2 < <(parse_version "$version2")
    
    # Compare major
    if [[ "$major1" -lt "$major2" ]]; then
        echo "-1"
        return 0
    elif [[ "$major1" -gt "$major2" ]]; then
        echo "1"
        return 0
    fi
    
    # Compare minor
    if [[ "$minor1" -lt "$minor2" ]]; then
        echo "-1"
        return 0
    elif [[ "$minor1" -gt "$minor2" ]]; then
        echo "1"
        return 0
    fi
    
    # Compare patch
    if [[ "$patch1" -lt "$patch2" ]]; then
        echo "-1"
        return 0
    elif [[ "$patch1" -gt "$patch2" ]]; then
        echo "1"
        return 0
    fi
    
    # Equal
    echo "0"
}

# bump_version: Bump a semantic version by type
# Args:
#   $1 - version: Current version (e.g., "v1.2.3" or "1.2.3")
#   $2 - bump_type: Type of bump (major, minor, patch)
# Returns:
#   New version with 'v' prefix (e.g., "v1.3.0")
# Example:
#   new_version=$(bump_version "v1.2.3" "minor")  # Returns "v1.3.0"
bump_version() {
    local version="${1:-}"
    local bump_type="${2:-}"
    
    if [[ -z "$version" ]] || [[ -z "$bump_type" ]]; then
        echo "Error: version and bump_type are required" >&2
        return 1
    fi
    
    # Validate bump_type
    if [[ ! "$bump_type" =~ ^(major|minor|patch)$ ]]; then
        echo "Error: invalid bump_type: $bump_type (expected: major, minor, or patch)" >&2
        return 1
    fi
    
    # Parse version
    local major minor patch
    read -r major minor patch < <(parse_version "$version")
    
    # Bump according to type
    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# validate_version: Validate a semantic version format
# Args:
#   $1 - version: Version string to validate
# Returns:
#   0 if valid, 1 if invalid
# Example:
#   if validate_version "v1.2.3"; then echo "valid"; fi
validate_version() {
    local version="${1:-}"
    
    if [[ -z "$version" ]]; then
        return 1
    fi
    
    # Remove 'v' prefix if present
    version="${version#v}"
    
    # Check format
    if [[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        return 0
    else
        return 1
    fi
}

# get_version_from_tag: Extract version from a module tag
# Args:
#   $1 - tag: Tag name (e.g., "cache/inmemory/v1.2.3")
# Returns:
#   Version string (e.g., "v1.2.3")
# Example:
#   version=$(get_version_from_tag "cache/inmemory/v1.2.3")
get_version_from_tag() {
    local tag="${1:-}"
    
    if [[ -z "$tag" ]]; then
        echo "Error: tag is required" >&2
        return 1
    fi
    
    # Extract version part (everything after last /)
    local version="${tag##*/}"
    
    # Validate it's a version
    if ! validate_version "$version"; then
        echo "Error: tag does not contain valid version: $tag" >&2
        return 1
    fi
    
    echo "$version"
}

# Export functions for use in other scripts
export -f parse_version
export -f compare_versions
export -f bump_version
export -f validate_version
export -f get_version_from_tag
