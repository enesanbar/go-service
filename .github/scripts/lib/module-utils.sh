#!/usr/bin/env bash
# module-utils.sh - Go module utility functions
# Functions: discover_modules, validate_module, get_module_name

set -euo pipefail

# discover_modules: Find all Go modules in the repository
# Returns:
#   JSON array of module objects with path and name
# Example:
#   modules=$(discover_modules)
discover_modules() {
    local modules=()
    
    # Find all go.mod files
    while IFS= read -r go_mod; do
        # Get module directory (remove ./go.mod suffix)
        local module_path="${go_mod%/go.mod}"
        module_path="${module_path#./}"
        
        # Skip if empty (root go.mod)
        if [[ -z "$module_path" ]]; then
            continue
        fi
        
        # Get module name from go.mod
        local module_name
        module_name=$(get_module_name "$go_mod") || continue
        
        # Add to array
        modules+=("{\"path\":\"$module_path\",\"name\":\"$module_name\"}")
    done < <(find . -name "go.mod" -not -path "*/vendor/*" -not -path "*/.git/*" | sort)
    
    # Output as JSON array
    if [[ ${#modules[@]} -eq 0 ]]; then
        echo "[]"
    else
        echo "[$(IFS=,; echo "${modules[*]}")]"
    fi
}

# validate_module: Check if a path contains a valid Go module
# Args:
#   $1 - module_path: Path to check
# Returns:
#   0 if valid module, 1 otherwise
# Example:
#   if validate_module "cache/inmemory"; then echo "valid"; fi
validate_module() {
    local module_path="${1:-}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    # Check if go.mod exists
    if [[ ! -f "${module_path}/go.mod" ]]; then
        echo "Error: no go.mod found at ${module_path}/go.mod" >&2
        return 1
    fi
    
    # Check if module directive exists in go.mod
    if ! grep -q "^module " "${module_path}/go.mod"; then
        echo "Error: invalid go.mod at ${module_path}/go.mod (no module directive)" >&2
        return 1
    fi
    
    return 0
}

# get_module_name: Extract module name from go.mod file
# Args:
#   $1 - go_mod_path: Path to go.mod file
# Returns:
#   Module name as defined in go.mod
# Example:
#   name=$(get_module_name "cache/inmemory/go.mod")
get_module_name() {
    local go_mod_path="${1:-}"
    
    if [[ -z "$go_mod_path" ]]; then
        echo "Error: go_mod_path is required" >&2
        return 1
    fi
    
    if [[ ! -f "$go_mod_path" ]]; then
        echo "Error: go.mod not found at $go_mod_path" >&2
        return 1
    fi
    
    # Extract module name from first "module" directive
    local module_name
    module_name=$(grep "^module " "$go_mod_path" | head -n 1 | awk '{print $2}')
    
    if [[ -z "$module_name" ]]; then
        echo "Error: could not extract module name from $go_mod_path" >&2
        return 1
    fi
    
    echo "$module_name"
}

# get_go_version: Extract Go version requirement from go.mod
# Args:
#   $1 - module_path: Path to module directory
# Returns:
#   Go version (e.g., "1.21") or empty if not specified
# Example:
#   go_version=$(get_go_version "cache/inmemory")
get_go_version() {
    local module_path="${1:-}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    local go_mod="${module_path}/go.mod"
    if [[ ! -f "$go_mod" ]]; then
        echo "Error: go.mod not found at $go_mod" >&2
        return 1
    fi
    
    # Extract Go version from "go X.Y" directive
    local go_version
    go_version=$(grep "^go " "$go_mod" | awk '{print $2}' | head -n 1)
    
    echo "$go_version"
}

# has_tests: Check if module has test files
# Args:
#   $1 - module_path: Path to module directory
# Returns:
#   0 if tests exist, 1 otherwise
# Example:
#   if has_tests "cache/inmemory"; then echo "has tests"; fi
has_tests() {
    local module_path="${1:-}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    # Check for any *_test.go files
    if find "$module_path" -name "*_test.go" -type f | grep -q .; then
        return 0
    else
        return 1
    fi
}

# get_module_display_name: Get human-friendly display name for module
# Args:
#   $1 - module_path: Path to module (e.g., "cache/inmemory")
# Returns:
#   Display name (e.g., "cache-inmemory")
# Example:
#   display=$(get_module_display_name "cache/inmemory")
get_module_display_name() {
    local module_path="${1:-}"
    
    if [[ -z "$module_path" ]]; then
        echo "Error: module_path is required" >&2
        return 1
    fi
    
    # Replace slashes with hyphens
    echo "$module_path" | tr '/' '-'
}

# Export functions for use in other scripts
export -f discover_modules
export -f validate_module
export -f get_module_name
export -f get_go_version
export -f has_tests
export -f get_module_display_name
