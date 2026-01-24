#!/bin/bash
# check-dependencies.sh
# SessionStart hook to verify required dependencies for release-agent plugin

set -e

# Colors for output (if terminal supports it)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track missing dependencies
MISSING_REQUIRED=()
MISSING_OPTIONAL=()

# Check required dependencies
check_required() {
    local cmd=$1
    local desc=$2
    local install_hint=$3

    if command -v "$cmd" &> /dev/null; then
        echo -e "${GREEN}[OK]${NC} $cmd - $desc"
    else
        echo -e "${RED}[MISSING]${NC} $cmd - $desc"
        echo "         Install: $install_hint"
        MISSING_REQUIRED+=("$cmd")
    fi
}

# Check optional dependencies
check_optional() {
    local cmd=$1
    local desc=$2
    local install_hint=$3

    if command -v "$cmd" &> /dev/null; then
        echo -e "${GREEN}[OK]${NC} $cmd - $desc"
    else
        echo -e "${YELLOW}[OPTIONAL]${NC} $cmd - $desc"
        echo "            Install: $install_hint"
        MISSING_OPTIONAL+=("$cmd")
    fi
}

echo "=== Release Agent Dependency Check ==="
echo ""

# Required: Core tools
echo "Core Tools:"
check_required "git" "Version control" "brew install git"
check_required "gh" "GitHub CLI for CI status" "brew install gh"

echo ""

# Required: Release Agent tools
echo "Release Agent Tools:"
check_required "atrelease" "Release automation CLI" "go install github.com/agentplexus/agent-team-release/cmd/atrelease@latest"
check_required "schangelog" "Changelog generation" "go install github.com/grokify/structured-changelog/cmd/schangelog@latest"
check_required "sroadmap" "Roadmap generation" "go install github.com/grokify/structured-roadmap/cmd/sroadmap@latest"

echo ""

# Language-specific tools (optional based on project)
echo "Language Tools (detected by project type):"

# Check for Go projects
if [ -f "go.mod" ]; then
    echo "  Go project detected:"
    check_required "go" "Go compiler" "brew install go"
    check_required "golangci-lint" "Go linter" "brew install golangci-lint"
    check_optional "gocoverbadge" "Coverage badge generator" "go install github.com/grokify/gocoverbadge@latest"
fi

# Check for Node.js projects
if [ -f "package.json" ]; then
    echo "  Node.js project detected:"
    check_required "node" "Node.js runtime" "brew install node"
    check_required "npm" "Package manager" "brew install node"
    check_optional "eslint" "JavaScript linter" "npm install -g eslint"
    check_optional "prettier" "Code formatter" "npm install -g prettier"
fi

# Check for Python projects
if [ -f "pyproject.toml" ] || [ -f "setup.py" ] || [ -f "requirements.txt" ]; then
    echo "  Python project detected:"
    check_optional "python3" "Python runtime" "brew install python"
    check_optional "pytest" "Test runner" "pip install pytest"
    check_optional "ruff" "Fast Python linter" "pip install ruff"
fi

# Check for Rust projects
if [ -f "Cargo.toml" ]; then
    echo "  Rust project detected:"
    check_optional "cargo" "Rust package manager" "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
fi

echo ""

# Set environment variables for subsequent commands
if [ -n "$CLAUDE_ENV_FILE" ]; then
    echo "export RELEASE_AGENT_PLUGIN_LOADED=1" >> "$CLAUDE_ENV_FILE"
    echo "export RELEASE_REPO_ROOT=$(pwd)" >> "$CLAUDE_ENV_FILE"
fi

# Summary
echo "=== Summary ==="
if [ ${#MISSING_REQUIRED[@]} -eq 0 ]; then
    echo -e "${GREEN}All required dependencies are installed.${NC}"
    exit 0
else
    echo -e "${RED}Missing ${#MISSING_REQUIRED[@]} required dependencies: ${MISSING_REQUIRED[*]}${NC}"
    echo ""
    echo "The release-agent plugin may not function correctly."
    echo "Please install the missing dependencies and restart Claude Code."
    # Don't exit with error to allow Claude to continue
    exit 0
fi
