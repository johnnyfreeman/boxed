#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

usage() {
    echo "Usage: $0 <success|error> <version> [duration]"
    echo ""
    echo "Example:"
    echo "  $0 success v1.2.3 \"2m 34s\""
    echo "  $0 error v1.2.3 \"1m 12s\""
    exit 1
}

if [ $# -lt 2 ]; then
    usage
fi

status="$1"
version="$2"
duration="${3:-unknown}"

case "$status" in
    success)
        box_type="success"
        title="Deploy Complete"
        commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
        branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

        $BOXED "$box_type" \
            --title "$title" \
            --subtitle "$version" \
            --kv "Duration=$duration" \
            --kv "Commit=$commit" \
            --kv "Branch=$branch" \
            --kv "Environment=${ENVIRONMENT:-production}" \
            --footer "Deployed at $(date '+%Y-%m-%d %H:%M:%S')"
        ;;
    error)
        box_type="error"
        title="Deploy Failed"

        $BOXED "$box_type" \
            --title "$title" \
            --subtitle "$version" \
            --kv "Duration=$duration" \
            --kv "Environment=${ENVIRONMENT:-production}" \
            --footer "Failed at $(date '+%Y-%m-%d %H:%M:%S')"
        ;;
    *)
        echo "Error: Invalid status. Use 'success' or 'error'"
        usage
        ;;
esac
