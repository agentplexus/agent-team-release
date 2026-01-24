#!/bin/bash
# validate-release-command.sh
# PreToolUse hook to validate release-related bash commands

# This script receives the command being executed via stdin or arguments
# It can validate/modify commands or pass them through

# Read the tool input (command being executed)
COMMAND="${TOOL_INPUT:-$1}"

# If no command provided, pass through
if [ -z "$COMMAND" ]; then
    exit 0
fi

# Safety checks for release-related commands

# Check for dangerous git operations
if echo "$COMMAND" | grep -qE "git\s+(push\s+--force|reset\s+--hard|clean\s+-fd)"; then
    echo "Warning: Potentially destructive git command detected."
    echo "Command: $COMMAND"
    echo "Please confirm this is intentional."
    # Exit 0 to allow user to see warning but not block
    exit 0
fi

# Check for tag deletion
if echo "$COMMAND" | grep -qE "git\s+tag\s+-d|git\s+push.*--delete.*tag"; then
    echo "Warning: Tag deletion command detected."
    echo "Command: $COMMAND"
    echo "Deleting tags can cause issues with existing releases."
    exit 0
fi

# Validate version format for release commands
if echo "$COMMAND" | grep -qE "atrelease\s+release\s+"; then
    VERSION=$(echo "$COMMAND" | grep -oE "v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?")
    if [ -n "$VERSION" ]; then
        # Ensure version starts with 'v'
        if [[ ! "$VERSION" =~ ^v ]]; then
            echo "Note: Version '$VERSION' should typically start with 'v' (e.g., v$VERSION)"
        fi
    fi
fi

# Pass through - command is allowed
exit 0
